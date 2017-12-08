package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/OlivierBoucher/redis-cluster-operator/pkg/redis"
	// since everflow uses gcp, you can comment this for your own sake
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"flag"

	clientset "github.com/OlivierBoucher/redis-cluster-operator/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL  string
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func Main() int {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprint(os.Stderr, "instantiating logger failed: ", err)
		return 1
	}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building kubeconfig: %s", err.Error())
		return 1
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building kubernetes clientset: %s", err.Error())
		return 1
	}

	redisClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building redis-cluster clientset: %s", err.Error())
		return 1
	}

	ro, err := redis.NewOperator(kubeClient, redisClient)
	if err != nil {
		fmt.Fprint(os.Stderr, "instantiating redis-cluster controller failed: ", err)
		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error { return ro.Run(ctx.Done()) })

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		logger.Info("Received SIGTERM, exiting gracefully...")
	case <-ctx.Done():
	}

	cancel()
	if err := wg.Wait(); err != nil {
		logger.Error("Unhandled error received. Exiting...", zap.Error(err))
		return 1
	}

	return 0
}

func main() {
	os.Exit(Main())
}

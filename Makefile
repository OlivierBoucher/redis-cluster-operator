CURRENT_DIR=${CURDIR}
HACK_DIR ?= hack

generate:
	$(HACK_DIR)/update-codegen.sh
	# The following hack is required because of
	# https://github.com/kubernetes/code-generator/issues/20
	# we essentially copy the typed package from the lowercase generated
	# package and make sure that all the code in pkg/client uses
	# the good package name
	mv $(CURRENT_DIR)/../../olivierboucher/redis-cluster-operator/pkg/client/clientset/versioned/typed \
	$(CURRENT_DIR)/pkg/client/clientset/versioned/typed && rm -rf $(CURRENT_DIR)/../../olivierboucher
	find pkg/client -type f -print0 | xargs -0 sed -i 's|github.com/olivierboucher/redis-cluster-operator/|github.com/OlivierBoucher/redis-cluster-operator/|g'

.PHONY: generate


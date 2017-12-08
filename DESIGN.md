# Design

### Redis Spec

```
version: the tag of the official redis image to use
members: uint, min 3; the number of masters
replicationFactor: uint; the number of slaves for each master
```

### Implementation logic

Given N = members and R = replicationFactor, the total number of pods will be of order N*R

In order to determine whether a pod is a master or a slave, we use the ordinal index provided by the statefulset
and create windows for each masters, each master will be a multiple of R

Given 3 members and a replicationFactor of 2, our statefulset will consist of 6 pods.

```
Master --> Slave      Master --> Slave    Master --> Slave

[redis-0][redis-1]    [redis-2][redis-3]  [redis-4][redis-5]
```

The same windowing logic applies with a greater replication factor

Given 3 members and a replicationFactor of 3, our statefulset will consist of 9 pods.

```
Master --> Slave --> Slave     Master --> Slave --> Slave   Master --> Slave --> Slave

[redis-0][redis-1][redis-2]    [redis-3][redis-4][redis-5]  [redis-6][redis-7][redis-8]
```

WIP: more details to come

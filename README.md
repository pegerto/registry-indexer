# registry-indexer


The docker registry component is responsible only for the allocation and storage
of images data, but do not work effectivity managing metadata.

This PoC covers the required service to manage the registry metadata.


```

                      HTTPS

                        +
                        |
                        |
               +--------v-------------+
               |                      |
         +-----+        NGINX         +-----+
         |     |                      |     |
         |     |                      |     |
/v2/*    |     +----------------------+     | /v2/_catalog
         |                                  |
         |                                  |
         |                                  |
         |                                  |
         |                                  |
+--------v----------+           +-----------v--------+
|                   |           |                    |
|                   <-----------+                    |
|     REGISTRY      |           |      INDEXER       |
|                   |           |                    |
|                   |           |                    |
+---------+---------+           +----------^---------+
          |                                |
          |                                |
          |                                |
          |                                |
          |                                |
          +--------------------------------+
                    NOTIFY INDEXER


```

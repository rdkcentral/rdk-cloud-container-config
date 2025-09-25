# rdk-cloud-container-config
Hosts container files for RDK Cloud Components such as Webpa and XConf
=======
# rdk-cloud-container-config/xconf

Tree structure of xconf folder is as follows

├── cassandra
│   └── db
│       ├── db_create_tables.cql
│       └── db_init.cql
├── docker-compose.yaml
├── init-cassandra.sh
├── xconfadmin
│   ├── Dockerfile
│   └── config
│       └── sample_xconfadmin.conf
├── xconfui
│   ├── Dockerfile
│   └── config
│       └── sample_xconfui.conf
└── xconfwebconfig
    ├── Dockerfile
    └── config
        └── sample_xconfwebconfig.conf


Clone the rdk-cloud-container-config repo to your workspace

From the xconf folder,
run below command

$ docker compose up --build -d


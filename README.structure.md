### project structure

```

├── cmd
│   └── aapiservice // name of main server
├── config // all configuration beusing by main server
├── docker
│   ├── mongodb // initial data
│   ├── mongodb_installer //installer for multiple replication and master DB
│   ├── mysql_installer // handle installer
│   │   ├── installer // installer config and scripts
│   │   ├── master // master db config and scripts
│   │   └── replicas // replication db config and script
│   └── mysqlseed //initial data for mysql
├── externals
│   // all external services will be there
│       
├── integration_protos
│   └── api_protos
│       └── v1
│           ├── mobile // protobufs for mobile
│           └── web // protobufs for web
├── internal
│   ├── interceptors // middleware that will be handle all request before executed by controller
│   └── validations // validation for request, rpc calls
├── kubernetes // all kubernetes deployment
├── migrations // all sql, mongodb migrations 
├── pkg //imagine as Common library
│   ├── grpc_errors // grpc middleware to handle errors 
│   ├── jaeger // monitoring 
│   ├── kafka_srv // message streaming middleware
│   ├── logger // logger middleware
│   ├── metrics // monitoring system middleware
│   ├── mime_types // common type of request middle ware
│   ├── minio_srv // imaging middleware
│   ├── mongodb //mongodb middleware
│   ├── mysql //mysql middleware
│   ├── rabbitmq //messaging middleware
│   ├── redis //redis middleware
│   ├── trace //tracing middleware
│   └── utils //others
├── protos
│   └── v1 //protobufs merge all protobufs of services to provide full protos all in one
├── scripts // all scripts to start the application
├── services
│   ├── authentication // authentication service
│   │   ├── models //models
│   │   ├── protos //related protos 
│   │   │   ├── authentication //proto of current service
│   │   │   │   └── v1 //version of proto
│   │   │   ├── device //proto of another service
│   │   │   │   └── v1 //version of proto
│   │   │   ├── administrator //proto of another service
│   │   │   │   └── v1 //version of proto
│   │   │   └── user //proto of another service
│   │   │       └── v1 //version of proto
│   │   ├── repositories //database communication and caching 
│   │   ├── services //handle request to authentication service, and decide to combine data from repos, usercases
│   │   └── usercases // business logic
│   ├── device // device service
│   │   ├── models //models
│   │   ├── protos //related protos 
│   │   │   └── device //proto of current service
│   │   │       └── v1 //version of proto
│   │   ├── repositories //database communication and caching 
│   │   └── services //handle request to device service, and decide to combine data from repos, usercases
│   ├── logging  // logging service (background service)
│   │   ├── models //models
│   │   ├── protos //related protos 
│   │   │   └── log //proto of current service
│   │   │       └── v1 //version of proto
│   │   ├── repositories //database communication and caching 
│   │   └── services //handle request to logging service, and decide to combine data from repos, usercases
│   ├── administrator
│   │   ├── models
│   │   ├── protos
│   │   │   └── administrator
│   │   │       └── v1
│   │   ├── repositories
│   │   ├── services
│   │   └── usecases
│   └── user
│       ├── models
│       ├── protos
│       │   ├── ai_engine
│       │   │   └── v1
│       │   ├── faiss
│       │   │   └── v1
│       │   └── user
│       │       └── v1
│       ├── repositories
│       └── services
└── shared //all models, functions, custom middleware, consts, cache, keys...
    ├── constants
    ├── http_clients
    └── utilities

```
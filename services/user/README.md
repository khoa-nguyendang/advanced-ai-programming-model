DEVELOP service

give permission for proto-gen.sh (first time)

```
chmod +x ./proto-gen.sh

# from root dir
chmod +x ./services/user/proto-gen.sh
```

generate proto by trigger

```
./proto-gen.sh input_proto_path

#example
./proto-gen.sh ./protos/user/v1/user.proto
./proto-gen.sh ./protos/faiss/v1/faiss.proto
./proto-gen.sh ./protos/ai_engine/v1/ai_engine.proto

# from root dir
./services/user/proto-gen.sh services/user/protos/user/v1/user.proto
./services/user/proto-gen.sh services/user/protos/faiss/v1/faiss.proto
./services/user/proto-gen.sh services/user/protos/ai_engine/v1/ai_engine.proto
```
# App FAISS Index Service

[![CircleCI](https://circleci.com/gh/AppCommercial/aapi-faiss-service.svg?style=svg&circle-token=c1b422a1c7711b5ce80bb3fe9be9e6fffac03caf)](https://circleci.com/gh/AppCommercial/aapi-faiss-service)

## Why do we need this?

Facebook FAISS index can be used directly in code as a library. A very basic implementation of it could be simply done as part of the cloud service code. However, `faiss` has a couple of main problems that prevent it from being scalable to our level of satisfaction.

- The index works the best in-memory, which means each Flask worker process would need its own `faiss` instance. This makes the memory requirement for `faiss` quickly grow 5x or 10x per API server
- In order to use `faiss` in memory in the first place, we also need to load it from disk, which is an expensive operation. We want to limit this to as few invocations as possible.

Besides, we also noticed that `faiss` only needs CPU to run so in theory it can scale much more cheaply if done as an independent microservice. With this approach, we can build many more features such as load management, index sharding, failover, replica distribution etc in the future.

## API Documentation

### Server config

- Currently, this is running as a Flask server.
- All server configuration can be overridden inside `config.py`.
- In production, we can deploy these config via our `fab` deployment script.

### Documentation

#### Vectors

`faiss` performs nearest-neighbour searches on vectors. Of course, numpy vectors/arrays are not valid JSON objects so for the input vectors, we define each vector as an array of numbers.

#### `GET /health/`

- This is a simple `GET` request for healthcheck.
- Sample response:
```
Success: HTTP 200 (the server has loaded its index from filesystem successfully)
{}

Failure: HTTP 500
{"reason": "Index is not loaded correctly"}
```

#### `POST /register/`

- This request takes an `application/json` header.
- In order for `faiss` to perform nearest-neighbour searches, its index needs to be registered with input vectors in the first place.
- This request can register a batch of vectors along with their IDs to `faiss`.
- Notes:
    - `faiss` only accepts `int64` values for IDs
    - Each `faiss` index has a fixed dimension that needs to be specified upon creation, so if any input vector has a different dimension than this number, it will fail to register and result in an error. The default value for this is `512` (see `config.py`)
- Each vector is represented as a tuple [vector_id, vector]
- Sample request body:
```
{
    "vectors": [
        [267400127, [24.0, 1.2, 12671.0 ...]]
        ...
    ]
```
- Sample response:
```
Success: HTTP 200
{}

Failure: HTTP 400, failure to parse JSON
{"reason": "Malformed request body"}

Failure: HTTP 400, `vector` has wrong data type
{"reason": "`vector` at ID `267400127` should be an array"}

Failure: HTTP 400, `vector` has wrong dimension
{"reason": "`vector` at ID `267400127` should contain 512 entries"}
```

#### `POST /deregister/`

- This request takes an `application/json` header.
- This request deregisters a batch of vector IDs from the index.
- This means subsequent searches will not return the deregistered ID in the results.
- Sample request body:
```
{
    "vector_ids": [267400127, 267400128]
}
```
- Sample response:
```
Success: HTTP 200
{}
```

#### `POST /search/`

- This request takes an `application/json` header.
- This request performs a nearest-neighbour search on the `faiss` index service.
- It matches a single vector at a time for now. In the future, we should perform this search in batches to make the best use of `faiss`.
- Sample request body:
```
{
    "vector": [24.0, 1.2, 12671.0 ...]
```
- Sample response:
```
Success: HTTP 200
{}

Failure: HTTP 400, failure to parse JSON
{"reason": "Malformed request body"}

Failure: HTTP 400, `vector` has wrong data type
{"reason": "`vector` should be an array"}

Failure: HTTP 400, `vector` has wrong dimension
{"reason": "`vector` should contain 512 entries"}
```

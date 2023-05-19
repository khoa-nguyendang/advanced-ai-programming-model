To view this diagram in vscode, 
try to install extension: bierner.markdown-mermaid
then view this file with preview mode.

```mermaid
sequenceDiagram
%% Enrollment
    participant w  as Web
    participant  ag as ApiGateway
    participant  us as UserService
    participant d as Object Detect Model
    participant fr as Face Regcoganize
    participant fs as Faiss Service
    %%participant lg as LoggingService
    participant mg as MongoDb
    %%participant kf as Kafka
    participant ms as MySQL

    

    w->>+ag: Enroll
    ag->>+us: Check user exists
    us->>+d: retrieve huma
    d->>+fr: extract face feature
    fr->>+fs: find similar vectors
    fs->>fs: load index from database
    fs->>+fs: normalize_L2 to <br/> search cosine similarity.
    alt is Exists
        fs-->>-fr: already exists
        fr-->>us: return already exists error
        us-->>ag: return already exists error
        ag-->>w: already exists error
    else is new faces
        fs->>+mg: add new embedding vector
        mg-->>-fs: return uuid 
        fs-->>-fr: return uuid
        fr-->>-us: return uuid
        us->>+ms: add user's info
        ms-->>-us: return state changes
        us-->>-ag: return enroll result
        ag-->>-w: return enroll result
        w->>w: refresh to get <br/> new users
    end
```
    



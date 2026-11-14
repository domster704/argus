# Argus

> Argus is a many-eyed giant from ancient Greek mythology, as well as the name of a government information system and a planet in a popular universe

## Start

1. Server
```bash
go run .\cmd\collector\ --listen=127.0.0.1:50052    
```

2. Client
```Bash
go run .\cmd\agent\ --id=home-pc --server=127.0.0.1:50052 --interval=1s
```
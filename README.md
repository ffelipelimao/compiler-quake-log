# Quake log Compiler


##  Why compiler?


```mermaid
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
```

## Test

`go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out`
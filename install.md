### windows

```bat
##编译为linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go  build -o mvncleaner main.go
##编译为macosx
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o mvncleaner main.go
##编译为windows
SET GOOS=windows
go build -o mvncleaner.exe  main.go
```

### macosx

```shell
##编译为windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mvncleaner.exe main.go
##编译为linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mvncleaner main.go
```

### linux

```shell
##编译为windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mvncleaner.exe main.go
##编译为macosx
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o mvncleaner main.go
```


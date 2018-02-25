# GeekApk-Server

## Initialization

```shell
git clone https://github.com/geekapk/GeekApk-Server
cd GeekApk-Server
export GOPATH=`pwd`

go get ...

```

## Start

```shell
go run launch.go
```

## Package management

Currently we just use the repository root as GOPATH and put the server code in src/ to keep things simple. Feel free to open an issue if you've got a better dependency management strategy!

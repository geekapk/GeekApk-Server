# GeekApk-Server

## Initialization

```shell
git clone https://github.com/geekapk/GeekApk-Server
cd GeekApk-Server
export GOPATH=`pwd`

rm -rf src/github.com src/golang.org

go get github.com/gorilla/context github.com/gorilla/securecookie github.com/gorilla/sessions github.com/jinzhu/gorm github.com/jinzhu/inflection github.com/lib/pq github.com/satori/go.uuid

go get golang.org/x/crypto

```

## Start

```shell
go run launcher
```


I think maybe `vendor` is better than full gopath
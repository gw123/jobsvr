-include .env

APPNAME=jobsvr

genpb:
	protoc --go_out=plugins=grpc:.  *.proto

server:
	go run entry/main.go job

.PHONY: send
send:
	go run demo/demo.go -action send

.PHONY: listen
listen:
	go run demo/demo.go -action listen

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s' -o $(APPNAME) ./entry/main.go &&\
	upx -9 -f -o upload $(APPNAME)

.PHONY: upload
upload:
	 scp upload root@sh2:/data/apps/jobsvr

restart:
	ssh root@sh2 supervisorctl restart jobsvr

all: build upload restart


build-alpine:
	@docker run --rm -v "$(PWD)":/go/src/github.com/gw123/jobsvr \
	    -e GOPROXY=https://goproxy.cn \
	    -e GOPRIVATE=github.com/gw123/jobsvr \
		-w /go/src/github.com/gw123/jobsvr \
		golang:1.15.2-alpine3.12 \
		go build -v -ldflags '-w -s' -o $(APPNAME) github.com/gw123/jobsvr/entry &&\
		upx -6 -f -o upload ./$(APPNAME)

build-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags  '-w -s  -extldflags "-static"' -o $(APPNAME)  ./entry/main.go &&\
	upx -6 -f -o upload ./$(APPNAME)

docker-image:
	chmod +x upload &&\
    cp upload $(DOCKER_BUILD_PATH)/$(APPNAME)  &&\
    cp Dockerfile $(DOCKER_BUILD_PATH)/ &&\
	@docker build -t $(REMOTE_USER_API_TAG) $(DOCKER_BUILD_PATH)
	@docker push $(REMOTE_USER_API_TAG)

## 在宿主机器上静态打包， 打包体积大但是速度快， 适合开发阶段
docker-all: build-static docker-image
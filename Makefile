-include .env

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


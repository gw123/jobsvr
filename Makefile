genpb:
	protoc --go_out=plugins=grpc:.  *.proto

server:
	go run cmd/main.go

client:
	go run demo/demo.go
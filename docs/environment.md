# protobuf and grpc setting
install protobuf and grpc-go in host os
(macos): brew install protobuf
(ubutun): sudo apt-get install protobuf
> govendor get -u github.com/grpc/grpc-go
> go get -a github.com/golang/protobuf/protoc-gen-go

generate protobuffer files
> cd registry/api
> protoc -I. registry.proto --go_out=plugins=grpc:.

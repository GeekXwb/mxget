protoc --proto_path=api/grpc --proto_path=third_party --gogo_out=plugins=grpc:pkg/api api/grpc/*.proto
protoc --proto_path=api/grpc --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api api/grpc/*.proto

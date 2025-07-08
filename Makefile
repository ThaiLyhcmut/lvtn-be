proto-auth:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       proto/auth/auth.proto
proto-common:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       proto/common/common.proto
proto-async:
	protoc --go_out=. --go_opt=paths=source_relative        --go-grpc_out=. --go-grpc_opt=paths=source_relative       proto/asynq/asynq.proto
service-common:
	go run services/_common/main.go
service-auth:
	go run services/auth/main.go
generate-grpc:
	protoc --proto_path=./proto \
	  --go_out=./auth_service/services/ \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=./auth_service/services/ \
	  --go-grpc_opt=paths=source_relative \
	  --go-grpc_opt=require_unimplemented_servers=false \
	  ./proto/auth.proto
	protoc --proto_path=./proto \
	  --go_out=./task_service/services/ \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=./task_service/services/ \
	  --go-grpc_opt=paths=source_relative \
	  ./proto/auth.proto

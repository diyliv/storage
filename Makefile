gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative proto/storage/storage.proto
start-redis:
	docker pull redis && docker run --name rediska -d redis && docker exec -it rediska /bin/bash	
test:
	cd tests/ && docker-compose up
run:
	docker-compose -f deployments/docker-compose.yaml up --build --remove-orphans
generate:
	protoc pb/*.proto --proto_path=. \
         --go_out=pb/ --go_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb \
         --go-grpc_out=pb/ --go-grpc_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb
	go generate ./...

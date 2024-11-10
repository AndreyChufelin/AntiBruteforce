run:
	docker-compose -f deployments/docker-compose.yaml up --build --remove-orphans
generate:
	protoc pb/*.proto --proto_path=. \
         --go_out=pb/ --go_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb \
         --go-grpc_out=pb/ --go-grpc_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb
	go generate ./...

ingtegration-tests:
	docker-compose --env-file ./deployments/.env.tests -f ./deployments/docker-compose.yaml run --build --rm tests


.PHONY: run generate build test integration-tests cli
COMPOSE_FILE=deployments/docker-compose.yaml

run:
	docker-compose -f ${COMPOSE_FILE} up --build --remove-orphans

build:
	go build -o ./bin/antibruteforce ./cmd/ratelimiter

generate:
	protoc pb/*.proto --proto_path=. \
         --go_out=pb/ --go_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb \
         --go-grpc_out=pb/ --go-grpc_opt=module=github.com/AndreyChufelin/AntiBruteforce/pb
	go generate ./...

test:
	go test -race -count 100 -v ./...

integration-tests:
	./tests/integration/start.sh

cli:
	docker-compose -f deployments/docker-compose.yaml run cli

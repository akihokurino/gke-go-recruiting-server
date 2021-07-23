MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export PATH := $(ROOT)/scripts:$(PATH)

# proceed-work-status
BATCH := ""

# reindex-search
# reorder-work
JOB := ""

vendor:
	go mod tidy

gen-proto:
	mkdir -p proto/go
	rm -rf proto/go/*
	protoc --proto_path=proto/. --twirp_out=proto/go --go_out=proto/go proto/*.proto

gen-injector:
	cd di && wire

gen: gen-proto gen-injector

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o .build/main -a -installsuffix cgo entrypoint/api/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o .build/main -a -installsuffix cgo entrypoint/batch/main.go

build_script:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o .build/main -a -installsuffix cgo script/sync-image/main.go

run-local:
	docker-compose up

stop-local:
	docker-compose stop

run-local-batch:
	docker-compose run --rm batch make gen && go run /app/entrypoint/batch/main.go ${BATCH}

run-local-job:
	docker-compose run --rm batch make gen && go run /app/entrypoint/batch/main.go ${JOB}

setup-k8s:
	setup-k8s.sh

deploy-secret:
	deploy-secret.sh

deploy-api: deploy-secret
	deploy-api.sh

deploy-batch: deploy-secret
	deploy-batch.sh

deploy-job: deploy-secret
	deploy-job.sh k8s/job/${JOB}.yaml

deploy: deploy-secret
	deploy-api.sh
	deploy-batch.sh

proxy_db:
	cloud_sql_proxy -credential_file=./config/gcp-key.json -instances=gke-go-sample:asia-northeast1:db=tcp:0.0.0.0:3306

format:
	find . -print | grep --regex '.*\.go' | xargs goimports -w -local
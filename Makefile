.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/publishsns main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/recivemsg sqs_lambda/main.go
clean:
	rm -rf ./bin

deploy: build
	sls deploy --verbose

deploy-local: build
	serverless deploy --stage local --verbose

sls-debug:
	SLS_DEBUG=* sls offline --useDocker



default: test
version=local

COVERAGE=80

build: covercheck swagger
	@docker build -t walissoncasonatto/user-service:${version} .

publish: build
	docker push walissoncasonatto/user-service:${version}
	
covercheck: cover
	@sh -c "'$(CURDIR)/scripts/coverage.sh' ${COVERAGE}"

cover: test
	@sh -c "'$(CURDIR)/scripts/cover.sh'"

test: lint regenerate
	@sh -c "'$(CURDIR)/scripts/test.sh'"

lint:
	@sh -c "'$(CURDIR)/scripts/lint.sh'"

run:
	@docker-compose up

build-run: build run

coverreport: cover
	@xdg-open .coverage/coverage.html
	@sh -c "'$(CURDIR)/scripts/coverage.sh' ${COVERAGE}"

clean-mock:
	find internal -iname '*_mock.go' -exec rm {} \;

wire:
	go get github.com/google/wire/cmd/wire
	@wire ./...

generate: wire
	go install github.com/golang/mock/mockgen@v1.5.0
	@go generate ./...

regenerate: clean-mock generate

swagger:
	@go get github.com/swaggo/swag/cmd/swag
	@go mod vendor
	@swag init -d cmd/ --parseDependency --parseDepth 6

install-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0

e2e-test:
	docker build -f dockerfile.e2e -t walissoncasonatto/user-service:e2e --build-arg ENV_FILE=e2e/environment.json .
	docker run --rm walissoncasonatto/user-service:e2e run --environment="/etc/environment.json"  /etc/newman/collection.json
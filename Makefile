.PHONY: main build unit  test mockery swagger swagger-ui swagger-client start openapi 

MODULES_TO_UNITTEST := $(shell go list ./... | grep -v "test\|mocks")

unit :
	@echo "Running unit tests for the modules: $(MODULES_TO_UNITTEST)"
	go test $(MODULES_TO_UNITTEST)

mockery :
	@echo "Generating mocks for storages"
	@echo "Add --Name or --All to generate mocks"
	mockery --all --output="internal/mocks" --with-expecter

start:
	go build && ./order-service
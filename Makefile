go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

# Run all the code generators in the project
.PHONY: generate
generate: go.sum
	go generate ./...


.PHONY: build-check
build-check: go.sum
		go build -o ./bin/check -mod=readonly $(BUILD_FLAGS) ./checks/main

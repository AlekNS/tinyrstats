usage:
	@echo "Use targets: clean, build, linter, format, dep_ensure, $(cmd)"

targets = $(addprefix bin/, $(cmd))

dep_ensure:
ifeq ($(skip_dep), false)
	@echo "Running go mod"
	@go mod vendor
	@echo "Done"
else
	@echo "go mod disabled"
endif

$(go_bin)/golint:
	@echo "Installing golint"
	@go get -v -u golang.org/x/lint/golint
	@echo "Done"

bin/%: $(go_src) | dep_ensure
	$(eval bin_name=$(patsubst bin/%,%,$@))
	@mkdir -p bin
	@echo "Building bin $(bin_name)"
	@cd ./cmd/$(bin_name) && go build -o ../../bin/$(bin_name) -ldflags "-w" main.go
	@echo "Done"

test: | dep_ensure
	@go test -cpu 2 -v -parallel 2 -cover -race ./internal/... ./pkg/...
	@echo "Done"

clean:
	@echo "Clean $(targets)"
	@rm -vf $(targets)
	@echo "Done"

build: $(targets)

linter: $(go_bin)/golint
	@echo "Running golint"
	@$(foreach dir,$(lintfolders),golint -set_exit_status $(dir)/...;)
	@echo "Done"

format:
	@echo "Formating all go code"
	@$(foreach dir,$(lintfolders),gofmt -s -w $(dir);)
	@echo "Done"

$(cmd): bin/$(cmd)
	@echo "Done $@"

PHONY: usage clean linter format build dep_ensure $(cmd)

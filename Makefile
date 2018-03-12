ALL_SRC                = $(shell find . -name "*.go" | grep -v -e vendor)
BINARY                 = $(shell echo $${PWD\#\#*/})
PACKAGES               = $(shell go list ./... | grep -v -E 'vendor')
NON_GENERATED_PACKAGES = $(shell go list ./... | grep -v -E 'vendor|client|apis') # TODO: By excluding /apis/, we are excluding types.go which is a non-generated file.
RACE                   = -race
GOTEST                 = go test -v $(RACE)
GOLINT                 = golint
GOVET                  = go vet
GOFMT                  = gofmt
ERRCHECK               = errcheck -ignoretests
FMT_LOG                = fmt.log
LINT_LOG               = lint.log
DEP_VERSION            = v0.4.1
PASS                   = $(shell printf "\033[32mPASS\033[0m")
FAIL                   = $(shell printf "\033[31mFAIL\033[0m")
COLORIZE               = sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(ALL_SRC) unit_test lint codegen_verify

.PHONY: up
up:
	./hack/kanali-up.sh

.PHONY: unit_test
unit_test:
	@bash -c "set -e; set -o pipefail; $(GOTEST) $(PACKAGES) | $(COLORIZE)"

.PHONY: e2e_test
e2e_test:
	@./hack/e2e.sh

.PHONY: test
test: unit_test e2e_test

.PHONY: install
install:
	(dep version | grep $(DEP_VERSION)) || (mkdir -p $(GOPATH)/bin && DEP_RELEASE_TAG=$(DEP_VERSION) curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep version)
	dep ensure -v -vendor-only # assumes updated Gopkg.lock

.PHONY: fmt
fmt:
	@$(GOFMT) -e -s -l -w $(ALL_SRC)
	@./hack/updateLicenses.sh

.PHONY: cover
cover:
	@./hack/cover.sh $(shell go list $(NON_GENERATED_PACKAGES))
	@go tool cover -html=cover.out -o cover.html

.PHONY: binary
binary:
	@./hack/binary.sh $(VERSION)

.PHONY: lint
lint:
	@$(GOVET) $(NON_GENERATED_PACKAGES)
	@$(ERRCHECK) $(NON_GENERATED_PACKAGES)
	@cat /dev/null > $(LINT_LOG)
	@$(foreach pkg, $(NON_GENERATED_PACKAGES), $(GOLINT) $(pkg) >> $(LINT_LOG) || true;)
	@[ ! -s "$(LINT_LOG)" ] || (echo "Lint Failures" | cat - $(LINT_LOG) && false)
	@$(GOFMT) -e -s -l $(ALL_SRC) > $(FMT_LOG)
	@[ ! -s "$(FMT_LOG)" ] || (echo "Go Fmt Failures, run 'make fmt'" | cat - $(FMT_LOG) && false)

.PHONY: install_ci
install_ci: install
	go get github.com/wadey/gocovmerge
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover
	go get golang.org/x/lint/golint
	go get github.com/kisielk/errcheck

.PHONY: codegen_verify
codegen_verify:
	@./hack/verify-codegen.sh

.PHONY: codegen_update
codegen_update:
	@./hack/update-codegen.sh
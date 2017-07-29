ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor \
        -e ".*/\..*" \
        -e ".*/_.*" \
        -e ".*/mocks.*")

BINARY=$(shell echo $${PWD\#\#*/})
FILES = $(shell go list ./... | grep -v /vendor/)
PACKAGES := $(shell glide novendor)

RACE=-race
GOTEST=go test -v $(RACE)
GOLINT=golint
GOVET=go vet
GOFMT=gofmt
ERRCHECK=errcheck -ignoretests
FMT_LOG=fmt.log
LINT_LOG=lint.log

PASS=$(shell printf "\033[32mPASS\033[0m")
FAIL=$(shell printf "\033[31mFAIL\033[0m")
COLORIZE=sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(ALL_SRC) test fmt lint

.PHONY: install
install:
	glide --version || go get github.com/Masterminds/glide
	glide install

.PHONY: build
build:
	GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o $(BINARY)

.PHONY: fmt
fmt:
	$(GOFMT) -e -s -l -w $(ALL_SRC)
	./scripts/updateLicenses.sh

.PHONY: cover
cover:
	./scripts/cover.sh $(shell go list $(PACKAGES))
	go tool cover -html=cover.out -o cover.html

.PHONY: test
test:
	bash -c "set -e; set -o pipefail; $(GOTEST) $(PACKAGES) | $(COLORIZE)"

.PHONY: lint
lint:
	@$(GOVET) $(PACKAGES)
	@$(ERRCHECK) $(PACKAGES)
	@cat /dev/null > $(LINT_LOG)
	@$(foreach pkg, $(PACKAGES), $(GOLINT) $(pkg) >> $(LINT_LOG) || true;)
	@[ ! -s "$(LINT_LOG)" ] || (echo "Lint Failures" | cat - $(LINT_LOG) && false)
	@$(GOFMT) -e -s -l $(ALL_SRC) > $(FMT_LOG)
	@[ ! -s "$(FMT_LOG)" ] || (echo "Go Fmt Failures, run 'make fmt'" | cat - $(FMT_LOG) && false)

.PHONY: install_ci
install_ci: install
	go get github.com/wadey/gocovmerge
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover
	go get github.com/golang/lint/golint
	go get github.com/kisielk/errcheck

.PHONY: test_ci
test_ci:
	@./scripts/cover.sh $(shell go list $(PACKAGES))
	make lint

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
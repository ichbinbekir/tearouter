GO := go

run:
	$(GO) run $(CURDIR)/cmd/tearouter

build:
	$(GO) build -o $(CURDIR)/bin/ $(CURDIR)/cmd/tearouter

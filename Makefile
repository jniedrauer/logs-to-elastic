# vim: set noet

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Build config
BUILDFLAGS = -v
BUILDDIR = build
APPDIR = internal/app
TARGETS = cloudwatch elb
BINARIES = $(addprefix $(BUILDDIR)/, $(TARGETS))
OUTPUTS = $(addsuffix .zip, $(BINARIES))

all: test build

build: $(BINARIES) $(OUTPUTS) deps

$(BUILDDIR)/%.zip: $(BUILDDIR)/%
	cd $(BUILDDIR) && zip $(notdir $@) $(notdir $<)

$(BUILDDIR)/%: $(APPDIR)/%/main.go
	$(GOBUILD) $(BUILDFLAGS) -o $@ $<

test: deps
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARIES) $(OUTPUTS)

deps:
	$(GOGET) ./...

.PHONY: build test clean deps

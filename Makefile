# vim: set noet

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean -cache -testcache
GOTEST = $(GOCMD) test
DEP = dep
DEPENSURE = $(DEP) ensure

# Build config
BUILDFLAGS = -v
BUILDDIR = build
APPDIR = cmd
TARGETS = cloudwatch elb
BINARIES = $(addprefix $(BUILDDIR)/, $(TARGETS))
OUTPUTS = $(addsuffix .zip, $(BINARIES))

all: test build

build: deps $(BINARIES) $(OUTPUTS)

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
	$(DEPENSURE) -v

.PHONY: build test clean deps

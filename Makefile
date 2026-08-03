PRI_DIR = ./internal/
PUB_DIR = ./pkg/
BIN_DIR = ./bin/
RELEASE_DIR = ./release/
EXEC = roguefront

#Get OS and configure based on OS
ifeq ($(OS),Windows_NT)
    DISTRO ?=windows
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        ARCH ?=amd64
    else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		    ARCH ?=amd64
		endif
		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		    ARCH ?=ia32
		endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
   		DISTRO ?=linux
    endif
    ifeq ($(UNAME_S),Darwin)
   		DISTRO ?=darwin
    endif
    ifeq ($(UNAME),Solaris)
	   	DISTRO ?=solaris
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        ARCH ?=amd64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        ARCH ?=ia32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        ARCH ?=arm64
    endif
endif


.PHONY: run
#: Starts the project
run: build
	@$(BIN_DIR)$(EXEC)

.PHONY: build
#: Performs a clean build of the project
build: .deps $(PRI_DIR)** $(PUB_DIR)**
	@CGO_ENABLED=1 \
	GOOS=$(DISTRO) \
	GOARCH=$(ARCH) \
	CGO_CFLAGS=$(CFLAGS) \
	CGO_CPPFLAGS=$(CFLAGS) \
	CGO_LDFLAGS=$(LDFLAGS) \
	go build -o $(BIN_DIR)$(EXEC) cmd/main.go

build-wasm:
	@DISTRO=js \
	ARCH=wasm \
	$(MAKE) --no-print-directory build
build-macos:
	@DISTRO=darwin \
	ARCH=amd64 \
	$(MAKE) --no-print-directory build
build-ubuntu:
	@DISTRO=linux \
	ARCH=amd64 \
	$(MAKE) --no-print-directory build
build-windows:
	@DISTRO=windows \
	ARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	CXX=x86_64-w64-mingw32-g++ \
	$(MAKE) --no-print-directory build

.PHONY: release
#: packages release target
release: build .deps

.PHONY: clean
#: Cleans build files from project
clean:
	@rm $(BIN_DIR)$(EXEC) || true;
	@rm $(RELEASE_DIR)$(EXEC) || true;

.PHONY: clean-all
#: Cleans slate for project
clean-all: clean
	@rm .deps || true;

# deps include target
.PHONY: deps
.deps:
	@$(MAKE) --no-print-directory deps
	@go mod tidy

#: Install dependencies for running this project
deps:
	@sudo apt-add-repository --update ppa:longsleep/golang-backports
	@sudo apt install golang
	@sudo apt install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libglx-dev libgl1-mesa-dev libxxf86vm-dev
	@touch .deps

.PHONY: help
#: Lists available commands
help:
	@echo "Available Commands for project:"
	@grep -B1 -E "^[a-zA-Z0-9_\.-]+\:" Makefile | grep -A1 -E "^#\: \w+" \
	 | grep -v -- -- \
	 | sed 'N;s/\n/###/' \
	 | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	 | column -t  -s '###'

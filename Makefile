BUILD_DIR = ./build/
PRI_DIR = ./internal/
PUB_DIR = ./pkg/
BIN_DIR = ./bin/
RELEASE_DIR = ./release/
EXEC = roguefront
GCC = gcc
GCP = g++
CFLAGS = -I$(CURDIR)/include/raylib-go/raylib
LDFLAGS = -L$(CURDIR)/lib

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

DEPS ?=linux-amd64
ifeq ($(DISTRO),linux)
	ifeq ($(ARCH),amd64)
		DEPS ?=linux-amd64
	endif
	ifeq ($(ARCH),arm64)
		DEPS ?=linux-arm64
	endif
endif
ifeq ($(DISTRO),windows)
	DEPS ?=windows
endif

.PHONY: run
#: Starts the project
run: build .deps
	@$(BIN_DIR)$(EXEC)

.PHONY: build
#: Performs a clean build of the project
build: .dev-deps $(PRI_DIR)** $(PUB_DIR)**
	@CGO_ENABLED=1 \
	GOOS=$(DISTRO) \
	GOARCH=$(ARCH) \
	CGO_CFLAGS=$(CFLAGS) \
	CGO_CPPFLAGS=$(CFLAGS) \
	CGO_LDFLAGS=$(LDFLAGS) \
	go build -tags x11,opengl43 -o $(BIN_DIR)$(EXEC) cmd/main.go

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

.PHONY: clean-all
#: Cleans slate for project
clean-all:
	@rm .dev-deps || true;
	@rm .deps || true;

# deps include target
.PHONY: deps
.deps:
	@$(MAKE) --no-print-directory deps

#: Install dependencies for running this project
deps:
	@touch .deps

# dev-deps include target
.dev-deps:
	@$(MAKE) --no-print-directory dev-deps

# dev-deps for linux
.PHONY: .dev-deps-linux-amd64
.dev-deps-linux-amd64:
	@sudo apt-get install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev
	@sudo apt-get install -y libgl-dev libx11-dev xorg-dev libxxf86vm-dev

.PHONY: .dev-deps-linux-arm64
.dev-deps-linux-arm64:
	@sudo apt-get install -y libgl1-mesa-dev:arm64 libxi-dev:arm64 libxcursor-dev:arm64 libxrandr-dev:arm64 libxinerama-dev:arm64 libwayland-dev:arm64 libxkbcommon-dev:arm64
	@sudo apt-get install -y libgl-dev:arm64 libx11-dev:arm64 xorg-dev:arm64 libxxf86vm-dev:arm64 mingw-w64

# dev-deps for windows
.PHONY: .dev-deps-windows
.dev-deps-windows:
	@cd include/raylib-go/raylib/src && make PLATFORM=PLATFORM_DESKTOP

.PHONY: dev-deps
#: Install dependencies for compiling targets in this makefile
dev-deps: .deps
	@git submodule update --init --recursive
	@go mod tidy
	@$(MAKE) --no-print-directory .dev-deps-$(DEPS)
	@touch .dev-deps

.PHONY: test
#: Perfrom unit tests for application
test:
	@gofmt -e -w .
	@go test ./... -cover || true;

.PHONY: help
#: Lists available commands
help:
	@echo "Available Commands for project:"
	@grep -B1 -E "^[a-zA-Z0-9_\.-]+\:" Makefile | grep -A1 -E "^#\: \w+" \
	 | grep -v -- -- \
	 | sed 'N;s/\n/###/' \
	 | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	 | column -t  -s '###'

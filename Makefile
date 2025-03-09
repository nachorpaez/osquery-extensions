all: build

APP_NAME = osquery_extension
PACKAGE_VERSION= 0.0.2
PKGDIR_TMP = ${TMPDIR}golang

init:
	go mod init github.com/nachorpaez/osquery-extensions

.PHONY: deps
deps:
	go mod download
	go mod verify
	go mod vendor
	go mod tidy

.PHONY: .pre-build
.pre-build: clean
	mkdir -p build/darwin
	mkdir -p build/windows
	mkdir -p build/linux

clean:
	rm -rf build/
	rm -rf ${PKGDIR_TMP}_darwin

.PHONY: test
test:
	go test -v ./...

.PHONY:
apple: build/darwin/$(APP_NAME).ext
build/darwin/$(APP_NAME).ext: tables/**/*.go pkg/**/*.go main.go
	@mkdir -p $(@D)
	GOOS=darwin GOARCH=amd64 CGO_ENABLE=0 go build -o build/darwin/$(APP_NAME).amd64 -ldflags "-X main.packageVersion=$(PACKAGE_VERSION)"
	GOOS=darwin GOARCH=arm64 CGO_ENABLE=0 go build -o build/darwin/$(APP_NAME).arm64 -ldflags "-X main.packageVersion=$(PACKAGE_VERSION)"
	lipo -create -output build/darwin/$(APP_NAME).ext build/darwin/$(APP_NAME).amd64 build/darwin/$(APP_NAME).arm64

.PHONY:
win: build/windows/$(APP_NAME).ext.exe
build/windows/$(APP_NAME).ext.exe: tables/**/*.go pkg/**/*.go main.go
	@mkdir -p $(@D)
	GOOS=windows GOARCH=amd64 CGO_ENABLE=0 go build -o $(@) -ldflags "-X main.packageVersion=$(PACKAGE_VERSION)"

.PHONY:
osqueryi: build/darwin/$(APP_NAME).ext
	@echo build/darwin/$(APP_NAME).ext > build/darwin/extensions.load
	sleep 2
	sudo osqueryi --verbose --extensions_autoload=build/darwin/extensions.load --allow_unsafe
default:
	# go test -v
	go test -v ./compile
	# cd tests && go test -v
setup:
	# git clone https://github.com/evanw/esbuild
	cd cmd/zassets && go install

# static linking will prevent some problems with libsass when using different kernels
build:
	mkdir -p build
	cd cmd/zassets && GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-linkmode external -extldflags -static" -o build/zassets
build-windows:
	mkdir -p build
	cd cmd/zassets && GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-linkmode external -extldflags -static" -o build/zassets-windows
build-macos:
	mkdir -p build
	cd cmd/zassets && GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "-linkmode external -extldflags -static" -o build/zassets-macos

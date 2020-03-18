default:
	go test -v
	# go test -v ./compile
	# cd tests && go test -v
setup:
	# git clone https://github.com/evanw/esbuild
	cd cmd/zassets && go install

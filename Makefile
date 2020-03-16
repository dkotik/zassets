default:
	go test -v
	# cd tests && go test -v
setup:
	# git clone https://github.com/evanw/esbuild
	cd cmd/zassets && go install

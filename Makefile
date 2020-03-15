default:
	# go test -v
	# cd tests && go test -v
	# cd zassets && go test -v
	cd zassets/compile && go test -v
setup:
	git clone https://github.com/evanw/esbuild

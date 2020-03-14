module goresminpack

go 1.14

require (
	github.com/OneOfOne/xxhash v1.2.7
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/cobra v0.0.6
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/wellington/go-libsass v0.9.2
)

//replace ./minify=>../minify

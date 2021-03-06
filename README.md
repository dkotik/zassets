# Zassets, An Elegant Resource Bundler for Go

> **Disclaimers:** this project has been archived, because its main function has been replaced with Go embedding since version ~1.13.
>
> The resource debugger and compile APIs are unstable. The project carries some technical debt from some minifiers, which are not adapted to work well with virtual file systems, and a particular ESNext compiler ESBuild, which is incredibly fast but [quite hacky itself](https://github.com/evanw/esbuild/issues/13#issuecomment-587111778). Pull requests for more general usage and optimization are welcome.

The program generates static asset packs as embedded Zip archives. It allows for multiple asset packs to exist harmoniously side by side in the same module, unlike most other packers. Giving up individual gzip encoding, common in other resource packers, for a Zip bundle speeds up compilation for larger resource packs. It also allows easy loading of assets from an external archive by calling `zassets.FromArchive("zipFile.zip")` when needed. The loss of speed for directory transversal operations is justified by the assumption that static assets are rarely accessed directly. They should be shielded by a proper caching layer or mirror to a CDN. This project is inspired by the following excellent packages:

- https://github.com/shurcooL/vfsgen
- https://github.com/tdewolff/minify
- https://github.com/evanw/esbuild
- https://github.com/wellington/go-libsass

## Installation

Do not forget to install [esbuild](https://github.com/evanw/esbuild) if you need it for Javascript compilation.
There are build recipes for Linux, Windows, and MacOS in the Makefile. If you have the Go environment configured with $GOROOT/bin included in $PATH, enter the following command in your terminal:

```sh
go get -u github.com/spf13/cobra
go get -u github.com/dkotik/zassets
go install github.com/dkotik/zassets/cmd/zassets
```

## Usage

```
zassets [files or directories]... [flags]...

Flags:
  -c, --comment [lines]      Include a comment with the variable definition.
  -d, --debug                Readable refined output.
  -e, --embed                Embed provided files and directories.
  -i, --ignore [pattern]     Skip files and directories that match provided pattern.
                             You can use this flag multiple times.
  -o, --output [file]        Write program output to this location.
  -p, --package [name]       Assets will belong to this package.
  -r, --refine               Apply default refiners to assets before embedding.
  -s, --sum [algorithm]      Include a hash table sum.* in the embedded archive.
                             Choose from xxh64, md5, and sha256.
  -t, --tag [name]           Specify a build tag. You can use this flag multiple times.
  -v, --var [name]           Assets will be accessible using this variable name.
      --version              version for zassets
```

Point Zassets at resource files or directories by adding the following line under the package declaration in one of your Go files. When you run `go generate`, the program will create the _assets.gen.go_ file, which will contain the specified _Assets_ variable acting as _http.FileSystem_. To read an asset call _Assets.Open("assetpath.txt")_.

```go
//go:generate zassets dir1 dir2 file1 --output assets.gen.go --package mypackage --var Assets --embed
```

You may use shell redirection for the output in this manner: _//go:generate sh -c "zassets dir1 dir2 file1 --var Assets --package mypackage --embed > assets.gen.go"_. But if an error occurs due to a missing file or during refinement, the shell will truncate _assets.gen.go_ and cause package execution errors because of a missing variable. For this reason, it is wiser to specify _--output_ parameter explicitly every time instead of using shell redirection.

### Transformations

The following default file transformations happen when the `--refine` flag is activated. You may configure the compiler with custom refiners manually by using `zassets/compile` package.

- `*.js` files are compiled to ESNext `*.js` bundles.
- `*.sass` and `*.scss` files are compiled to `*.css`.
- `*.html`, `*.svg`, and `*.css` files are minified.
- `*.jpg`, `*.jpeg`, `*.png`, `*.webp` images are resized to 1080p and re-compressed for the Web.
- `*.min.*` files are served without any changes.

### Public Assets

If you would like to include a hash map for a given asset pack, use the `--sum [algorithm]` parameter. The resulting archive will contain a hash map file _sum.[algorithm]_, which you can use for validation or for serving assets by hash name. The second is handy for serving multiple versions of the same file in continual deployment setups or for a CDN or when multiple packages share some of the same assets. Files with identical content will be associated with the same hash name, when the _sum.[algorithm]_ file is included. Experimental helpers you may try:

- **`zassets.Public`**: An _http.FileSystem_ object that will open assets by an associated hash name. Serve files over HTTP from this object.
- **`zassets.PublicRegister`**: Link all assets from a given pack with _zassets.Public_. If the pack contains a supported _sum.[algorithm]_ file, the assets are associated using the hash map and their original file extension. Otherwise, the debug mode is assumed and a flimsy hash is generated from name space and path.
- **`zassets.PublicName`**: Returns the hash name associated with a resource. Use this function with your template engine to point to correct assets.
- **`zassets.Handler`**: An _http.Handler_ object that serves files from _zassets.Public_.

## Roadmap

- [ ] If an entry is missing, does not seem to raise an error. Packs all the previous entries as if everything is normal.
- [ ] Explain major benefit of multiple resource packs: a huge pile of icons takes a long time to zip up and they do not need to be re-compressed.
- [ ] Solidify Refiner API:
  - [ ] Javascript io.Pipe is clumsy. Needs re-writing.
  - [ ] Add image quality setting for the _Compiler_.
  - [ ] Add a refiner for Typescript.
  - [ ] Add a refiner for Clojure.
  - [ ] What if zassets points to the same folder "." - the pathing gets screwed up with a panic. Recursive inclusion bug possibility - should I fight it by filtering out .go files by default?!
  - [ ] (Default ignore setting?) If not set, ignore paths that begin with "." or "\_" or "node_modules".
- [ ] Solidify Debugger API:
  - [x] Watched directories do not respond to new files being created.
  - [ ] Currently copying files to a temporary directory. Building a virtual path tree would be faster.
- [ ] When the APIs are solid, hedge them with test suites.

## License

Zassets is released under the MIT license. The author would also like to add the SQLite blessing:

> May you do good and not evil. May you find forgiveness for yourself and forgive others. May you share freely, never taking more than you give.

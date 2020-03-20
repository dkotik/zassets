# Zassets, An Elegant Resource Bundler for Go
> **v0.0.1 Disclaimer:** The API is unstable. The project carries some technical debt from minifiers, which are not adapted to work well with virtual file systems, and a particular ESNext compiler ESBuild, which is incredibly fast but [quite hacky itself](https://github.com/evanw/esbuild/issues/13#issuecomment-587111778). Pull requests for more general usage and optimization are welcome.

The program generates an embedded static asset pack as a Zip archive next to a given directory. It allows for multiple asset packs to exist harmoniously side by side in the same module, unlike most other packers. This project is inspired by the following excellent packages:

- https://github.com/shurcooL/vfsgen
- https://github.com/tdewolff/minify
- https://github.com/markbates/pkger
- https://github.com/evanw/esbuild
- https://github.com/wellington/go-libsass

## Installation
Enter the following command in your terminal:
``` sh
go install github.com/dkotik/zassets/cmd/zassets
```

## Usage
Point Zassets at resource files or directories by adding the following line under the package declaration:

``` go
//go:generate zassets --output assets.gen.go --package test --var Assets --embed dir1 dir2 file1
```

<!-- _The program will create two files, the deployment file `assets.gen.go` and the development file `assets.dev.gen.go`. Each file contains an object that satisfies `http.FileSystem` interface named after the input directory. The development file can be activated using either **dev** or **debug** build tag and points directly to the transformed, but not minified, source files on disk, allowing you to edit them live._ -->

You may use shell redirection for the output in this manner: `//go:generate sh -c "zassets --var Assets --package test --embed dir1 dir2 file1 > assets.gen.go"`. But if an error occurs due to a missing file or during refinement, the shell will truncate `assets.gen.go` and cause package execution errors because of a missing variable. For this reason, it is wiser to specify `--output` parameter explicitly every time.

## Parameters
- `--hashwith`
- `--debug`
<!-- document all the parameters -->

## Transformations
- `*.sass` and `*.scss` files are compiled to `*.css`.
- `*.js` files are compiled to ESNext `*.js` bundles.
- `*.html`, `*.svg`, and `*.css` files are minified.
- `*.jpg`, `*.jpeg`, `*.png`, `*.webp` images are resized and re-compressed for the Web.
- `*.min.*` files are served without any changes.

## Public Assets
All files matching `/public/**` glob are registered to a content-based hash map. `goresminpack.PublicName` function returns the associated hash with the appropriate extension. It is handy for encoding asset paths in your template engine. Use `goresminpack.PublicHTTPHandler` to present all public files through a single handler.

## Roadmap
- [ ] // TODO: I need to capture error from that go func somehow
- [ ] Just use debug parameter? Hot-swapable <directory>.dev.gen.go driver that emulates serving of assets directly from disk, when launching in `debug` or `dev` build tags.
- [ ] add comments to all exported functions and classes
- [ ] add a blessing
- [ ] set logger for the compiler

// The directory operations are slow! Zip is not the right file
// format for frequent tree transversal. Put Store behind a
// a proper caching layer, if speed is a requirement.

there are build recipes in Makefile

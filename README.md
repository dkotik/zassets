# GoResMinPack, An Opinionated Resource Compiler

The program generates an embedded static asset pack next to a given directory. It is tailored completely for the convenience of its author at the moment, batteries included. **Warning:** It is hacky in many places with certain hard-wired defaults, a particular ESNext compiler ESBuild, which is [hacky itself](https://github.com/evanw/esbuild/issues/13#issuecomment-587111778), and the inclusion of [Bulma CSS][bulma] framework. The project will be streamlined later. Pull requests for more general usage and optimization are welcome. This project is indebted to the following excellent packages:

- github.com/shurcooL/vfsgen
- github.com/tdewolff/minify
- github.com/wellington/go-libsass
- github.com/evanw/esbuild

## Usage
Point GoResMinPack at a static resource directory by adding the following line under the package declaration:
``` go
//go:generate goresminpack assets
```
The program will create two files, the deployment file `assets.gen.go` and the development file `assets.dev.gen.go`. Each file contains an object that satisfies `http.FileSystem` interface named after the input directory. The development file can be activated using either **dev** or **debug** build tag and points directly to the transformed, but not minified, source files on disk, allowing you to edit them live.

## Transformations
- `*.sass` and `*.scss` files are compiled to `*.css`.
- `*.js` files are compiled to ESNext `*.js` bundles.
- `*.tmpl`, `*.html`, `*.svg`, `*.css`, and `*.sql` files are minified.
- `*.min.*` files are served without any changes.

## Public Assets
All files matching `/public/**` glob are registered to a content-based hash map. `goresminpack.PublicName` function returns the associated hash with the appropriate extension. It is handy for encoding asset paths in your template engine. Use `goresminpack.PublicHTTPHandler` to present all public files through a single handler.

## Roadmap
- [ ] Hot-swapable <directory>.dev.gen.go driver that emulates serving of assets directly from disk, when launching in `debug` or `dev` build tags.
- [ ] Webp image compression support.
- [ ] Check is `esbuild` is in Exec path.

## Debug Mode

[bulma]: https://github.com/jgthms/bulma

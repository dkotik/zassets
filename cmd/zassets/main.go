package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/alecthomas/kong"
	"github.com/dkotik/zassets"
	"github.com/dkotik/zassets/compile"
)

var cli struct {
	Entries  []string `kong:"arg"`
	Variable string   `kong:"flag,name='var',short='v',help='Assets will be accessible using this variable name.',default='Assets'"`
	// set default comment
	Comment       string   `kong:"flag,name='comment',short='c',help='Include a comment with the variable definition.',default='Assets ...'"`
	Tags          []string `kong:"flag,name='tag',short='t',help='Specify build tags.'"`
	Package       string   `kong:"flag,name='package',short='p',help='Assets will belong to this package.'"`
	Output        string   `kong:"flag,name='output',short='o',type='path',help='Write program output to this location.'"`
	Ignore        []string `kong:"flag,name='ignore',short='i',help='Skip files and directories that match provided patterns.'"`
	HashAlgorythm string   `kong:"flag,name='sum',short='s',help='Include a hash table sum.* in the embedded archive.',enum='xxh64,md5,sha256'"`
	Embed         bool     `kong:"flag,name='ember',short='e',help='Embed provided files and directories.'"`
	Refine        bool     `kong:"flag,name='refine',short='r',help='Apply default refiners to assets before embedding.'"`
	Debug         bool     `kong:"flag,name='debug',short='d',help='Readable refined output.',env='DEBUG'"`
}

func main() {
	err := func() error {
		c, err := kong.New(&cli,
			kong.Description(`An elegant resource bundler for Go. Compile and embed a resource directory into a Go module. Embedded resources are stored in an object	that satisfies http.FileSystem interface.`),
			kong.Vars{"version": "0.0.2"},
			kong.ConfigureHelp(kong.HelpOptions{
				Compact: true,
				Summary: true,
			}))
		if err != nil {
			return err
		}
		_, err = c.Parse(os.Args[1:])
		if err != nil {
			return err
		}
		em := &zassets.Embed{
			Comment:       cli.Comment,
			Debug:         cli.Debug,
			Entries:       cli.Entries,
			HashAlgorythm: cli.HashAlgorythm,
			Ignore:        cli.Ignore,
			Package:       cli.Package,
			Tags:          cli.Tags,
			Variable:      cli.Variable,
		}
		em.SetTemplate("") // TODO: clumsy here, should I change template live or debug?
		if len(os.Args) <= 1 {
			ctx, err := c.Parse([]string{`--help`})
			if err != nil {
				return err
			}
			return ctx.Run()
		}

		if em.Debug && cli.Embed {
			w := os.Stdout
			if cli.Output != "" {
				w, err = os.Create(cli.Output)
				if err != nil {
					return err
				}
				defer w.Close()
			}
			return em.Reader(w, nil)
		}

		iterator, err := compile.NewIterator(cli.Entries, cli.Ignore)
		if err != nil {
			return err
		}

		if cli.Refine && !cli.Embed {
			c, err := compile.NewCompiler(compile.WithDefaultOptions())
			if err != nil {
				return err
			}
			if em.Debug {
				compile.WithDebug()(c)
			}
			return c.Run(cli.Output, iterator)
		}

		if em.Refine {
			t, err := ioutil.TempDir(os.TempDir(), `zassets-*`)
			if err != nil {
				return err
			}
			c, err := compile.NewCompiler(compile.WithDefaultOptions())
			if err != nil {
				return err
			}
			if em.Debug {
				compile.WithDebug()(c)
			}
			err = c.Run(t, iterator)
			defer os.RemoveAll(t)
			if err != nil {
				return err
			}
			iterator = &compile.Iterator{
				Entries: []string{t},
				Ignore:  []*regexp.Regexp{}}
		}
		w := os.Stdout
		if cli.Output != "" {
			w, err = os.Create(cli.Output)
			if err != nil {
				return err
			}
			defer w.Close()
		}

		return em.Iterator(w, iterator)
	}()
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s.\n", err.Error())
	os.Exit(1)
}

package zassets

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/OneOfOne/xxhash"
	"github.com/dkotik/zassets/compile"
)

var defaultTemplate = `
{{- define "header" -}}
// Code generated by github.com/dkotik/zassets. DO NOT MODIFY.{{ with .Tags }}
// +build {{ . | tags }}{{ end }}

package {{ .Package }}

import "github.com/dkotik/zassets"

{{ if eq .Comment "" -}}
	// {{.Variable}} contains static assets.
{{- else -}}
	{{ comment .Comment }}
{{- end }}
var {{ .Variable }} = zassets.Must(zassets.FromBytes([]byte("
{{- end -}}

{{- define "footer" -}}
")))
{{- end -}}

{{- define "sum" -}}
// {{.Variable}}HashTable associates each entry with a content-based {{ .HashAlgorythm }} hash.
var {{ .Variable }}HashTable = map[string]string{
    {{- range $k, $v := .HashTable }}
    "{{ $k }}": "{{ printf "%x" $v }}",
    {{- end }}
}
{{- end -}}`

// Embed converts a binary stream or a path set to a Go asset Store.
type Embed struct {
	Variable      string
	Package       string
	Comment       string
	Tags          []string
	HashAlgorythm string
	HashTable     map[string][]byte

	template *template.Template
}

// SetTemplate sets up definition blocks for output rendering.
// If an empty string is given, the default template is used instead.
func (e *Embed) SetTemplate(t string) {
	e.template = template.New(`default`).Funcs(template.FuncMap{
		"tags": func(tags []string) string {
			return strings.Join(tags, `,`)
		},
		"quote": strconv.Quote,
		"comment": func(s string) string {
			result := strings.SplitAfter(s, "\n")
			return `// ` + strings.Join(result, `// `)
		},
	})
	if t == "" {
		e.template = template.Must(e.template.Parse(defaultTemplate))
		return
	}
	e.template = template.Must(e.template.Parse(t))
}

func (e *Embed) captureHash(w io.Writer, r io.Reader, p string) (err error) {
	var h hash.Hash
	switch e.HashAlgorythm {
	default: // do nothing
		_, err = io.Copy(w, r)
		return err
	case `md5`:
		h = md5.New()
	case `sha256`:
		h = sha256.New()
	case `xxhash`:
		h = xxhash.New64()
	}
	w = io.MultiWriter(h, w)
	_, err = io.Copy(w, r)
	e.HashTable[p] = h.Sum([]byte{})
	return err
}

// Iterator zips and embeds the contents of all paths.
func (e *Embed) Iterator(w io.Writer, i *compile.Iterator) (err error) {
	e.HashTable = make(map[string][]byte)
	pr, pw := io.Pipe() // TODO: this does not appear to be elegant
	// TODO: I need to capture error from that go func somehow
	go func() {
		defer pw.Close()
		z := zip.NewWriter(pw)
		defer z.Close()
		z.SetComment(`Resource pack generated by github.com/dkotik/zassets.`)
		err = i.Walk(func(target, relative string, info os.FileInfo) error {
			r, err := os.Open(target)
			if err != nil {
				return err
			}
			defer r.Close()
			h, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			h.Name = path.Clean(relative)
			h.Method = zip.Deflate
			w, err := z.CreateHeader(h)
			if err != nil {
				return err
			}
			return e.captureHash(w, r, relative)
		})

		if len(e.HashTable) > 0 {
			ws, err := z.CreateHeader(&zip.FileHeader{
				Name:    `sum.` + e.HashAlgorythm,
				Comment: `Table of hash values for each archived file.`,
			})
			if err != nil { // TODO: this does not appear elegant at all
				log.Fatal(err)
			}
			err = e.template.ExecuteTemplate(ws, `sum`, e)
			if err != nil { // TODO: this does not appear elegant at all
				log.Fatal(err)
			}
		}

		if err != nil { // TODO: this does not appear elegant at all
			log.Fatal(err)
		}
	}()
	return e.Reader(w, pr)
}

// Reader writes binary data through a template.
func (e *Embed) Reader(w io.Writer, r io.Reader) error {
	if e.Variable == "" {
		return errors.New("variable name must be specified")
	}
	if e.Package == "" {
		return errors.New("package name must be specified")
	}
	if e.HashAlgorythm != "" && e.HashAlgorythm != "md5" && e.HashAlgorythm != "sha256" && e.HashAlgorythm != "xx" {
		return errors.New("unknown hash algorythm, choose from xx, md5, sha")
	}

	err := e.template.ExecuteTemplate(w, `header`, e)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	var i, n int
	for {
		n, err = r.Read(buffer)
		if err != nil {
			break
		}
		for i = 0; i < n; i++ {
			_, err = fmt.Fprintf(w, `\x%02x`, buffer[i])
			if err != nil {
				return err
			}
		}
	}
	return e.template.ExecuteTemplate(w, `footer`, e)
}

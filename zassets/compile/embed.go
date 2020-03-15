package compile

import (
	"io"
	"strconv"
	"strings"
	"text/template"
)

// EmbedFunctions provides helpers for the file generation throug template.
var EmbedFunctions = template.FuncMap{
	"quote": strconv.Quote,
	"comment": func(s string) string {
		result := strings.SplitAfter(s, "\n")
		return `// ` + strings.Join(result, `// `)
	},
}
var defaultTemplate = template.Must(template.New(`default`).Funcs(EmbedFunctions).Parse(`package {{ .Package }}{{ with .Tags }}

// +build {{ . }}
{{end }}

// TODO: autogeneration sig.

import zassets{{ with .Comment }}

{{ comment . }}{{ end }}

var {{ .Name }} = zassets.Must(zassets.FromBytes([]byte("{{ range .Data }}{{ printf "\\x%02x" . }}{{ end }}")))
`))

// EmbedValues contains fields required for the template.
type EmbedValues struct {
	Name    string
	Package string
	Comment string
	Tags    string
	Data    []byte
}

// Embed writes encoded binary data through a template.
func Embed(w io.Writer, v *EmbedValues, t *template.Template) error {
	if t == nil {
		t = defaultTemplate
	}
	// Validate EmbedValues
	return t.Execute(w, v)
}

// demo implements a WASM module that can be used to format markdown
// with the goldmark-frontmatter extension.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

func main() {
	js.Global().Set("formatMarkdown",
		js.FuncOf(func(this js.Value, args []js.Value) any {
			var req request
			req.Decode(args[0])

			res, err := run(&req)
			if err != nil {
				res = &response{Error: err.Error()}
			}
			return res.Encode()
		}))
	select {}
}

type request struct {
	Markdown string
	YAML     bool
	TOML     bool
}

func (r *request) Decode(v js.Value) {
	r.Markdown = v.Get("markdown").String()
	r.YAML = v.Get("yaml").Bool()
	r.TOML = v.Get("toml").Bool()
}

type response struct {
	HTML        string
	Frontmatter string
	Error       string
}

func (r *response) Encode() js.Value {
	return js.ValueOf(map[string]any{
		"html":        r.HTML,
		"frontmatter": r.Frontmatter,
		"error":       r.Error,
	})
}

func run(req *request) (*response, error) {
	var formats []frontmatter.Format
	if req.YAML {
		formats = append(formats, frontmatter.YAML)
	}
	if req.TOML {
		formats = append(formats, frontmatter.TOML)
	}
	if len(formats) == 0 {
		// frontmatter.Extender will use a default set of formats
		// if none are specified.
		// We don't want the demo to do that.
		// Introduce a placeholder format to force an error.
		formats = []frontmatter.Format{
			{
				Name:  "placeholder",
				Delim: 0, // invalid delimiter
			},
		}
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Formats: formats,
			},
		),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err := md.Convert([]byte(req.Markdown), &buf, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("convert markdown: %w", err)
	}

	var fm string
	if data := frontmatter.Get(ctx); data != nil {
		var meta map[string]any
		if err := data.Decode(&meta); err != nil {
			return nil, fmt.Errorf("decode frontmatter: %w", err)
		}

		formatted, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("format frontmatter: %w", err)
		}

		fm = string(formatted)
	}

	return &response{
		HTML:        buf.String(),
		Frontmatter: fm,
	}, nil
}

package frontmatter

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Extender adds support for front matter to a Goldmark Markdown parser.
//
// Use it by installing it into the [goldmark.Markdown] object upon creation.
// For example:
//
//	goldmark.New(
//		// ...
//		goldmark.WithExtensions(
//			// ...
//			&frontmatter.Extender{},
//		),
//	)
type Extender struct {
	// Formats lists the front matter formats
	// that are supported by the extender.
	//
	// If empty, DefaultFormats is used.
	Formats []Format

	// TODO:
	// Bit map to opt into rendering/setting metadata?
}

var _ goldmark.Extender = (*Extender)(nil)

// Extend extends the provided Goldmark Markdown.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&Parser{
				Formats: e.Formats,
			}, 0),
		),
	)
}

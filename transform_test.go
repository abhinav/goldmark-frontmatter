package frontmatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

func TestMetaTransformer(t *testing.T) {
	t.Parallel()

	t.Run("empty context", func(t *testing.T) {
		t.Parallel()

		doc := ast.NewDocument()
		ctx := parser.NewContext()

		(&MetaTransformer{}).Transform(doc, nil, ctx)

		assert.Empty(t, doc.Meta())
	})

	t.Run("decode error", func(t *testing.T) {
		t.Parallel()

		doc := ast.NewDocument()
		ctx := parser.NewContext()

		(&Data{
			raw:    []byte("invalid"),
			format: YAML,
		}).set(ctx)

		(&MetaTransformer{}).Transform(doc, nil, ctx)

		assert.Empty(t, doc.Meta())
	})

	t.Run("", func(t *testing.T) {
		t.Parallel()

		doc := ast.NewDocument()
		ctx := parser.NewContext()

		(&Data{
			raw:    []byte("foo: bar"),
			format: YAML,
		}).set(ctx)

		(&MetaTransformer{}).Transform(doc, nil, ctx)

		assert.Equal(t, map[string]any{
			"foo": "bar",
		}, doc.Meta())
	})
}

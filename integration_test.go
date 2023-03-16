package frontmatter_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"gopkg.in/yaml.v3"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/integration.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc string         `yaml:"desc"`
		Give string         `yaml:"give"`
		Want string         `yaml:"want"`
		Data map[string]any `yaml:"data"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			md := goldmark.New(
				goldmark.WithExtensions(&frontmatter.Extender{
					Mode: frontmatter.SetMetadata,
				}),
			)

			src := []byte(tt.Give)

			t.Run("Data.Decode", func(t *testing.T) {
				t.Parallel()

				ctx := parser.NewContext()
				var got bytes.Buffer
				require.NoError(t,
					md.Convert(src, &got, parser.WithContext(ctx)))
				assert.Equal(t,
					strings.TrimSuffix(tt.Want, "\n"),
					strings.TrimSuffix(got.String(), "\n"),
				)

				meta := frontmatter.Get(ctx)
				require.NotNil(t, meta)

				var gotData map[string]any
				require.NoError(t, meta.Decode(&gotData))
				assert.Equal(t, tt.Data, gotData)
			})

			t.Run("Document.Meta", func(t *testing.T) {
				t.Parallel()

				doc := md.Parser().Parse(text.NewReader(src)).OwnerDocument()
				assert.Equal(t, tt.Data, doc.Meta())
			})
		})
	}
}

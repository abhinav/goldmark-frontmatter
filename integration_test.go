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
				goldmark.WithExtensions(&frontmatter.Extender{}),
			)

			ctx := parser.NewContext()
			var got bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &got, parser.WithContext(ctx)))
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
	}
}

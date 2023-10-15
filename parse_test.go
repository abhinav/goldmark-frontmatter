package frontmatter

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"gopkg.in/yaml.v3"
)

func TestParser(t *testing.T) {
	t.Parallel()

	testdata, err := os.ReadFile("testdata/parser.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc    string `yaml:"desc"`
		Give    string `yaml:"give"`
		Formats []struct {
			Name  string `yaml:"name"`
			Delim string `yaml:"delim"`
		} `yaml:"formats"`

		WantFormat string `yaml:"wantFormat"`
		WantRaw    string `yaml:"wantRaw"`
	}
	require.NoError(t, yaml.Unmarshal(testdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			var formats []Format
			for _, f := range tt.Formats {
				require.Len(t, f.Delim, 1,
					"bad format %q: delim must be a single character", f.Name)
				delim := []byte(f.Delim)[0]
				formats = append(formats, Format{
					Name:  f.Name,
					Delim: delim,
					Unmarshal: func(data []byte, v interface{}) error {
						t.Fatalf("unexpected call to %v.Unmarshal", f.Name)
						return errors.New("unreachable")
					},
				})
			}

			rdr := text.NewReader([]byte(tt.Give))
			ctx := parser.NewContext()
			p := parser.NewParser(
				parser.WithBlockParsers(
					util.Prioritized(&Parser{
						Formats: formats,
					}, 500),
				),
			)

			p.Parse(rdr, parser.WithContext(ctx))

			data := Get(ctx)
			if tt.WantFormat == "" {
				assert.Nil(t, data)
				return
			}

			assert.Equal(t, tt.WantFormat, data.format.Name)
			assert.Equal(t, tt.WantRaw, string(data.raw))
		})
	}
}

func TestFrontmatterNode_Dump(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	stdoutr, stdoutw, err := os.Pipe()
	require.NoError(t, err)

	src := []byte("title: Hello World")
	node := frontmatterNode{
		Format:  YAML,
		Segment: text.NewSegment(0, len(src)),
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		os.Stdout = stdoutw
		node.Dump(src, 0)
		os.Stdout = oldStdout
		assert.NoError(t, stdoutw.Close())
	}()
	<-done

	stdout, err := io.ReadAll(stdoutr)
	require.NoError(t, err)
	assert.NoError(t, stdoutr.Close())

	assert.Contains(t, string(stdout), "Format: YAML\n")
	assert.Contains(t, string(stdout), "Data: title: Hello World\n")
}

func TestLineDelim(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		give  string
		want  byte
		count int
	}{
		{desc: "empty"},
		{desc: "no delim", give: "foo"},
		{desc: "too short", give: "--"},
		{desc: "mismatch", give: "---_-"},
		{
			desc:  "trailing crlf",
			give:  "---\r\n",
			want:  '-',
			count: 3,
		},
		{
			desc:  "trailing lf",
			give:  "---\n",
			want:  '-',
			count: 3,
		},
		{
			desc:  "three",
			give:  "---",
			want:  '-',
			count: 3,
		},
		{
			desc:  "many",
			give:  "-------",
			want:  '-',
			count: 7,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			delim, count := lineDelim([]byte(tt.give))
			assert.Equal(t, tt.want, delim)
			assert.Equal(t, tt.count, count)
		})
	}
}

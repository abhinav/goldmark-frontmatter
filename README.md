# goldmark-frontmatter

[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/goldmark/frontmatter.svg)](https://pkg.go.dev/go.abhg.dev/goldmark/frontmatter)
[![CI](https://github.com/abhinav/goldmark-frontmatter/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/goldmark-frontmatter/actions/workflows/ci.yml)

goldmark-frontmatter is an extension for the [goldmark] Markdown parser
that adds support for parsing YAML and TOML front matter from Markdown documents.

  [goldmark]: http://github.com/yuin/goldmark

## Features

- Parses YAML and TOML front matter out of the box
- Allows defining your own front matter formats
- Exposes front matter in via a types-safe API

### Demo

A web-based demonstration of the extension is available at
<https://abhinav.github.io/goldmark-frontmatter/demo/>.

## Installation

```bash
go get go.abhg.dev/goldmark/frontmatter@latest
```

## Usage

To use goldmark-frontmatter, import the `frontmatter` package.

```go
import "go.abhg.dev/goldmark/frontmatter"
```

Then include the `frontmatter.Extender` in the list of extensions
when constructing your [`goldmark.Markdown`].

  [`goldmark.Markdown`]: https://pkg.go.dev/github.com/yuin/goldmark#Markdown

```go
goldmark.New(
  goldmark.WithExtensions(
    // ...
    &frontmatter.Extender{},
  ),
).Convert(src, out)
```

By default, this won't have any effect except stripping the front matter
from the document.
See [Accessing front matter](#accessing-front-matter) on how to read it.

### Syntax

Front matter starts with three or more instances of a delimiter,
and must be the first line of the document.

The supported delimiters are:

- YAML: `-`

    For example:

    ```
    ---
    title: goldmark-frontmatter
    tags: [markdown, goldmark]
    description: |
      Adds support for parsing YAML front matter.
    ---

    # Heading 1
    ```

- TOML: `+`

    For example:

    ```
    +++
    title = "goldmark-frontmatter"
    tags = ["markdown", "goldmark"]
    description = """\
      Adds support for parsing YAML front matter.\
      """
    +++

    # Heading 1
    ```

The front matter block ends with the same number of instances of the delimiter.
So if the opening line used 10 occurrences, so must the closing.

    ---------------------------
    title: goldmark-frontmatter
    tags: [markdown, goldmark]
    ---------------------------

### Accessing front matter

To access the front matter parsed by goldmark-frontmatter,
you must pass in a `parser.Context`
when you call `Markdown.Convert` or `Parser.Parse`.

```go
md := goldmark.New(
  goldmark.WithExtensions(&frontmatter.Extender{}),
  // ...
)

ctx := parser.NewContext()
md.Convert(src, out, parser.WithContext(ctx))

d := frontmatter.Get(ctx)
```

From here, you can decode the front matter into a struct.

```go
var meta struct {
  Title string   `yaml:"title"`
  Tags  []string `yaml:"tags"`
  Desc  string   `yaml:"description"`
}
if err := d.Decode(&meta); err != nil {
  // ...
}
```

#### Decoding all fields

Decode into a `map[string]any` if you need access to everything
in the front matter.

```go
var meta map[string]any
if err := fm.Decode(&meta); err != nil {
  // ...
}
```

## Similar projects

- [yuin/goldmark-meta](https://github.com/yuin/goldmark-meta)
  provides support for YAML front matter.

## License

This software is made available under the MIT license.

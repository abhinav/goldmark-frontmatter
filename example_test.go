package frontmatter_test

import (
	"fmt"
	"log"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

func ExampleData_Decode() {
	const src = `---
title: Foo
tags: [bar, baz]
---

This page is about foo.
`

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)

	ctx := parser.NewContext()
	if err := md.Convert([]byte(src), os.Stdout, parser.WithContext(ctx)); err != nil {
		log.Fatal(err)
	}

	fm := frontmatter.Get(ctx)
	if fm == nil {
		log.Fatal("no frontmatter found")
	}

	var data struct {
		Title string   `yaml:"title"`
		Tags  []string `yaml:"tags"`
	}
	if err := fm.Decode(&data); err != nil {
		log.Fatal(err)
	}

	fmt.Println("---")
	fmt.Println("Title: ", data.Title)
	fmt.Println("Tags: ", data.Tags)

	// Output:
	// <p>This page is about foo.</p>
	// ---
	// Title:  Foo
	// Tags:  [bar baz]
}

func ExampleMode_setMetadata() {
	const src = `+++
title = "Foo"
tags = ["bar", "baz"]
+++

This page is about foo.
`

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Mode: frontmatter.SetMetadata,
			},
		),
	)

	doc := md.Parser().Parse(text.NewReader([]byte(src)))
	meta := doc.OwnerDocument().Meta()

	fmt.Println("Title: ", meta["title"])
	fmt.Println("Tags: ", meta["tags"])

	// Output:
	// Title:  Foo
	// Tags:  [bar baz]
}

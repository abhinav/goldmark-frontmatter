module go.abhg.dev/goldmark/frontmatter/demo

go 1.23.0

toolchain go1.24.5

replace go.abhg.dev/goldmark/frontmatter => ../

require (
	github.com/yuin/goldmark v1.7.12
	go.abhg.dev/goldmark/frontmatter v0.2.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

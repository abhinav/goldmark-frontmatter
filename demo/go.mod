module go.abhg.dev/goldmark/frontmatter/demo

go 1.24.0

toolchain go1.25.5

replace go.abhg.dev/goldmark/frontmatter => ../

require (
	github.com/yuin/goldmark v1.7.16
	go.abhg.dev/goldmark/frontmatter v0.3.0
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

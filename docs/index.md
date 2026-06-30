# go-ruby-rexml documentation

**Ruby's REXML XML library in pure Go — DOM, parser, serialisers, XPath subset; MRI byte-exact, no cgo.**

`go-ruby-rexml/rexml` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's
[`REXML`](https://docs.ruby-lang.org/en/master/REXML.html), matching reference Ruby (MRI) byte-for-byte. The module
path is `github.com/go-ruby-rexml/rexml`.

It was **extracted from rbgo's internals into a reusable standalone library**:
the module is standalone and importable by any Go program, and it is the backend
bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo`
as a native module — a sibling of [go-ruby-yaml](https://github.com/go-ruby-yaml/yaml) (Psych), [go-ruby-regexp](https://github.com/go-ruby-regexp/regexp) (Onigmo) and [go-ruby-erb](https://github.com/go-ruby-erb/erb) (ERB). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete — MRI byte-exact"
    A wide XML corpus parsed and re-serialised here and by the system `ruby` (`REXML::Document.new(xml).to_s`, `Formatters::Pretty`, `REXML::XPath.match`), compared exactly via stdin/stdout binmode; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches.

## Quick taste

```go
d, _ := rexml.ParseDocument(`<config><server name="web" port="8080"/></config>`)

fmt.Println(d.ToString())
// <config><server name='web' port='8080'/></config>  (attrs sorted, single-quoted)

fmt.Println(d.Pretty(2))
// <config>
//   <server name='web' port='8080'/>
// </config>

srv := rexml.XPathFirst(d, "//server").(*rexml.Element)
name, _ := srv.Attr("name")
fmt.Println(name) // web

// Build a tree programmatically.
doc := rexml.NewDocument()
root := doc.AddElement("root")
root.AddElement("item").AddText("hello & <world>")
fmt.Println(doc.ToString())
// <root><item>hello &amp; &lt;world&gt;</item></root>
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`rexml`](https://github.com/go-ruby-rexml/rexml) | the library — Ruby's REXML in pure Go |
| [`docs`](https://github.com/go-ruby-rexml/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-rexml.github.io`](https://github.com/go-ruby-rexml/go-ruby-rexml.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-rexml/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-rexml/rexml](https://github.com/go-ruby-rexml/rexml).

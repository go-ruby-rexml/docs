# Usage & API

The public API lives at the module root (`github.com/go-ruby-rexml/rexml`). It is **Ruby-shaped but Go-idiomatic**: the node types mirror `REXML::Document` / `Element` / `Attribute` / `Text` / …, while the surface follows Go conventions — value types, explicit `error`, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-rexml/rexml`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-rexml/rexml
```

## Worked example

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

## Shape

```go
// Parse & construct
func ParseDocument(s string) (*Document, error)  // REXML::Document.new
func NewDocument() *Document

// Serialise
func (d *Document) ToString() string             // to_s (compact)
func (d *Document) Pretty(indent int) string     // Formatters::Pretty
func (d *Document) PrettyString() string
func (d *Document) Write(o WriteOptions) string

// DOM
func (e *Element) AddElement(name string) *Element
func (e *Element) AddAttribute(name, value string)
func (e *Element) AddText(s string) *Element
func (e *Element) Elements(path string) []*Element
func (e *Element) Attr(name string) (string, bool)
func (e *Element) Text() string
func (d *Document) Root() *Element

// XPath subset
func XPathFirst(ctx Node, path string) Node
func XPathEach(ctx Node, path string, fn func(Node))
func XPathMatch(ctx Node, path string) []Node
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-rexml/rexml` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-yaml](https://github.com/go-ruby-yaml/yaml) (Psych), [go-ruby-regexp](https://github.com/go-ruby-regexp/regexp) (Onigmo) and [go-ruby-erb](https://github.com/go-ruby-erb/erb) (ERB) are bound. The dependency runs the
other way: this library has no dependency on the Ruby runtime.

// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ruby-rexml/rexml"
)

// doc is a fixed, reproducible XML document: nested elements, attributes, mixed
// whitespace text, and a namespace (xmlns:meta / meta:rev). It mirrors
// benchmarks/ruby/rexml.rb byte-for-byte so the two drivers parse identical
// input and their outputs can be diffed before any timing is trusted.
const doc = `<?xml version='1.0' encoding='UTF-8'?>
<catalog xmlns:meta='http://example.com/meta' id='c1'>
  <book id='b1' meta:rev='2'>
    <title>Go in Action</title>
    <author>Kennedy</author>
    <price currency='USD'>39.99</price>
  </book>
  <book id='b2'>
    <title>The Go Programming Language</title>
    <author>Donovan</author>
    <price currency='USD'>44.50</price>
  </book>
</catalog>`

// mustParse parses the fixed document, panicking on error (the input is a
// compile-time constant, so a failure is a bug in the library, not the data).
func mustParse() *rexml.Document {
	d, err := rexml.Parse(doc)
	if err != nil {
		panic(err)
	}
	return d
}

// titleTexts returns the text of every <title> reachable by the descendant
// query //title, in document order.
func titleTexts(d *rexml.Document) []string {
	var out []string
	for _, n := range rexml.XPathMatch(d, "//title") {
		if e, ok := n.(*rexml.Element); ok {
			out = append(out, e.Text())
		}
	}
	return out
}

// verify prints the canonical output of every benchmarked operation so it can be
// diffed against the Ruby side before any timing is trusted.
func verify() {
	d := mustParse()
	fmt.Println("== parse ==")
	fmt.Printf("root=%s\n", d.Root().QName())
	fmt.Printf("books=%d\n", len(rexml.XPathMatch(d, "//book")))
	fmt.Println("== xpath //title ==")
	titles := titleTexts(d)
	fmt.Printf("%d\n", len(titles))
	for _, t := range titles {
		fmt.Println(t)
	}
	fmt.Println("== xpath /catalog/book ==")
	fmt.Printf("%d\n", len(rexml.XPathMatch(d, "/catalog/book")))
	fmt.Println("== serialize (to_s) ==")
	fmt.Println(d.ToString())
	fmt.Println("== pretty (indent 2) ==")
	fmt.Println(strings.TrimRight(d.Pretty(2), "\n"))
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "verify" {
		verify()
		return
	}
	d := mustParse()
	bench("parse", 500, func() { v, _ := rexml.Parse(doc); sink = v })
	bench("xpath-descendant", 2000, func() { sink = rexml.XPathMatch(d, "//title") })
	bench("xpath-absolute", 2000, func() { sink = rexml.XPathMatch(d, "/catalog/book") })
	bench("serialize", 1000, func() { sink = d.ToString() })
	bench("pretty", 1000, func() { sink = d.Pretty(2) })
}

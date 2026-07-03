# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "rexml/document"
require "rexml/formatters/pretty"
require_relative "_harness"

# Fixed, reproducible XML document — mirrors benchmarks/go/main.go byte-for-byte:
# nested elements, attributes, mixed whitespace text, and a namespace.
DOC = <<~XML.chomp
  <?xml version='1.0' encoding='UTF-8'?>
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
  </catalog>
XML

def pretty(doc)
  out = +""
  REXML::Formatters::Pretty.new(2).write(doc, out)
  out
end

if ARGV[0] == "verify"
  d = REXML::Document.new(DOC)
  puts "== parse =="
  puts "root=#{d.root.name}"
  puts "books=#{REXML::XPath.match(d, '//book').size}"
  puts "== xpath //title =="
  titles = REXML::XPath.match(d, "//title").map(&:text)
  puts titles.size
  titles.each { |t| puts t }
  puts "== xpath /catalog/book =="
  puts REXML::XPath.match(d, "/catalog/book").size
  puts "== serialize (to_s) =="
  puts d.to_s
  puts "== pretty (indent 2) =="
  puts pretty(d).chomp
  exit
end

D = REXML::Document.new(DOC)
bench("parse",            500) { REXML::Document.new(DOC) }
bench("xpath-descendant", 2000) { REXML::XPath.match(D, "//title") }
bench("xpath-absolute",   2000) { REXML::XPath.match(D, "/catalog/book") }
bench("serialize",        1000) { D.to_s }
bench("pretty",           1000) { pretty(D) }

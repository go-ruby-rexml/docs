<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-rexml` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-rexml` library**
against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby, TruffleRuby). It
measures the **library primitive** through its Go API, isolated from the rbgo
interpreter, so the numbers answer: *is the pure-Go implementation as fast as the
reference runtime's own `rexml`?*

## Layout

- `go/`            — self-contained Go driver; `go.mod` pins the published
  library by pseudo-version (no `replace`). The built `bench` binary is
  git-ignored.
- `ruby/rexml.rb`  — the equivalent workload; `ruby/_harness.rb` is the shared
  timer.
- `run.sh`         — runs every available runtime and prints one Markdown table
  per sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Operations

The **same fixed `<catalog>` document** (nested elements, attributes, mixed
whitespace text, and an `xmlns:meta` namespace) is driven through five
operations: `Document.new` (parse), a descendant XPath query (`//title`), an
absolute XPath query (`/catalog/book`), the compact serialiser (`Document#to_s`),
and `Formatters::Pretty` at indent 2. The two XPath expressions stay inside the
library's supported REXML subset (`child`/`descendant` steps, positional and
`@attr` predicates, `@attr` / `text()`).

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script parse **identical input bytes**
and each exposes a `verify` mode whose canonical output (root name, XPath match
counts and title texts, and the full `to_s` / `Pretty(2)` serialisations) is
checked identical to MRI before timing:

```sh
(cd benchmarks/go && go run . verify) > /tmp/go.txt
ruby benchmarks/ruby/rexml.rb verify  > /tmp/rb.txt
diff /tmp/go.txt /tmp/rb.txt   # must be empty
```

Results are published, dated, in [`../docs/performance.md`](../docs/performance.md).

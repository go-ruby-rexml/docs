# Performance

`go-ruby-rexml/rexml` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `rexml` — a
tree model, an XML parser, the two serialisers (the default compact formatter
and `Formatters::Pretty`), and a widely-used subset of XPath, all with no cgo.
This page records a **comparative, library-level benchmark** of that module
against the reference Ruby runtimes, part of the ecosystem-wide per-module parity
suite.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This measures the **pure-Go library directly, through its Go API**, isolated from
any interpreter dispatch, answering the parity question head-on: *is the pure-Go
implementation as fast as the reference runtime's own `rexml`?* The **same
workload, same fixed input document, same iteration counts** run through the Go
library and through each reference runtime's `rexml` stdlib; every operation's
output was checked **byte-identical to MRI** before any timing (see *Reproduce*
below).

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
  All runtimes measured **on the host**, no VM.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` (the oracle) · MRI + YJIT ·
  JRuby 10.1.0.0 (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Fixed input (reproducible, identical bytes in both drivers):** a small
  `<catalog>` document — an XML declaration, two nested `<book>` elements each
  with `id` / `currency` attributes and `<title>`/`<author>`/`<price>` text
  children, a `meta:` namespace (`xmlns:meta` + `meta:rev`), and the whitespace
  between elements.
- **Operations:** `Document.new` (parse), a descendant XPath query (`//title`),
  an absolute XPath query (`/catalog/book`), the default compact serialiser
  (`Document#to_s`), and `Formatters::Pretty` at indent 2.
- **XPath scope:** the library implements the widely-used REXML subset —
  `child`/`descendant` steps, `[n]` positional, `[@attr]` / `[@attr='v']`
  predicates, `@attr` and `text()`. The two queries benchmarked (`//title`,
  `/catalog/book`) sit inside that subset; axes such as `ancestor` /
  `following-sibling` and functions like `count()` / `position()` are out of
  scope and are not benchmarked.
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### parse

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1883.4 | 0.02× |
| MRI | 95424.0 | 1.00× |
| MRI + YJIT | 47020.0 | 0.49× |
| JRuby | 155751.1 | 1.63× |
| TruffleRuby | 73867.1 | 0.77× |

#### xpath-descendant

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 303.4 | 0.01× |
| MRI | 47122.0 | 1.00× |
| MRI + YJIT | 29298.0 | 0.62× |
| JRuby | 49504.4 | 1.05× |
| TruffleRuby | 21546.3 | 0.46× |

#### xpath-absolute

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 157.0 | 0.01× |
| MRI | 18170.5 | 1.00× |
| MRI + YJIT | 13004.5 | 0.72× |
| JRuby | 19157.4 | 1.05× |
| TruffleRuby | 14717.8 | 0.81× |

#### serialize

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 771.1 | 0.05× |
| MRI | 14249.0 | 1.00× |
| MRI + YJIT | 7325.0 | 0.51× |
| JRuby | 11631.8 | 0.82× |
| TruffleRuby | 11302.4 | 0.79× |

#### pretty

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1053.9 | 0.06× |
| MRI | 17031.0 | 1.00× |
| MRI + YJIT | 11192.0 | 0.66× |
| JRuby | 8844.0 | 0.52× |
| TruffleRuby | 17223.5 | 1.01× |

REXML is a **pure-Ruby** XML library — parser, XPath engine, and serialisers are
all interpreted Ruby — so a compiled pure-Go implementation dominates every
operation. Parsing is **~51× faster than MRI** (0.02×): MRI's `Document.new`
tokenises and builds the tree in Ruby, while the Go library scans bytes directly.
The XPath queries are the widest gaps — **~100–150× faster** (both 0.01×) —
because REXML's XPath compiles and walks an expression tree in Ruby per call,
whereas the Go subset resolves the path against the node set directly.
Serialisation is **~18× faster** (0.05×) and the Pretty formatter **~16×**
(0.06×). YJIT roughly halves MRI on parse and serialise (0.49×–0.66×) but cannot
close an order-of-magnitude interpreter gap. JRuby's warm-up shows on `parse`,
where it lands slower than cold MRI (1.63×) inside this fixed pass budget.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-rexml/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`
    pseudo-version — no `replace`), the equivalent `ruby/rexml.rb` workload, and
    `run.sh`. Run `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass
    budget and `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries. Both
    drivers accept a `verify` argument (`(cd benchmarks/go && go run . verify)`
    vs `ruby benchmarks/ruby/rexml.rb verify`) that prints each operation's
    canonical output — the parsed root name, the `//book` and `//title` /
    `/catalog/book` match counts and title texts, and the full `to_s` and
    `Pretty(2)` serialisations; the two were confirmed **byte-identical** before
    timing.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly JRuby on `parse`, a cold-JIT figure, not
    steady-state. Sub-microsecond rows (the two XPath queries) carry the most
    relative noise; treat those ratios as order-of-magnitude. Every number here
    is a **real measured value** from the dated run above — nothing is
    fabricated, estimated, or cherry-picked. The go-ruby column is the pure-Go
    library; every other column is that interpreter's own `rexml` stdlib doing
    the equivalent work.

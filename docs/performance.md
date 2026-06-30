# Performance

`go-ruby-rexml/rexml` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `rexml`. This
page records the **methodology** of the ecosystem-wide per-module parity suite —
how this module is benchmarked against the reference Ruby runtimes — without
quoting numbers here, so the figures never drift out of date.

## What is measured

The **same** Ruby script — an XML parse + re-serialise round-trip (`REXML::Document.new(xml).to_s`, `Formatters::Pretty`, `REXML::XPath.match`) — is run under every runtime.
`rbgo`'s number reflects **this pure-Go library doing the work**; every other
column is that interpreter's own `rexml` (or equivalent) implementation. So the
comparison is the **Ruby-visible operation**, apples-to-apples across
interpreters. The script prints a deterministic checksum and its output is
checked **byte-identical to MRI** before timing.

## Method

- **Best-of-N wall time** (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** MRI (the oracle) and MRI + YJIT; JRuby (OpenJDK); TruffleRuby
  (GraalVM CE Native). JRuby and TruffleRuby are timed **cold, single-shot**, so
  they carry JVM / Graal startup on every run — read them as one-shot
  `ruby file.rb` costs, the same way `rbgo` and MRI are measured, not as
  steady-state JIT numbers.
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`rexml.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-rexml) | 20 | 0.04× |
| MRI (ruby 4.0.5) | 510 | 1.00× |
| MRI + YJIT | 310 | 0.61× |
| JRuby 10.1.0.0 | 2850 | 5.59× |
| TruffleRuby 34.0.1 | 1490 | 2.92× |

rbgo runs on **go-ruby-rexml** and is **~25x faster than MRI** here (0.04x): MRI's REXML is a pure-Ruby XML parser/serializer, so the compiled pure-Go library dominates the parse+serialize loop. Second-biggest win of the wave-3 suite.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.

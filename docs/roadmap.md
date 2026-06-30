# Roadmap

`go-ruby-rexml/rexml` is grown **test-first**, each capability differential-tested against
MRI rather than built in isolation. Ruby's REXML — the deterministic,
interpreter-independent slice extracted from rbgo's internals — is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Parse | `ParseDocument` building the full node tree, preserving raw text / entities / XML declaration / DOCTYPE, returning a `*ParseError` with byte offset on malformed input. | **Done** |
| DOM | `Element` accessors and mutators (`Add*`, `Elements`, `EachElement`, `Text`, `ElementAt`, `QName` / `NamespaceURI`, `Root`), with an ordered `Attributes` map. | **Done** |
| Serialise | Compact `ToString`, `Write(WriteOptions)`, and the `Pretty` formatter — exact attribute quoting / ordering, entity escaping, CDATA / comment / PI round-trip, self-closing empties. | **Done** |
| XPath subset | `XPathFirst` / `XPathEach` / `XPathMatch` over the widely-used subset; axes / functions outside the documented boundary are out of scope. | **Done** |
| Node model | The concrete `Document` / `Element` / `Attribute` / `Text` / `Comment` / `CData` / `Instruction` / `DocType` / `XMLDecl` types a host binds, with raw round-trip fidelity. | **Done** |
| Differential oracle & coverage | A wide XML corpus re-serialised here and by the system `ruby`, compared byte-for-byte; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **XPath subset boundary.** Only the widely-used subset is implemented: axes beyond child / descendant (no `ancestor`, `following-sibling`, `parent`, `..`), functions beyond `text()` (no `count()`, `position()`, `last()`, `name()`, `contains()`, …), arithmetic, and compound boolean predicates (`and` / `or`) are out of scope; an unrecognised predicate leaves the node set unchanged.
- **No interpreter.** The library implements the deterministic parse / serialise / XPath algorithm; it never runs arbitrary Ruby. Anything needing a live binding is the consumer's job — that is why `rbgo` binds this module rather than the reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's `rexml` 3.4.4, pinned by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime; the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.

# SpiceDB Schema Parser & Validator CLI Tool

spicedb-dsl-validator is a CLI tool which parse and validate SpiceDB [schema](https://authzed.com/docs/guides/schema).

## About SpiceDB
[SpiceDB](https://authzed.com) is an open source, Google Zanzibar-inspired database for creating and managing security-critical application permissions.

# Build
`make binary`

# Usage

Parsed Correctly
```agsl
./spicedb-dsl-validator parse --file-path tests/empty.zed 
```

Broken schema
```agsl
./spicedb-dsl-validator parse --file-path tests/broken.zed
E1220 22:57:11.609797    4216 parse.go:135] Complied error: parse error in `schema`, line 1, column 1: Expected end of statement or definition, found: TokenTypeError

```
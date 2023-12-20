# Spicedb Schema Parser CLI

A simple cli utility to verification the Spicedb schema DSL can be parsed or not.

# Build
`make binary`

# Usage

Parsed Correctly
```agsl
./spicedb-dsl-validator parse --file-path tests/empty.zed 
Parsed correctly

```

Broken schema
```agsl
./spicedb-dsl-validator parse --file-path tests/broken.zed
Extracted Error Message: Expected end of statement or definition, found: TokenTypeError

```
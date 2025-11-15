# TOON Specification Reference

This Go package implements the TOON specification as defined in the official spec repository.

> **Fork Notice**: This is a community fork providing a Go implementation. It is not affiliated with or maintained by the original TOON creators.

## Specification Source

**Authoritative Specification**: [github.com/toon-format/spec](https://github.com/toon-format/spec)

**Current Version**: 1.4 (2025-11-05)

## Quick Links

- **[Full Specification](https://github.com/toon-format/spec/blob/main/SPEC.md)** - Complete technical specification
- **[Changelog](https://github.com/toon-format/spec/blob/main/CHANGELOG.md)** - Version history
- **[Examples](https://github.com/toon-format/spec/tree/main/examples)** - Example TOON files
- **[Conformance Tests](https://github.com/toon-format/spec/tree/main/tests)** - Language-agnostic test fixtures

## This Repository

This repository (`github.com/soy4rias/toongo`) is a **community fork** that provides a Go implementation of the TOON specification. It achieves 100% conformance with the official specification tests.

### Original Implementations

- **Official TypeScript Implementation**: [toon-format/toon](https://github.com/toon-format/toon) - The reference implementation by the original creators
- **Official Specification**: [toon-format/spec](https://github.com/toon-format/spec) - The authoritative source for TOON format definition

### This Implementation

- **Go Package**: `github.com/soy4rias/toongo`
- **Conformance**: 100% (323/323 tests passing)
- **Status**: Complete and production ready
- **Maintainer**: Community maintained (not affiliated with original creators)

## Conformance

This Go implementation has been tested against all official conformance tests:

- Decode: 185/185 tests passing (100%)
- Encode: 138/138 tests passing (100%)
- Total: 323/323 tests passing (100%)

See [CONFORMANCE.md](./CONFORMANCE.md) for detailed results.

## Reporting Issues

- **Specification Issues**: Report to [toon-format/spec](https://github.com/toon-format/spec/issues)
- **Go Implementation Issues**: Report to this repository's issues
- **TypeScript Implementation Issues**: Report to [toon-format/toon](https://github.com/toon-format/toon/issues)

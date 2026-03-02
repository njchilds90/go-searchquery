# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - Initial Release
### Added
- Zero-dependency Lexer and Parser for GitHub-style string queries.
- Support for extracting key-value pairs (`key:value`).
- Support for parsing quoted phrases (`"some phrase"`).
- Support for exclusion modifiers (`-key:value` or `-"excluded phrase"`).
- Deterministic `String()` reconstruction methods to turn ASTs back into secure query strings.
- GoDoc documentation and table-driven unit tests.

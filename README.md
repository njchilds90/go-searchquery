# go-searchquery

[![CI](https://github.com/njchilds90/go-searchquery/actions/workflows/ci.yml/badge.svg)](https://github.com/njchilds90/go-searchquery/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-searchquery.svg)](https://pkg.go.dev/github.com/njchilds90/go-searchquery)

A zero-dependency, deterministic parser for GitHub-style search queries in Go. 

Safely turn user search strings (or AI Agent prompt outputs) like `is:open label:bug "system error"` into a highly structured Abstract Syntax Tree without relying on brittle regex.

## Installation

```bash
go get github.com/njchilds90/go-searchquery

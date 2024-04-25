# prdy

A command line tool to do the things a developer can forget before submitting a pull-request. Prdy scans for specified keywords, often related to debugging or temporary code (like `console.log` or `var_dump`), and highlights any occurrences. It also provides functionalities for managing configuration settings to tailor the search according to user-specific requirements, or to automatically run the test suite after scanning.

## Purpose

1. Quality of life application to improve the work I do, even if only by a few percent
2. Increase comfort with core features of Go language and experiment with application structure

## Built using
- Go standard library
- [go-git-ignore](github.com/sabhiram/go-git-ignore) package to enhance configuration

## Technical features
- File system iteraction
- Configuration management
- Code scanning and pattern matching

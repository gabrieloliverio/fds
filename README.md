[![Go](https://github.com/gabrieloliverio/fds/actions/workflows/go.yml/badge.svg)](https://github.com/gabrieloliverio/fds/actions/workflows/go.yml)

# `fds`

Modern and opinionated find/replace CLI programme. Short version of **F**in**d** and **S**ubstitute, read as /fɔ.dɐ.s(ɨ)/ :).

# Features

- Find and replace text (and RegEx) from stdin 
- Inline replace in files[1]
- Modern PCRE RegEx, same as you use on `rg` and your favourite programming languages
- Use RegEx groups as replacement
- Case-insensitive matching
- String-literal mode - no RegEx and escaping characters when you don't need RegEx
- Find files and directories using glob double-start patterns
- Replace with interactive mode, similar to `git patch` and vim replace `/c`
- Ignore files and directories with glob double-star patterns

[1] When provided a file, it creates a temporary file, writes the content and replaces the original file, following symlinks by default.

# Usage

```
$ echo subject | fds [ options ] search_pattern replace
$ fds [ options ] search_pattern replace ./file
$ fds [ options ] search_pattern replace ~/directory
$ fds [ options ] search_pattern replace ~/directory/**/somepattern*

Options:

	-l, -literal        Treat pattern as a regular string instead of as Regular Expression
	-i, -insensitive    Ignore case on search
	-c, -confirm        Confirm each substitution
	-v, -verbose        Print debug information
	-ignore             Ignore glob patterns, comma-separated. Ex. -ignore "vendor/**,node_modules/lib/**.js"
```

# Roadmap

- [x] Stdin (pipe) + replacement as string
- [x] Positional parameters
- [x] Replacement as RegEx
- [x] "Inline" replace for single file, writing content into temp file and renaming it
- [x] Support for symlinks (similar to sed's `--follow-symlink`)
- [x] Support for string-literal mode
- [x] Support for case-insensitive mode
- [x] Interactive mode
- [x] Include line numbers
- [x] Directories
- [x] Glob
- [x] Accept --ignore-glob
- [ ] Catch interrupt signal to clean up temp files
- [ ] Ignore files listed in .gitignore
- [ ] Multiple files, directories and/or globs
- [ ] Concurrency when reading/writing several files

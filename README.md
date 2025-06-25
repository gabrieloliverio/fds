[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://github.com/gabrieloliverio/fds/actions/workflows/go.yml/badge.svg)](https://github.com/gabrieloliverio/fds/actions/workflows/go.yml)
[![Sonar Cloud](https://sonarcloud.io/api/project_badges/measure?project=gabrieloliverio_fds&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabrieloliverio_fds)

# `fds`

Modern and opinionated find/replace CLI program. Short version of **F**in**d** and **S**ubstitute, read as /fɔ.dɐ.s(ɨ)/ :).

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
echo subject | fds [ options ] search_pattern replace
fds [ options ] search_pattern replace ./file
fds [ options ] search_pattern replace ~/directory
fds [ options ] search_pattern replace ~/directory/**/somepattern*

Options:

	-l, --literal        Treat pattern as a regular string instead of as Regular Expression
	-i, --insensitive    Ignore case on search
	-c, --confirm        Confirm each substitution
	-v, --verbose        Print debug information
	--ignore-globs       Ignore glob patterns, comma-separated. Ex. --ignore-globs "vendor/**,node_modules/lib/**.js"
	--workers            Number of workers created to process the substitutions. Default value: 4

Examples:

# From stdin
echo "baz bar" | fds "baz" "foo" # prints out "foo bar"

# Replace in a file
fds "foo" "bar" ./file.txt

# Replace in files from a directory
fds "foo" "bar" ./dir

# Replace in files from a directory using 8 workers instead of 4
fds "foo" "bar" ./dir --workers 8

# Confirm each replacement
fds -c "foo" "bar" ./file.txt

# Literal mode
fds -l "->" "=>" ./file.txt

# Insensitive mode
fds -i "foo" "bar" ./file.txt

# Replace recursively in .txt files
fds "foo" "bar" ./dir/**/*.txt
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
- [x] Accept --ignore-globs
- [x] Concurrency when reading/writing several files
- [ ] Ignore files listed in .gitignore
- [ ] Multiple files, directories and/or globs

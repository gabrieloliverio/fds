# `fds`

Modern and opinionated find/replace CLI programme. Short version of **F**in**d** and **S**ubstitute, read as /fɔ.dɐ.s(ɨ)/ :).

# Features

- Find and replace text  RegEx) from stdin. 
- Inline replace in files[1]
- Modern PCRE RegEx, same as you use on `rg` and your favourite programming languages
- Use RegEx groups as replacement
- String-literal mode - no RegEx and escaping characters when you don't need RegEx
- Replace with interactive mode, similar to `git patch` and vim replace `/i`

[1] When provided a file, it creates a temporary file, writes the content and replaces the original file, following symlinks by default.

# Usage

```
$ echo "some text" | fds text replacement       # Using stdin
$ fds text replacement ./a_file                 # Reading file
$ fds text replacement ./a_directory            # Reading directory
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
- [ ] Multiple files, directories and/or globs
- [ ] Ignore files listed in .gitignore
- [ ] Accept --ignore-glob
- [ ] Catch interrupt signal to clean up temp files
- [ ] Concurrency when reading/writing several files

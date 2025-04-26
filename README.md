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
$ fds "some text" text replacement              # Using positional parameters
$ fds ./afile text replacement                  # Reading file
$ fds ./adir text replacement                   # Reading directory
$ fds dir/**/file.* text replacement            # Glob
```

# Roadmap

- [x] Stdin (pipe) + replacement as string
- [x] Positional parameters
- [x] Replacement as RegEx
- [x] "Inline" replace for single file, writing content into temp file and renaming it
- [x] Support for symlinks (similar to sed's `--follow-symlink`)
- [x] Support for string-literal mode
- [ ] Support for case-insensitive mode
- [ ] Interactive mode
- [ ] Multiple files, directories and Glob
- [ ] Concurrency when reading/writing several files

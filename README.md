# `fds`

Modern version of `sed` written in Go. Short version of **F**in**d** and **S**ubstitute, read as /fɔ.dɐ.s(ɨ)/ :).

# Features

- Replace text (and obviously RegEx) from stdin and file(s). 
- Modern PCRE RegEx - no need to learn (and decorate) all the differences of sed's RegEx
- Use RegEx groups as replacement
- String-literal mode - no RegEx and escaping characters when you don't need RegEx
- Replace with interactive mode, similar to git path and vim replace `/i`

# Usage

$ echo "some text" | fds text replacement # Using stdin
$ fds "some text" text replacement # Using positional parameters
$ fds ./afile text replacement # Reading file
$ fds ./adir text replacement # Reading directory

# Roadmap

- [x] Stdin (pipe) + replacement as string
- [x] Positional parameters
- [x] Replacement as RegEx
- [x] Single file
- [ ] Glob and directories
- [ ] Interactive mode
- [ ] Concurrent

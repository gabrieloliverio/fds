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

# Roadmap

- [x] Stdin (pipe) + replacement as string
- [x] Positional parameters
- [ ] Replacement as RegEx
- [ ] Single file
- [ ] Glob and directories
- [ ] Interactive mode
- [ ] Concurrent

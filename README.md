# Interpreter

An interpreter implementation in Go for learning.

## Story

I started reading [`Writing An Interpreter In Go (Thorsten Ball)`](https://interpreterbook.com) book. The book is about writing an interpreter for a fictional programming language called `Monkey` (the syntax is defined in that book). So, I wrote code and tried to understand how we can build an interpreter from the book. After finishing the book, I got a pretty good implementation of an interpreter. [Features from the book](#features-from-the-book)

But there were many missing features in the language that are generally present in all languages. So, I tried to complete the programming language with all those missing features. [New features I added](#new-features-added)

## Features from the book

- [x] `let` and `return` statements
- [x] Integers, booleans, strings, arrays, and hashes
- [x] Arithmetic and comparison operators
- [x] `if` / `else` conditionals
- [x] Expression evaluation with operator precedence via Pratt parsing
- [x] First-class functions and closures
- [x] Built-in functions: `len`, `first`, `last`, `rest`, `push`, `puts`

## New features added

- [ ] Line and file name tracking in generated tokens
- [ ] Full Unicode support (including emojis)
- [ ] Number formats other than integers (floats, hex, octal, etc.)
- [ ] Character escaping in strings

# Interpreter

An interpreter implementation in Go for learning.

## Story

I started reading [Writing An Interpreter In Go](https://interpreterbook.com) by Thorsten Ball. The book is about writing an interpreter for a fictional programming language called `Monkey` (the syntax is defined in the book). I wrote the code and tried to understand how to build an interpreter by following along. After finishing the book, I had a pretty good implementation of an interpreter. [Features from the book](#features-from-the-book)

But the language was missing many features that are generally present in all languages. So I tried to complete the programming language by adding those missing features. [New features I added](#new-features-added)

I also refactored the code to be more robust, maintainable, and readable.

## Features from the book

- [x] `let` and `return` statements
- [x] Integers, booleans, strings, arrays, and hashes
- [x] Arithmetic and comparison operators
- [x] `if` / `else` conditionals
- [x] Expression evaluation with operator precedence via [Pratt parsing](https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing)
- [x] First-class functions and closures
- [x] Built-in functions: `len`, `first`, `last`, `rest`, `push`, `puts`

## New features added

- [ ] Line and file name tracking in generated tokens
- [ ] Full Unicode support (including emojis)
- [ ] Number formats other than integers (floats, hex, octal, etc.)
- [ ] Character escaping in strings

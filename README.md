## TODO for later

1. Add line number and file name in generated token
2. Support full unicode (and emojis!)
3. Support for others numbers then integers (floats, hex notations, octal notations, etc)

## Notes

1. There is only two statements in the monkey programming language
    - let statement
    - return statement
2. The parser is a Pratt parser (top-down operator precedence parsing): each
   token type registers a nud (null denotation, parsed with no left-hand
   expression -- prefix position) and/or a led (left denotation, parsed
   given an already-parsed left-hand expression -- infix/postfix position),
   and precedence values decide how far a led keeps consuming to its right
   before control returns to the caller.

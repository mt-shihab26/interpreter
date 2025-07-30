# interpreter

A interpreter written in Go.

## Notes

### Parts of the Interpreter

- the lexer
- the parser
- the ast (Abstract Syntax Tree)
- the internal object system
- the evaluator

### Parsing

- AST (abstract syntax tree)
- Parser generator
    - CFG
    - BNF
    - EBNF
- There are two main strategies of parsing
    - top-down parsing
        - recursive descent parsing (Pratt parser, inventor Vaughan Pratt)
        - early parsing
        - predictive parsing
    - bottom-up parsing
- follow works
    - formal proof of its correctness
    - error-recovery process
    - detection of erroneous syntax

## Syntax of Monkey Programming Language

### Bind Value

```monkey
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2)
```

### Array

```monkey
let myArray = [1, 2, 3, 4, 5]; // array
myArray[0] // => 1
```

### Hash (Map)

```monkey
let throsten = {"name": "Throsten", age: "28"} // hash
throsten["name"] // => "Throsten"
```

### Function

```monkey
const add = fn(a, b) { return a + b; };
const add2 = fn(a, b) { a + b; }; // implicit return

add(2, 1) // calling function

let fibonacci = fn(x) {
    if (x) {
       0
    } else {
        if (x == 1) {
            1
        } else {
            fibonacci(x - 1) + fibonacci(x - 2)
        }
    }
};

// first class functions
let twice = fn(f, x) {
    return f(f(x))
}

let addTwo = fn(x) {
    return x + 2;
}

twice(addTwo, 2) // => 6
```

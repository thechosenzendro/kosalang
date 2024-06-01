# Basic data types
## Integers
Whole numbers (```1```), can be negated (```-47```) and can be separated with underscores (```1_000_000```).
## Floats
Decimals (```0.5```), can also be negated (```-42.3```) and seperated by underscores (```10_000.5```).
## Strings
UTF-8 encoded text delimited by quotation marks with support for interpolating expressions by putting them into curly brackets. Can be multiline.
```python
"👋 {world}!"
```
## Lists
A list of values of the same type, delimited by square brackets and separated by commas.
```python
[int] numbers = [1, 2, 3]
```
## Dictionares
*This is not that great of a syntax, i should reconsider it. Also, are dicts even neccessary?*
```python
{string: string} surnames = {"John": "Doe", "Jane": "Doe"}
```
## Bools
Simply ```true``` or ```false```. Can be negated with the ```not``` keyword.

# Custom types
Custom types are declared with the ```data``` keyword. They can have states with, or without entries. If there are entries at the top level, they need to be filled out regardless of the variant.
```python
data Person:
    int birth_date
    Alive:
        bool hungry
    Dead:
        int death_date
    Unknown
```

# Constants
All variables cannot be changed, so they are constants. The declaration can be optionally started with a type, then the constant name, equal sign, and the expression. 
```python
int one = 1
two = 2
```
U can use a constant by refering to it by its name, like ```one```.

# Control flow
Control flow is handled by the ```match``` expression. It can be used with or without a variable. It contains a list of matchees, if a matchee matches, it runs the code inside it, and then returns the resulting value. If no match is found, the _ matchee is used. The compiler enforces that all of the possible states are matched. You can match on a value, structure, type, or a data state.
```python
one = 1
match one:
    0: "zero"
    _: "default"
```
## Patterns
A pattern can be any concrete value. Like ```0```, ```"Hello, world!"``` or ```[1, 2, 3]```.
### List
```python
# Empty list
[]
# List that has the length of 1 and the number 47.
[47]
# List that has the length of 1 and any value
[_]
# List that has variable length and the number 47
[47, ...]
# List that has variable length, the number 47 and all the other values are assigned to the constant other (bad syntax, redesign.)
[47, ...[int]other]

```
### Alternative patterns
```python
match one:
    1 or 2: "it is either one or two"
    _: "default"
```

*Add more examples.*

# Sharing code
*Maybe rethink this?*
Modules are KoSALANG files. Types, functions and constants can be imported with the ```using``` keyword. You can import the module (thus having a namespace), or import just the stuff you want, with the .() syntax, and do all of that on a single line in multiples, separated by commas. Everything in the modules is private by default, and it can be made accessible by using the ```pub``` keyword.
```python
using http.(serve, HTTPStatusCode), random

pub hello = "World"
```

# Math and logical expression
```+```, ```-```, ```*```, ```/``` are self explanatory. ```%``` is a modulo operator. Parentheses can be used to make sure math expr is executed first.
```python
1 + 1 # 2
1 - 1 # 0
1 * 2 # 2
1 / 2 # 0.5
1 % 0.3 # 0.1
```

For equality, ```is``` and ```is not``` is used than, we just use ```<```, ```>```, ```<=``` and ```>=```. ```or``` and ```and``` are used for logical or and and.
```python
number = 2

match:
    number is 0: "zero"
    number <= 0: "negative"
    number >= 0: "positive"
```

# Types
## Unions
Unions can be made with ```or```.
```python
number = int or float
```
## Generics
1. "any"
If you dont wana restrict the type at all, just use a normal type variable
```python
T test_pass(T struct):
    struct
```
2. Restricting based on functions
What if we want only the types which i can call io.format and enum.map on?
```python
T test_pass(T(io.format, enum.map) struct):
    struct
```
# Functions
Functions are comprised of 3 parts.

1. Signature
Signature is a return type followed optionally by a function name, open parenthesis, list of arguments and a close parenthesis. It doesnt mean anything on its own, it is used either as a type, or in a function definition.
```python
# with name
int sum(int a, int b)
# without it
int (int a, int b)
```
2. Function definition
Function definition is just a signature with a block that implements it. KoSALANG doesnt have a return keyword, so the value of the last expression is used as a return value.
```python
int sum(int a, int b):
    a + b
```
You can also specify default values like this:
```python
int sum(int a: 0, int b: 1):
    a + b
```
The definition is an expression, so it can be used as an argument.
```python
# the definition for map defines all the argument types, so you dont need to specify them
[1, 3, 2].map((x): x*2)
```

3. Function call
Function call can be an ident or a function definition. You can call be the function by adding the ident and a list of arguments enclosed in parentheses. The arguments can also be named, but they need to go after the positional arguments.
```python
io.log("Hello World!")

io.log("Hello World!", end: "\n\n")
```
You can call it with a dot in the middle, if the left side has the type of the first argument like this:
```python
hello = "Hello World!".split(" ")
```
You can also chain calls like this:
```python
app = http.server()
    .add_route("/", index)
    .set_limit(10_000)
    .serve(8080)
# this is equivalent to this:
serve(set_limit(add_route(http.server(), "/", index), 10_000))
```
# To add
- some way of documentation

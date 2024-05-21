# Basic data types
## Integers
Whole numbers ```1```), can be negated (```-47```) and can be separated with underscores (```1_000_000```).
## Floats
Decimals (```0.5```), can also be negated (```-42.3```) and seperated by underscores (```10_000.5```).
## Strings
UTF-8 encoded text delimited by quotation marks with support for interpolating expressions by putting them into curly brackets.
```
"👋 {"world"}!"
```
## Lists
A list of values of the same type, delimited by square brackets and separated by commas.
```
int[1, 2, 3]
```
## Dictionares
*This is not that great of a syntax, i should reconsider it. Also, are dicts even neccessary?*
```
string: string{"John": "Doe", "Jane": "Doe"}
```
## Bools
Simply ```true``` or ```false```

# Custom types
Custom types are declared with the ```data``` keyword. They can have states with, or without entries. If there are entries at the top level, they need to be filled out regardless of the variant.
```
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
```
int one = 1
two = 2
```
U can use a constant by refering to it by its name, like ```one```.

# Control flow
Control flow is handled by the ```match``` expression. It can be used with or without a variable. It contains a list of matchees, if a matchee matches, it runs the code inside it, and then returns the resulting value. If no match is found, the _ matchee is used. The compiler enforces that all of the possible states are matched. You can match on a value, structure, type, or a data state.
```
one = 1
match one:
    0: "zero"
    _: "default"
```
*Add more examples.*

# Sharing code
*Maybe rethink this?*
Modules are KoSALANG files. Types, functions and constants can be imported with the ```using``` keyword. You can import the module (thus having a namespace), or import just the stuff you want, with the .() syntax, and do all of that on a single line in multiples, separated by commas. Everything in the modules is private by default, and it can be made accessible by using the ```pub``` keyword.
```
using http.(serve, HTTPStatusCode), random

pub hello = "World"
```
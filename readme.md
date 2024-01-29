# Python

## Good:
* readable syntax

## Bad
* types are just hints
* interpreted
* slow
* pip

# Go
## Good
* compiles to one .exe
* feature heavy std
* goroutines

## Bad
* imports can change behaviour
* nil
* error handling (errors as values are cool tho)

# Rust
## Good
* type system
* compiled
* enums
* cargo
* borrow checker
* macros

## Bad
* need to install crate to do anything
* dense syntax
* slow compile times
* BC rules can feel arbitrary

# Elixir

## Good 
* compiles to BEAM
* lightweight proccesses
* error handling
* pheonix live view
* metaprogramming
* modules
* examples double as tests
## Bad
* weird syntax
* no types
* bad lsp & tooling

# KoSALANG
## Should have/be...
* functional
* compiled to an executable
* types
* traits
* modules
* feature heavy std
* easy to understand imports (using and pub)
* some sort of autofree memory
* awesome tooling

# Ideas

# Examples
Examples arent valid syntax (for now)
## Struct definition
```python
# Make a "class"

struct Person:
    string name

greet(Person p):
    IO.log("Hello, my name is", p.name)

mark = Person(name: "Mark")
mark.greet()

# mark.greet() is the same as greet(mark)
# Maybe change this to a reciever?
```
## Modules
```python
module Math:
    int sum(int a, int b):
        a + b

three = Math.sum(1, 2)
```
## Easy to understand imports (using and pub)
```python
# Uses the Math module from earlier
# Makes the Math module available in scope
using Math

three = Math.sum(1, 2)
```

## Enums
```python
# values in structs that do not have a type are states
using DateTime

struct Person:
    Alive:
        bool is_hungry
    Dead:
        DateTime death_date

mark = Person.Alive(is_hungry: true)
```
## Error handling
```python
int or Error game(int value):
    if value == 5:
        return Error(msg: "Five is bad")
    else:
        return value

three = game(3)?
_error = game(5)?
# Previous statement returns the Error
```
## Traits



## Autofree
```python
[int] numbers = [1, 2, 3]
```
=>

```python
[int] numbers = [1, 2, 3]
free(numbers)
```
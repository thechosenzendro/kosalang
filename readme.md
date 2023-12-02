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
* gorutines

## Bad
* imports can change behaviour
* bad ecosystem

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
* MONADS?!?!
* readable syntax (C inspired)
* compiled to an executable
* types
* modules
* feature heavy std
* easy to understand imports (using)
* enums
* some sort of autofree memory
* macros (functions that run on compiletime and runtime.)
* awesome tooling (lsp the quality of Pylance)

# Examples
## Readable syntax (C inspired)
### Struct definition
```
struct person {
    int age,
    str name,
    @impl
    void greet (person self) {
        IO.log("Hi! My name is {} and i am {} years old.", self.name, self.age)
    }
}

person marek = person{"age": 15, "name": "Marek Pokorný"}
marek.greet()
```
## Modules
```
module Math {
    num sum(num a, num b) {
        a + b
    }    
}

Math::sum(1, 2)
```
## Easy to understand imports (using and pub)
```
# math.kosa
pub module Math {
    pub num sum(num a, num b) {
        a + b
    }    
}
# main.kosa
using math::Math
Math::sum(1, 2)
```

## Enums
```
enum Person {
    Alive {
        bool hungry
    }
    Dead {
        bool buried
    }
}
```

## Autofree
```
[num] numbers = [1, 2, 3]
```
=>

```
[num] numbers = [1, 2, 3]
free(numbers)
```
## Macros
Coming soon!


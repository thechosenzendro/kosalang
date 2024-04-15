The compiler for the KOSA programming language.
It's a mix of all of my favourite features of programming languages.

# Examples
```python
# Variable assignment
hello = "World"

# Functions
int add(int a, int b):
    # No return keyword
    a + b

# Structs
struct Person:
    int age
    string name
john = Person(34, "John Doe")

# Control flow
using io
cond = true
if cond:
    io.log("It works!")
else:
   io.error("What?")

names = ["John", "Jane", "Jacob"]
for name in names:
   io.log(f"Great name, {name}!")

# HTML literals
header = <h1>Hello World!</h1>
io.log(header.children)

# Pattern matching and error handling
using errors.error
result(string, error) get_last_name(string first_name):
    match first_name:
        "John": "Doe"
        "Jacob": "Smith"
        _: error("Uknown first name")

last_name = get_last_name("John")? # Doe
last_name = get_last_name("Adam")? # panic

   

```
# ...and more!
Language server is my biggest priority.

# gudge

Autotester for competitive programming problems written in Go 

## Usage:

The tests should be stored in the same directory as the solution with .txt file extension.
In this example the two input test cases are `1 2`, `2 3` with the expected outputs `3`, `5` respectively.:

```plaintext
1 2



3



2 3



5
```

Check ./testcases/ for the full example.

To run the tests using `<command> <solution>`.

```bash
#!/bin/bash
gudge <command> <solution>
```
Example:
```bash
#!/bin/bash
gudge python3 solution.py
gudge go run solution.go
```

## Installation, Build:

```bash
#!/bin/bash
cd ~/Projects
git clone git@github.com/TemaSaur/gudge.git
cd gudge
go build -o ~/bin
```


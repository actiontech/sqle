# Contributing Guidelines

## Golden rule

Speak up _before_ writing code. Comment on existing issue or create a new one. Discuss what
you want to implement _before_ implementing it.


## Getting code

Fork repository on GitHub and clone the source code:

```
git clone git@github.com:<your name>/reform.git
cd reform
make init
```

Thanks to Go modules, that will work in any directory. Make sure they are not disabled in your environment.

Please read the "Versioning and branching policy" section in README,
and send pull requests to the right branch.


## Makefile targets

* Run `make` without arguments to see all Makefile targets.
* Run `make env-up` or `make env-up-detach` to start databases with Docker Compose.
* Run `make init` to install development tools.
* Run `make test` to run all tests.
* Run `make lint` to run linters.

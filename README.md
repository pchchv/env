# **env** [![Go Reference](https://pkg.go.dev/badge/github.com/pchchv/golog.svg)](https://pkg.go.dev/github.com/pchchv/golog) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/env)](https://goreportcard.com/report/github.com/pchchv/env)
Go library which loads environment variables from .env files.  
Ruby project port [dotenv](https://github.com/bkeepers/dotenv).

> Storing configuration in the environment is one of the tenets of a [twelve-factor app](https://12factor.net). Anything that is likely to change between deployment environments–such as resource handles for databases or credentials for external services–should be extracted from the code into environment variables.
>
> But it is not always practical to set environment variables on development machines or continuous integration servers where multiple projects are run. Dotenv load variables from a .env file into ENV when the environment is bootstrapped.

It can be used as a library (to load into env for your own daemons, etc.) or as a bin command.

## Installation

As a library

```bash
go get github.com/pchchv/env
```

As a bin command

go >= 1.17
```bash
go install github.com/pchchv/env/cmd/env@latest
```

go < 1.17
```bash
go get github.com/pchchv/env/cmd/env
```

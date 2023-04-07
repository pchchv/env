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

## Usage

Add your application configuration to your `.env` file in the root of your project:

```bash
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

Then in your Go app you can do something like

```go
package main

import (
    "log"
    "os"

    "github.com/pchchv/env"
)

func main() {
  err := env.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  s3Bucket := os.Getenv("S3_BUCKET")
  secretKey := os.Getenv("SECRET_KEY")

  // now do something with s3 or whatever
}
```

Also, you can simply take advantage of the autoload package which will read `.env` when importing

```go
import _ "github.com/pchchv/env/autoload"
```

Although `.env` in the root of the project is used by default, you should not be restricted, both of the following examples are 100% acceptable

```go
env.Load("somerandomfile")
env.Load("filenumberone.env", "filenumbertwo.env")
```

If you want to be really fancy with your env file, you can make comments and export (below is the correct env file)

```bash
# I am a comment and that is OK
SOME_VAR=someval
FOO=BAR # comments at line end are OK too
export BAR=BAZ
```

Or finally you can do YAML(ish) style

```yaml
FOO: bar
BAR: baz
```

If you don't want env to change your environment, you can just get the map back

```go
var myEnv map[string]string
myEnv, err := env.Read()

s3Bucket := myEnv["S3_BUCKET"]
```

or from an `io.Reader` instead of a local file

```go
reader := getRemoteFile()
myEnv, err := env.Parse(reader)
```

... or from a `string` if you so desire

```go
content := getRemoteFileContent()
myEnv, err := env.Unmarshal(content)
```

### Precedence & Conventions

Existing envs take precedence of envs that are loaded later.

[convention](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use)
for managing multiple environments (i.e. development, testing, production)
is to create an environment named `{YOURAPP}_ENV` and load the environments in that order:

```go
env := os.Getenv("FOO_ENV")
if "" == env {
  env = "development"
}

env.Load(".env." + env + ".local")
if "test" != env {
  env.Load(".env.local")
}
env.Load(".env." + env)
env.Load() // The Original .env
```

If necessary, you can also use `env.Overload()` to break this convention and overwrite existing envs instead of just replacing them. Use with caution.

### Command Mode

Suppose you installed the command as above, and you have `$GOPATH/bin` in `$PATH`.

```
env -f /some/path/to/.env some_command with some args
```

If you do not specify `-f`, it will load `.env` into `PWD` by default.

By default it will not override existing environment variables; you can do this with the `-o` flag.

### Writing Env Files

env can also write a map representing the environment to a correctly-formatted and escaped file

```go
env, err := env.Unmarshal("KEY=value")
err := env.Write(env, "./.env")
```

or to a string

```go
env, err := env.Unmarshal("KEY=value")
content, err := env.Marshal(env)
```
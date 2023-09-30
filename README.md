# Config
Config is a Go package for loading configuration values from AWS SSM Parameter Store and environment variables into a struct.

## Installation

To install Config, use go get:

```
go get github.com/ddouglas/config
```

## Usage 


To use Config, first import it:

```
import "github.com/ddouglas/config"
```

Then, define a struct to hold your configuration values. The struct fields must be tagged with ssm or env to indicate the source of the value. For example:

```
type Config struct {
    DatabaseURL string `ssm:"/database_url"`
    APIKey      string `ssm:"/api_key,required"`
    Debug       bool   `env:"DEBUG"`
}
```

To load the configuration values into your struct, call the Load function:

```
var cfg Config
err := config.Load(context.Background(), &cfg)
if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}
```

By default, Load will look for SSM parameters with names that match the struct field tags. You can specify a prefix for SSM parameter names using the WithPrefix option:

```
opts := []config.LoadOptFunc{
    config.WithPrefix("/myapp"),
}
err := config.Load(context.Background(), &cfg, opts...)
```

You can also specify an SSM client to use using the WithSSMClient option. The Client must satify the `GetParametersAPIClient` interface:

```
type GetParametersAPIClient interface {
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
}

client := myCustomSSMClient{}
opts := []config.LoadOptFunc{
    config.WithSSMClient(client),
}
err := config.Load(context.Background(), &cfg, opts...)
```

If a required value is missing, or if there is an error fetching values from SSM, Load will return an error.

## Example
Here's an example of how to use Config to fetch environment variables and SSM parameters:

```
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ddouglas/config"
)

type Config struct {
    DatabaseURL string `ssm:"/database_url"`
    APIKey      string `ssm:"/api_key,required"`
    Debug       bool   `env:"DEBUG"`
}

func main() {
    var cfg Config
    err := config.Load(context.Background(), &cfg)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    fmt.Printf("DatabaseURL: %s\n", cfg.DatabaseURL)
    fmt.Printf("APIKey: %s\n", cfg.APIKey)
    fmt.Printf("Debug: %t\n", cfg.Debug)
}
```
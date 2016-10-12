
# Overview

`lambazon` attempts to provide a `net/http` request/response environment
based on the [caddy-awslambda](https://github.com/coopernurse/caddy-awslambda)
JSON envelope format.

It is intended to be used with [Apex](http://apex.run/) as a deploy tool.

This enables you to develop and test a Go web app locally as you would normally,
then deploy and run it at AWS Lambda *with no changes*.

This has been lightly tested and should be considered experimental.
Please file issues with ideas / problems.

## Example main.go file

This is how simple your Apex function `main.go` can be:

```go
package main

import (
	"github.com/coopernurse/lambazon"
	"github.com/coopernurse/lambazon/examples/calc/calcweb"
)

func main() {
   // Run accepts any http.Handler instance
	lambazon.Run(calcweb.NewRouterForLambda())
}
```

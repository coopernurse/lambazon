package main

import (
	"github.com/coopernurse/lambazon"
	"github.com/coopernurse/lambazon/examples/calc/calcweb"
)

func main() {
	lambazon.Run(calcweb.NewRouterForLambda())
}

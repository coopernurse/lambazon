package main

import (
	"github.com/coopernurse/lambazon/examples/calc/calcweb"
)

func main() {
	calcweb.NewRouter().Run()
}

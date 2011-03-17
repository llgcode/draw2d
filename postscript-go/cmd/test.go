package main

import (
	"fmt"
)

func main() {
	toto := make([]int, 2, 2)
	toto[0] = 1
	toto[1] = 2
	fmt.Printf("%v\n", toto)
	toto = toto[0:0]
	fmt.Printf("%v\n", toto)
	fmt.Printf("%v\n", cap(toto))
}

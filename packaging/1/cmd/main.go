package main

import (
	"fmt"

	"github.com/mmfalcao/go-expert/packaging/1/math"
)

func main() {
	m := math.NewMath(1, 2)
	fmt.Println("Soma entre os valores ", m, " é igual á ", m.Add())
	fmt.Println(math.X)
}

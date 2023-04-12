package main

import (
	"fmt"

	"github.com/mmfalcao/go-expert/packaging/3/math"
)

func main() {
	m := math.NewMath(1, 2)
	fmt.Println("Soma entre os valores ", m, " é igual á ", m.Add())
	fmt.Println(math.X)
}

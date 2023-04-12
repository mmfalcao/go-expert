package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mmfalcao/go-expert/packaging/4/math"
)

func main() {
	m := math.NewMath(1, 2)
	fmt.Println("Soma entre os valores ", m, " é igual á ", m.Add())
	fmt.Println(math.X)
	println(uuid.New().String())
}

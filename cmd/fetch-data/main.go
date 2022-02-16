package main

import (
	"fmt"

	"github.com/pjuzeliunas/nilan"
)

func main() {
	c := nilan.Controller{Config: nilan.Config{NilanAddress: "192.168.1.31:502"}}
	errors, _ := c.FetchErrors()
	fmt.Println(errors)
}

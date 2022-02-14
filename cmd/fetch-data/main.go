package main

import (
	"fmt"
	"time"

	"github.com/pjuzeliunas/nilan"
)

func main() {
	c := nilan.Controller{Config: nilan.Config{NilanAddress: "192.168.1.31:502"}}
	for i := 0; i <= 100; i++ {

		r, err := c.FetchReadings()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(r)

		s, err := c.FetchSettings()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(s)
		time.Sleep(time.Second)

	}

}

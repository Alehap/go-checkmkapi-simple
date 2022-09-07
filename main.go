package main

import (
	"fmt"
	"os"
	"hdapi/checkmkapi"
)

func main() {
	cmk := checkmkapi.New("https://monitor.com/check_mk/",os.Getenv("CMK_USER"),os.Getenv("CMK_SECRET"))
	fmt.Println(cmk.AddHost("test","1.2.3.4", "edge"))
	fmt.Println(cmk.DeleteHost("test"))
	// fmt.Println(cmk)
}

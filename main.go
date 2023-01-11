package main

import (
	"fmt"
	"os"
	// "time"
	"hdapi/checkmkapi"
)

func main() {
	cmk := checkmkapi.New("https://"+os.Getenv("CMK_HOST")+"/cmk/check_mk/",os.Getenv("CMK_USER"),os.Getenv("CMK_SECRET"))
	// netIn, netOut := cmk.GetAvgNetworkByHostname("abc.com")
	// In Mbps, Out Mbps
	// fmt.Println(netIn/1024/1024, netOut/1024/1024)
	cpu := cmk.GetCPUUtilByHostname("abc.com")
	// return percent, e.g. 0.55 = 0.55%; 23.55 = 23.55%
	fmt.Println(cpu)
	// fmt.Println(cmk.AddHost("test","1.2.3.4", "edge"))
	// fmt.Println(cmk.DeleteHost("test"))
	// fmt.Println(cmk)
}

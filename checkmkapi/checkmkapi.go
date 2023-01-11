package checkmkapi

import (
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"
	// "time"
	// b64 "encoding/base64"
	"encoding/json"
)

type cmk struct {
	result 		string
	result_code	int
}

type account struct {
    user			string
    secret			string
    cmkURL			string
}

func New(cmkURL string, user string, secret string) account {  
	a := account {user: user, secret: secret, cmkURL: cmkURL}
    return a
}

func (a account) httpRequest(endpoint string, request string) (bool,[][]interface{}) {
	// fmt.Println(request)
	client := &http.Client{}
	var data = strings.NewReader(request)
	req, err := http.NewRequest("POST", a.cmkURL + endpoint + ".py?_username="+a.user+"&_secret="+a.secret, data)
	if err != nil {
		log.Fatal(err)
		return false, nil
	}
	// req.Header.Set("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(a.user+":"+a.secret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err1 := client.Do(req)
	if err1 != nil {
		log.Fatal(err1)
		return false, nil
	}
	defer resp.Body.Close()
	bodyText, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
		return false, nil
	}
	// fmt.Printf("%s\n", bodyText)
	byt := []byte(bodyText)
	var obj [][]interface{}
	err3 := json.Unmarshal(byt, &obj);
	if err3 != nil {
        log.Fatal("Parse json err: ")
        log.Fatal(err3.Error())
        return false, nil
    }
    // fmt.Println(obj[1][32])
	return true, obj
}

func (a account) makeRequest(action string, request string) (bool,string) {
	// fmt.Println(action + " is called.")
	client := &http.Client{}
	var data = strings.NewReader(request)
	req, err := http.NewRequest("POST", a.cmkURL + "webapi.py?_username="+a.user+"&_secret="+a.secret+"&action="+action, data)
	if err != nil {
		log.Fatal(err)
		return false, err.Error()
	}
	// req.Header.Set("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(a.user+":"+a.secret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err1 := client.Do(req)
	if err1 != nil {
		log.Fatal(err1)
		return false, err1.Error()
	}
	defer resp.Body.Close()
	bodyText, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2)
		return false, err2.Error()
	}
	// fmt.Printf("%s\n", bodyText)
	// byt := []byte(bodyText)
	var obj map[string]interface{}
	err3 := json.Unmarshal(bodyText, &obj);
	if err3 != nil {
        log.Fatal("Parse json err: ")
        log.Fatal(err3)
        return false, err3.Error()
    }
    // fmt.Println(obj)
    return true, "ok"
}

func (a account) discoveryServices(hostname string) (bool, string) {
	request := `request={"hostname":"`+hostname+`","mode":"refresh"}`
	action  := "discover_services"
	return a.makeRequest(action, request)
}
func (a account) activeChanges(sitename string) (bool, string) {
	request := `request={"sites":["`+sitename+`"],"allow_foreign_changes":"1"}`
	action  := "activate_changes"
	return a.makeRequest(action, request)
}
func (a account) AddHost(hostname string, ip string, folder string) (bool, string) {
	request := `request={"hostname":"`+hostname+`","folder":"`+folder+`","attributes":{"ipaddress":"`+ip+`"},"create_folders":"1"}`
	action  := "add_host"
	stt1,err1 := a.makeRequest(action, request)
	if !stt1 {
		// fmt.Println("Loi ne")
		log.Fatal(stt1,err1)
		return false, err1
	}
	stt2,err2 := a.discoveryServices(hostname)
	if !stt2 {
		log.Fatal(err2)
		return false, err2
	}
	stt3,err3 := a.activeChanges("monitor")
	if !stt3 {
		log.Fatal(err3)
		return false, err3
	}
	return true, "ok"
}
func (a account) DeleteHost(hostname string) (bool, string) {
	request := `request={"hostnames":["`+hostname+`"]}`
	action  := "delete_hosts"
	stt1,err1 := a.makeRequest(action, request)
	if !stt1 {
		// fmt.Println("Loi ne")
		log.Fatal(stt1,err1)
		return false, err1
	}
	stt3,err3 := a.activeChanges("monitor")
	if !stt3 {
		log.Fatal(err3)
		return false, err3
	}
	return true, "ok"
}

func (a account) RenameHost(oldhostname string, newhostname string, ip string, folder string) (bool, string) {
	request := `request={"hostnames":["`+oldhostname+`"]}`
	action  := "delete_hosts"
	stt1,err1 := a.makeRequest(action, request)
	if !stt1 {
		// fmt.Println("Loi ne")
		log.Fatal(stt1,err1)
		return false, err1
	}
	return a.AddHost(newhostname, ip, folder)
}

func (a account) GetAvgNetworkByHostname(hostname string) (float64, float64) {
	request := `host=`+hostname+`&service=Interface+2&view_name=service&output_format=json`
	stt1,data := a.httpRequest("view", request)
	if !stt1 {
		// fmt.Println("Loi ne")
		// log.Fatal(stt1,err1)
		return 0.0, 0.0
	}
	lastIn, _ := strconv.ParseFloat(strings.Split(strings.Split(data[1][32].(string),";")[4],"=")[1], 64)
	lastOut, _ := strconv.ParseFloat(strings.Split(strings.Split(data[1][32].(string),";")[8],"=")[1], 64)
 	return lastIn*8, lastOut*8
}

func (a account) GetCPUUtilByHostname(hostname string) (float64) {
	request := `host=`+hostname+`&service=CPU+utilization&view_name=service&output_format=json`
	stt1,data := a.httpRequest("view", request)
	if !stt1 {
		// fmt.Println("Loi ne")
		// log.Fatal(stt1,data)
		return 0.0
	}
	// fmt.Println(data[1][32])
	cpu, _ := strconv.ParseFloat(strings.Split(strings.Split(data[1][32].(string),";")[12],"=")[1], 64)
 	return cpu
}

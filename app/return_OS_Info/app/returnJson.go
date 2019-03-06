package main

import (
	"encoding/json"
	"os"
	"log"
	"net/http"
	"github.com/comail/colog"
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"net"
)


var returner []byte

type healthInfo struct {
	Status	string	`json:"status"`
}

type osInfo struct {
    Ipaddress    []string  `json:"ipaddress"`
    Hostname  string `json:"hostname"`
    Os  string `json:"os"`
}

//functions for each handlers
func healthHandler (w http.ResponseWriter,r *http.Request) { //what does this star "*" mean?
	result,_ := returnHealth()
	w.Header().Set("Content/Type","Application/json")
	w.Write(result)
	log.Printf("info: Health check completed.")
}

func dirHandler (w http.ResponseWriter,r *http.Request) {

	var result []byte
	var logBody string
	var rtype string
	var ctype string

	switch method := r.Method; method {
	case "GET":
		result,rtype = returnJson()

		if rtype == "json" {
			ctype = "Application/json"
		} else {
			ctype = "text/plain"
		}

	case "POST":
		//http.Request.Body is type of io.ReadCloser. Converting io.ReadCloser to String.
		//See https://golangcode.com/convert-io-readcloser-to-a-string/
		rbody := &returner
		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		logBody = bufbody.String() //Converted. var "logBody" is for print body on log.
		*rbody = []byte(logBody)

		//checking json or not.
		body := []byte(logBody)
		rtype,_ := jsonCheck(body)

                if rtype == "json" {
                        ctype = "Application/json"
                } else {
                        ctype = "text/plain"
                }

		err := r.ParseForm()
		if err != nil {
			log.Printf("error: Got error in parseing POST form.")
		}

		result = []byte("CONTENT UPDATED. Contetnt type is "+ctype+".\n")

	case "PUT":
		fmt.Println("PUT method function comming soon.")

	case "DELETE":
		rbody := &returner
		*rbody = nil
		result = []byte("CONTENT DELETED.\n")
	}

	w.Header().Set("Content-Type",ctype)
	w.Write(result)
	log.Printf("info: Received access.")
	log.Printf("info: Protocol: %s",r.Proto)
	log.Printf("info: Method: %s",r.Method)
	log.Printf("info: Header: %s",r.Header)
	log.Printf("info: Body: %s",logBody)
	log.Printf("info: Host: %s",r.Host)
	log.Printf("info: Form: %s",r.Form)
	log.Printf("info: PostForm: %s",r.PostForm)
	log.Printf("info: RequestURI: %s",r.RequestURI)
	log.Printf("info: Response: %v",r.Response)
}

// for health check
func returnHealth () ([]byte,error) {
	p1 := &healthInfo{
		Status: "up",
	}

	p2,err := json.Marshal(p1)
	if err != nil {
		log.Fatal(err)
	}
	return p2,err
}

//returning json with uint8
func returnJson () ([]byte,string)  {

        iparray := getIpAddrs()
        hostname,_ := os.Hostname()
        os := runtime.GOOS
	osinfo,_ := returnOsInfo(iparray,hostname,os)

	//var parsedJson interface{}
	//var rtype string
	b1 := returner
	if len(b1) <= 0 {
		b1 = osinfo
	}
	j1 := b1
	rtype,_ := jsonCheck(j1)
	return j1,rtype
}

func jsonCheck (j0 []byte) (string, error) {
	var parsedJson interface{}
	var rtype string
	err := json.Unmarshal(j0,&parsedJson)
        if err != nil {
                rtype = "string"
        } else {
                rtype = "json"
        }
	return rtype,err
}


func getIpAddrs () []string {
        var (
                returnArray []string
                ip net.IP
        )

        ifaces, err := net.Interfaces()

        if err != nil {
                log.Fatal(err)
        }

        for _, i := range ifaces {
                addrs, err := i.Addrs()
                if err != nil {
                        log.Fatal(err)
                }
                for _,addr := range addrs {
                        switch v:= addr.(type) {
                        case *net.IPNet:
                                ip = v.IP
                                returnArray = append(returnArray,ip.String())
                        case *net.IPAddr:
                                ip = v.IP
                                returnArray = append(returnArray,ip.String())
                        }
                }
        }
        return returnArray
}

func returnOsInfo (iparray []string,hostname string,os string) ([]byte, error)  {
        var ipstring string

        for _,s := range iparray {
                ipstring = ipstring+","+s
        }

        ipstring = strings.TrimLeft(ipstring,",")

        p1 := &osInfo{
                Ipaddress: []string{ipstring},
                Hostname: string(hostname),
                Os: string(os),
        }

        p2, err := json.Marshal(p1)
        if err !=nil {
                log.Fatal(err)
        }
        return p2,err
}


func main () {
	address := ":"+os.Getenv("PORT")
	// logging configurations
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag: log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	// handlers
	http.HandleFunc("/dir",dirHandler)
	http.HandleFunc("/",healthHandler)

	log.Printf("debug: Application is started.")
	log.Fatal(http.ListenAndServe(address,nil))
}

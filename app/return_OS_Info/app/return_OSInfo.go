package main

import (
	"encoding/json"
	"net"
	"os"
	"log"
	"strings"
	"net/http"
	"github.com/comail/colog"
	"runtime"
)


//define constructions
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
        iparray := getIpAddrs()
        hostname,err := os.Hostname()
	os := runtime.GOOS

	if err != nil {
                log.Fatal(err)
		return
        }
        result,_ := returnOsInfo(iparray,hostname,os)
	w.Header().Set("Content-Type","Application/json")
	w.Write(result)
	log.Printf("info: Received access.")
	log.Printf("info: Protocol: %s",r.Proto)
	log.Printf("info: Method: %s",r.Method)
	log.Printf("info: Header: %s",r.Header)
	log.Printf("info: Body: %s",r.Body)
	log.Printf("info: Host: %s",r.Host)
	log.Printf("info: Form: %s",r.Form)
	log.Printf("info: PostForm: %s",r.PostForm)
	log.Printf("info: RequestURI: %s",r.RequestURI)
	log.Printf("info: Response: %v",r.Response)
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


func returnHealth () ([]uint8,error) {
	p1 := &healthInfo{
		Status: "up",
	}

	p2,err := json.Marshal(p1)
	if err != nil {
		log.Fatal(err)
	}
	return p2,err
}

func returnOsInfo (iparray []string,hostname string,os string) ([]uint8, error)  {
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

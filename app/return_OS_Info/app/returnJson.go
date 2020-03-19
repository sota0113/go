package main

import (
	"encoding/json"
	"os"
	"log"
	"net/http"
	"github.com/comail/colog"
	"github.com/peterbourgon/diskv"
	"bytes"
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
	var rcode int
	var result []byte
	if r.RequestURI == os.Getenv("HEALTH_DIR") {
		switch r.Method {
		case "GET":
			rcode = 200
	                w.Header().Set("Content/Type","Application/json")
	                log.Printf("info: Health check completed.")
	                result,_ = returnHealth()
		default:
			rcode = 400
		}
	} else {
		rcode = 404
	}
	w.WriteHeader(rcode)
	w.Write(result)
}

func apidirHandler (w http.ResponseWriter,r *http.Request) {

	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }
	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	d := diskv.New(diskv.Options{
	        BasePath:       "/opt/jpi-server",
	        Transform:      flatTransform,
	        CacheSizeMax: 1024 * 1024,
	});
	keyName := strings.TrimLeft(r.RequestURI,"/list/v1/api/")

        var result []byte
	var rcode int
        var logBody string
        var ctype string

	switch method := r.Method; method {

	case "GET":

		if strings.Index(keyName,"/") <= 0 {
			result, _ = d.Read(keyName)
			if len(result) <= 0 {
				rcode = 404
			} else {
				rcode = 200
				ctype,_ = jsonCheck(result)
			}
		} else {
			rcode = 404
		}

	case "PUT", "POST":
                //http.Request.Body is type of io.ReadCloser. Converting io.ReadCloser to String.
                //See https://golangcode.com/convert-io-readcloser-to-a-string/

		//r.url /vi/api/xxx/
		//trim xxx.
		// check xxx to kvs whether it already exists or not.
		// if exist, update value of key xxx with r.Body.
		// if not, create new key xxx and preserve its coresponding value r.Body to kvs.

		alwd := false
		alwdType := []string{"application/json","text/plain"}   //request returns Application Json.
		rqType := r.Header["Content-Type"][0] //application/json

		for _,v := range alwdType {
			if rqType == v {
				alwd = true
			}
		}

		if alwd == true {

			if strings.Index(keyName,"/") <= 0 {
		                bufbody := new(bytes.Buffer)
		                bufbody.ReadFrom(r.Body)
		                logBody = bufbody.String() //Converted. var "logBody" is for print body on log.

		                //checking json or not.
		                body := []byte(logBody)
		                ctype,_ = jsonCheck(body)
	                        result, _ = d.Read(keyName) //check 
	                        if len(result) <= 0 {
	                                rcode = 201 // If a request update new content, 201 would be returned.
	                        } else {
	                                rcode = 204 //If a request update existed content, 204 would be returned.
	                        }

				// Write
				d.Write(keyName, body)
				err := r.ParseForm()
				if err != nil {
					log.Printf("error: Got error on parseing PUT form.")
				}
				result = []byte("CONTENT UPDATED. Contetnt type is "+ctype+".\n")

			} else {
				rcode = 409
			}
		} else {
			rcode =409
		}

//	case "POST":	// POST request is not allowed.
//		rcode = 400

	case "DELETE":

                if strings.Index(keyName,"/") <= 0 {
                        result, _ = d.Read(keyName)

                        if len(result) <= 0 {
                                rcode = 404
                        } else {
				rcode = 204
				d.Erase(keyName)
                        }
		} else {
			rcode = 404
		}
                result = []byte("")
	default:
		rcode = 400
	}
        w.Header().Set("Content-Type",ctype)
	w.WriteHeader(rcode)
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

func dirHandler (w http.ResponseWriter,r *http.Request) {

	var result []byte
	var logBody string
	var ctype string
	var rcode int

	switch method := r.Method; method {
	case "GET":
		result,ctype = returnJson()

		if len(result) <= 0 {
			rcode = 404
		} else {
			rcode = 200
		}
        default:
                rcode = 400
	}
	w.Header().Set("Content-Type",ctype)
	w.WriteHeader(rcode)
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
	ctype,_ := jsonCheck(b1)
	return b1,ctype
}

func jsonCheck (j0 []byte) (string, error) {
	var parsedJson interface{}
	var ctype string
	err := json.Unmarshal(j0,&parsedJson)
        if err != nil {
                ctype = "text/plain"
        } else {
                ctype = "application/json"
        }
	return ctype,err
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
	apiDir := os.Getenv("MAIN_DIR")+"/api/v1/"
	// logging configurations
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag: log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	// handlers
	http.HandleFunc(os.Getenv("MAIN_DIR"),dirHandler)
	http.HandleFunc(apiDir,apidirHandler)
	http.HandleFunc(os.Getenv("HEALTH_DIR"),healthHandler)

	log.Printf("debug: Application is started.")
	log.Fatal(http.ListenAndServe(address,nil))
}

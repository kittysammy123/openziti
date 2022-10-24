package main

import "C"

import (
	"context"
	"encoding/json"
	"github.com/openziti/identity"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/openziti/sdk-golang/ziti/enroll"
	"io"
	"net"
	"net/http"
	"strings"
)

type Config struct {
	ZtAPI       string          `json:"ztAPI"`
	ID          identity.Config `json:"id"`
	ConfigTypes []string        `json:"configTypes"`
}

//export processEnrollment
func processEnrollment(src string) *C.char {
	temp := make([]byte, len(src))
	copy(temp, src)
	jwtToken := string(temp)

	var keyAlg config.KeyAlgVar = "RSA"
	var keyPath, certPath, idname, caOverride string

	tkn, _, err := enroll.ParseToken(jwtToken)
	flags := enroll.EnrollmentFlags{
		CertFile:      certPath,
		KeyFile:       keyPath,
		KeyAlg:        keyAlg,
		Token:         tkn,
		IDName:        idname,
		AdditionalCAs: caOverride,
	}

	conf, err := enroll.Enroll(flags)
	if err != nil {
		return C.CString("error")
	}
	bytes, err := json.Marshal(conf)
	if err != nil {
		return C.CString("error")
	}

	return C.CString(string(bytes))
}

var zitiContext ziti.Context
var httpConnect http.Client

func Dial(_ context.Context, _ string, addr string) (net.Conn, error) {
	service := strings.Split(addr, ":")[0] // will always get passed host:port
	return zitiContext.Dial(service)
}

//export createZitifiedHttpClient
func createZitifiedHttpClient(src string) {
	conf := make([]byte, len(src))
	copy(conf, src)

	c := Config{}
	err := json.Unmarshal(conf, &c)
	if err != nil {
		panic(err)
	}

	zitiContext = ziti.NewContextWithConfig((*config.Config)(&c))
	zitiTransport := http.DefaultTransport.(*http.Transport).Clone() // copy default transport
	zitiTransport.DialContext = Dial                                 //zitiDialContext.Dial
	httpConnect = http.Client{Transport: zitiTransport}
}

//export  zitiHttpGet
func zitiHttpGet(src string) *C.char {
	temp := make([]byte, len(src))
	copy(temp, src)
	url := string(temp)

	resp, e := httpConnect.Get(url)
	if e != nil {
		return C.CString("")
	}
	body, _ := io.ReadAll(resp.Body)
	return C.CString(string(body))
}
func main() {
	//jwt := "C:\\Users\\Administrator\\Desktop\\333.jwt"
	//config, _ := ioutil.ReadFile(jwt)
	//str1, str2 := processEnrollment(string(config))
	//fmt.Println(str1)
	//fmt.Println(str2)
}

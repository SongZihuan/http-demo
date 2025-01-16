package flagparser

import (
	"fmt"
	"strings"
)

var HttpAddress string = ":3366"
var HttpsAddress string = ""
var HttpsDomain = ""
var HttpsEmail = ""
var HttpsCertDir = "./ssl-certs"
var HttpsAliyunKey string
var HttpsAliyunSecret string
var DryRun = false
var Verbose = false

func Print() {
	fmt.Println("HttpAddress:", HttpAddress)
	fmt.Println("HttpsAddress:", HttpsAddress)
	fmt.Println("HttpsDomain:", HttpsDomain)
	fmt.Println("HttpsEmail:", HttpsEmail)
	fmt.Println("HttpsCertDir:", HttpsCertDir)
	fmt.Println("HttpsAliyunKey:", HttpsAliyunKey)
	fmt.Println("HttpsAliyunSecret:", strings.Repeat("*", len(HttpsAliyunSecret)))
}

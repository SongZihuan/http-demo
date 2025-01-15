package flagparser

import "fmt"

var HttpAddress string = ":3366"
var HttpsAddress string = ""
var HttpsDomain = ""
var HttpsEmail = ""
var HttpsCertDir = "./ssl-certs"
var ACMEAddress = ""
var DryRun = false
var Verbose = false

func Print() {
	fmt.Println("HttpAddress:", HttpAddress)
	fmt.Println("HttpsAddress:", HttpsAddress)
	fmt.Println("HttpsDomain:", HttpsDomain)
	fmt.Println("HttpsEmail:", HttpsEmail)
	fmt.Println("HttpsCertDir:", HttpsCertDir)
	fmt.Println("ACMEAddress:", ACMEAddress)
}

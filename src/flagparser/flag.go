package flagparser

import (
	"flag"
	"fmt"
)

func InitFlag() (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
	}()

	flag.StringVar(&HttpAddress, "address", HttpAddress, "http server address")
	flag.StringVar(&HttpAddress, "a", HttpAddress, "http server address")

	flag.StringVar(&HttpAddress, "http-address", HttpAddress, "http server address")
	flag.StringVar(&HttpAddress, "h", HttpAddress, "http server address")

	flag.StringVar(&HttpsAddress, "https-address", HttpsAddress, "https server address")
	flag.StringVar(&HttpsDomain, "https-domain", HttpsDomain, "https server domain")
	flag.StringVar(&HttpsEmail, "https-email", HttpsEmail, "https cert email")
	flag.StringVar(&HttpsCertDir, "https-cert-dir", HttpsCertDir, "https cert save dir")

	flag.Parse()

	return nil
}

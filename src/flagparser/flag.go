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

	flag.StringVar(&HttpAddress, "address", HttpAddress, "http server listen address")
	flag.StringVar(&HttpAddress, "a", HttpAddress, "http server listen address")

	flag.StringVar(&HttpAddress, "http-address", HttpAddress, "http server listen address")
	flag.StringVar(&HttpAddress, "h", HttpAddress, "http server listen address")

	flag.StringVar(&HttpsAddress, "https-address", HttpsAddress, "https server listen address")
	flag.StringVar(&HttpsDomain, "https-domain", HttpsDomain, "https server domain")
	flag.StringVar(&HttpsEmail, "https-email", HttpsEmail, "https cert email")
	flag.StringVar(&HttpsCertDir, "https-cert-dir", HttpsCertDir, "https cert save dir")
	flag.StringVar(&ACMEAddress, "acme-address", ACMEAddress, "acme https cert listen address")

	flag.BoolVar(&DryRun, "dry-run", DryRun, "only parser the options")

	flag.Parse()

	Print()
	return nil
}

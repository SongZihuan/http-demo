package flagparser

import (
	"flag"
	"fmt"
	"os"
)

func initFlag() (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
			return
		}
	}()

	flag.CommandLine.SetOutput(os.Stdout)

	flag.StringVar(&HttpAddress, "address", HttpAddress, "http server listen address")
	flag.StringVar(&HttpAddress, "http-address", HttpAddress, "http server listen address")

	flag.StringVar(&HttpsAddress, "https-address", HttpsAddress, "https server listen address")
	flag.StringVar(&HttpsDomain, "https-domain", HttpsDomain, "https server domain")
	flag.StringVar(&HttpsEmail, "https-email", HttpsEmail, "https cert email")
	flag.StringVar(&HttpsCertDir, "https-cert-dir", HttpsCertDir, "https cert save dir")

	flag.StringVar(&HttpsAliyunKey, "https-aliyun-dns-access-key", HttpsAliyunKey, "aliyun access key")
	flag.StringVar(&HttpsAliyunSecret, "https-aliyun-dns-access-secret", HttpsAliyunSecret, "aliyun access secret")

	flag.BoolVar(&DryRun, "dry-run", DryRun, "only parser the options")

	flag.BoolVar(&Version, "version", Version, "show the version")
	flag.BoolVar(&Version, "v", Version, "show the version")

	flag.BoolVar(&License, "license", License, "show the license")
	flag.BoolVar(&License, "l", License, "show the license")

	flag.BoolVar(&Report, "report", Report, "show the report")
	flag.BoolVar(&Report, "r", Report, "show the report")

	flag.BoolVar(&ShowOption, "show-option", ShowOption, "show the option")
	flag.BoolVar(&ShowOption, "s", ShowOption, "show the option")

	flag.Parse()

	return nil
}

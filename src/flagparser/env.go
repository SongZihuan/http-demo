package flagparser

import "os"

const ENV_PREFIX = "DH_"
const HTTP_PREFIX = "HTTP_"
const HTTPS_PREFIX = "HTTPS_"
const ALIYUN_PREFIX = "ALIYUN_"

const (
	HTTP_ADDRESS = ENV_PREFIX + HTTP_PREFIX + "ADDRESS"

	HTTPS_ADDRESS       = ENV_PREFIX + HTTPS_PREFIX + "ADDRESS"
	HTTPS_DOMAIN        = ENV_PREFIX + HTTPS_PREFIX + "DOMAIN"
	HTTPS_EMAIL         = ENV_PREFIX + HTTPS_PREFIX + "EMAIL"
	HTTPS_CERT_DIR      = ENV_PREFIX + HTTPS_PREFIX + "CERT_DIR"
	HTTPS_ALIYUN_KEY    = ENV_PREFIX + HTTPS_PREFIX + ALIYUN_PREFIX + "KEY"
	HTTPS_ALIYUN_SECRET = ENV_PREFIX + HTTPS_PREFIX + ALIYUN_PREFIX + "SECRET"
)

func initEnv() error {
	_HttpAddress := os.Getenv(HTTP_ADDRESS)
	if _HttpAddress != "" {
		HttpAddress = _HttpAddress
	}
	HttpsAddress = os.Getenv(HTTPS_ADDRESS)
	HttpsDomain = os.Getenv(HTTPS_DOMAIN)
	HttpsEmail = os.Getenv(HTTPS_EMAIL)
	_HttpsCertDir := os.Getenv(HTTPS_CERT_DIR)
	if _HttpsCertDir != "" {
		HttpsCertDir = _HttpsCertDir
	}
	HttpsAliyunKey = os.Getenv(HTTPS_ALIYUN_KEY)
	HttpsAliyunSecret = os.Getenv(HTTPS_ALIYUN_SECRET)
	return nil
}

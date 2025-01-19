package flagparser

import (
	"flag"
	"fmt"
	resource "github.com/SongZihuan/http-demo"
	"github.com/SongZihuan/http-demo/src/utils"
	"strings"
)

func PrintLicense() (int, error) {
	title := utils.FormatTextToWidth(fmt.Sprintf("License of %s:", utils.GetArgs0Name()), utils.NormalConsoleWidth)
	license := utils.FormatTextToWidth(resource.License, utils.NormalConsoleWidth)
	return fmt.Fprintf(flag.CommandLine.Output(), "%s\n%s\n", title, license)
}

func PrintVersion() (int, error) {
	version := utils.FormatTextToWidth(fmt.Sprintf("Version of %s: %s", utils.GetArgs0Name(), resource.Version), utils.NormalConsoleWidth)
	return fmt.Fprintf(flag.CommandLine.Output(), "%s\n", version)
}

func PrintReport() (int, error) {
	// 不需要title
	report := utils.FormatTextToWidth(resource.Report, utils.NormalConsoleWidth)
	return fmt.Fprintf(flag.CommandLine.Output(), "%s\n", report)
}

func PrintLF() (int, error) {
	return fmt.Fprintf(flag.CommandLine.Output(), "\n")
}

func Print() {
	fmt.Println("HttpAddress:", HttpAddress)
	fmt.Println("HttpsAddress:", HttpsAddress)
	fmt.Println("HttpsDomain:", HttpsDomain)
	fmt.Println("HttpsEmail:", HttpsEmail)
	fmt.Println("HttpsCertDir:", HttpsCertDir)
	fmt.Println("HttpsAliyunKey:", HttpsAliyunKey)
	fmt.Println("HttpsAliyunSecret:", strings.Repeat("*", len(HttpsAliyunSecret)))
}

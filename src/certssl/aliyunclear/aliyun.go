package aliyunclear

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"strconv"
	"strings"
)

var domainList []string
var mainDomain string
var rr string

var aliyunAccessKey string
var aliyunAccessSecret string

func InitAliyun(_aliyunAccessKey string, _aliyunAccessSecret string, domain string) (err error) {
	if aliyunAccessKey != "" {
		return nil
	}

	err = aliyunSDKCreateClient(_aliyunAccessKey, _aliyunAccessSecret)
	if err != nil {
		return fmt.Errorf("create aliyun client failed: %s", err.Error())
	}

	aliyunAccessKey = _aliyunAccessKey
	aliyunAccessSecret = _aliyunAccessSecret

	domainList, err = getDomainList()
	if err != nil {
		return fmt.Errorf("get domain list failed: %s", err.Error())
	}

	mainDomain, rr, err = getMainDomain(domain)
	if err != nil {
		return fmt.Errorf("get main domain failed: %s", err.Error())
	}

	return nil
}

func getDomainList() ([]string, error) {
	const pageSize = int64(100)
	var pageNumber = int64(0)

	res := make([]string, 0, 10)

	for {
		pageNumber += 1
		resp1, err := aliyunSDKDescribeDomains(pageNumber, pageSize)
		if err != nil {
			return nil, fmt.Errorf("get domain list failed: %s", err.Error())
		}

		for _, d := range resp1.Body.Domains.Domain {
			res = append(res, tea.StringValue(d.DomainName))
		}

		if int64(len(res)) == tea.Int64Value(resp1.Body.TotalCount) {
			break
		}
	}

	return res, nil
}

func getMainDomain(subDomainName string) (string, string, error) {
	for _, domainName := range domainList {
		suffix := fmt.Sprintf(".%s", domainName)

		if domainName == subDomainName {
			RR := "@"
			return domainName, RR, nil
		} else if strings.HasSuffix(subDomainName, suffix) {
			RR := strings.TrimSuffix(subDomainName, suffix)
			return domainName, RR, nil
		}
	}

	return "", "", fmt.Errorf("domain not found")
}

func clearACMEDNS01TXTRecord() error {
	resp1, err := aliyunSDKDeleteSubDomainRecords(mainDomain, rr, "TXT")
	if err != nil {
		return err
	}

	totalCount, err := strconv.ParseInt(tea.StringValue(resp1.Body.TotalCount), 10, 64)
	if err != nil {
		return err
	}

	if totalCount <= 0 {
		return fmt.Errorf("clear failed: total count is 0")
	} else if totalCount != 1 {
		return fmt.Errorf("clear may have some problem: total count is %d", totalCount)
	}

	return nil
}

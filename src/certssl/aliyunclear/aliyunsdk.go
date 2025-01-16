package aliyunclear

import (
	dns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var client *dns.Client

func aliyunSDKCreateClient(aliyunAccessKey string, aliyunAccessSecret string) (err error) {
	client, err = _aliyunSDKCreateClient(aliyunAccessKey, aliyunAccessSecret)
	if err != nil {
		return err
	}
	return nil
}

func _aliyunSDKCreateClient(aliyunAccessKey string, aliyunAccessSecret string) (*dns.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(aliyunAccessKey),
		AccessKeySecret: tea.String(aliyunAccessSecret),
	}
	config.RegionId = tea.String("cn-hangzhou")
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	res := &dns.Client{}
	res, err := dns.NewClient(config)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func aliyunSDKDescribeDomains(pageNumber int64, pageSize int64) (*dns.DescribeDomainsResponse, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}

	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}

	req := &dns.DescribeDomainsRequest{}
	req.PageNumber = tea.Int64(pageNumber)
	req.PageSize = tea.Int64(pageSize)
	resp, tryErr := func() (resp *dns.DescribeDomainsResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		resp, err = client.DescribeDomains(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}()
	if tryErr != nil {
		return nil, tryErr
	}
	return resp, nil
}

func aliyunSDKDeleteSubDomainRecords(domainName string, RR string, Type string) (*dns.DeleteSubDomainRecordsResponse, error) {
	req := &dns.DeleteSubDomainRecordsRequest{}
	req.DomainName = tea.String(domainName)
	req.RR = tea.String(RR)
	req.Type = tea.String(Type)
	resp, tryErr := func() (resp *dns.DeleteSubDomainRecordsResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		resp, err = client.DeleteSubDomainRecords(req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}()
	if tryErr != nil {
		return nil, tryErr
	}
	return resp, nil
}

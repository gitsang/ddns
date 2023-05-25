package main

import (
	"ddns/pkg/api"
	"encoding/json"
	"fmt"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	AccessKeyId     = ""
	AccessKeySecret = ""

	Domain = "home.c8g.top"
	RR     = "@"
	Type   = "AAAA"
)

func Describe(client *alidns20150109.Client) {
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(Domain),
		RRKeyWord:  tea.String(RR),
		Type:       tea.String(Type),
	}
	resp, err := client.DescribeDomainRecords(describeDomainRecordsRequest)
	if err != nil {
		panic(err)
	}

	respBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("resp", string(respBytes))
}

func Find(client *alidns20150109.Client) {
	rec, err := api.FindRecord(client, Domain, RR, Type)
	if err != nil {
		panic(err)
	}
	recBytes, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("rec", string(recBytes))
}

func main() {
	client, err := api.CreateClient(AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
	}

	Describe(client)

	Find(client)
}

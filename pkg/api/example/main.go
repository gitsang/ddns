package main

import (
	"ddns/pkg/api"
	"encoding/json"
	"fmt"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	AccessKeyId     = "LTAI5t5dPz2s7CAbpJNaGKBj"
	AccessKeySecret = "klbtdIHUzrCV3LxJ5LEZ7aUr5y3698"
)

func main() {
	client, err := api.CreateClient(AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
	}

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName:  tea.String("home.c8g.top"),
		RRKeyWord:   tea.String("@"),
		TypeKeyWord: tea.String("A"),
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

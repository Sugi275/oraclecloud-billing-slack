package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// OracleBillingClient OracleBillingClient
type OracleBillingClient struct {
	URL                    string
	UserName               string
	Password               string
	IdentityCloudServiceID string
}

func main() {
	url := "https://itra.oraclecloud.com/metering/api/v1/usagecost/cacct-8da51d892ccf453cae8c82145fcbc345?startTime=2019-03-26T00:00:00.000Z&endTime=2019-03-27T00:00:00.000Z&timeZone=UTC&usageType=DAILY"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-ID-TENANT-NAME", "idcs-4833af834c9a44e0b20738c3029777a6")
	req.SetBasicAuth("api", "gojs90jgr908ydsio'&(IOHi")
	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	buf := bytes.NewBuffer(body)
	html := buf.String()
	fmt.Println(html)
}

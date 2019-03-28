package oraclecloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	envOraclecloudAccountID = "ORACLECLOUD_ACCOUNT_ID"
	envOraclecloudIDCSID    = "ORACLECLOUD_IDCS_ID"
	envOraclecloudUsername  = "ORACLECLOUD_USERNAME"
	envOraclecloudPassword  = "ORACLECLOUD_PASSWORD"
)

// OracleBillingClient OracleBillingClient
type OracleBillingClient struct {
	TodayURL               string
	MonthURL               string
	UserName               string
	Password               string
	IdentityCloudServiceID string
}

// OracleBillingData OracleBillingData
type OracleBillingData struct {
	ServiceName  string
	ResourceName string
	Billing      int
	Currency     string
}

// OracleBillingDataList OracleBillingDataList
type OracleBillingDataList struct {
	TodayBilling   []OracleBillingData
	MonthBilling   []OracleBillingData
	TodayStartDate time.Time
	TodayEndDate   time.Time
	MonthStartDate time.Time
	MonthEndDate   time.Time
}

// GetOracleBillingData GetOracleBillingData
func GetOracleBillingData() (OracleBillingDataList, error) {
	oracleBillingDataList := OracleBillingDataList{}
	return oracleBillingDataList, nil
}

// getOracleBillingClient getOracleBillingClient
func getOracleBillingClient() (OracleBillingClient, error) {
	client := OracleBillingClient{}

	oraclecloudAccountID, ok := os.LookupEnv(envOraclecloudAccountID)
	if !ok {
		err := fmt.Errorf("can not read envOraclecloudAccountID from environment variable %s", envOraclecloudAccountID)
		loglib.Sugar.Error(err)
		return nil, err
	}

	client.TodayURL = "https://itra.oraclecloud.com/metering/api/v1/usagecost/cacct-8da51d892ccf453cae8c82145fcbc345?startTime=2019-03-26T00:00:00.000Z&endTime=2019-03-27T00:00:00.000Z&timeZone=UTC&usageType=DAILY"

	return client, nil
}

func testmain() {
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

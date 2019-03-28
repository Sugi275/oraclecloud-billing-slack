package oraclecloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sugi275/oraclecloud-billing-slack/src/loglib"
)

const (
	envOraclecloudAccountID = "ORACLECLOUD_ACCOUNT_ID"
	envOraclecloudIDCSID    = "ORACLECLOUD_IDCS_ID"
	envOraclecloudUsername  = "ORACLECLOUD_USERNAME"
	envOraclecloudPassword  = "ORACLECLOUD_PASSWORD"
)

// OracleBillingClient OracleBillingClient
type OracleBillingClient struct {
	YesterdayURL       string
	MonthURL           string
	UserName           string
	Password           string
	AccountID          string
	IDCSID             string
	YesterdayStartDate time.Time
	YesterdayEndDate   time.Time
	MonthStartDate     time.Time
	MonthEndDate       time.Time
}

// OracleBillingData OracleBillingData
type OracleBillingData struct {
	ServiceName  string
	ResourceName string
	Billing      float64
	Currency     string
	GsiProductID string
}

// OracleBillingDataList OracleBillingDataList
type OracleBillingDataList struct {
	YesterdayBilling   []OracleBillingData
	MonthBilling       []OracleBillingData
	YesterdayTotal     float64
	MonthTotal         float64
	Currency           string
	YesterdayStartDate time.Time
	YesterdayEndDate   time.Time
	MonthStartDate     time.Time
	MonthEndDate       time.Time
	BillingPageURL     string
}

// OracleResponseJSON OracleCloudのMeterAPIから取得するJSONを定義する構造体
type OracleResponseJSON struct {
	AccountID string `json:"accountId"`
	Items     []struct {
		SubscriptionID       string `json:"subscriptionId"`
		SubscriptionType     string `json:"subscriptionType"`
		ServiceName          string `json:"serviceName"`
		ResourceName         string `json:"resourceName"`
		Currency             string `json:"currency"`
		GsiProductID         string `json:"gsiProductId"`
		ServiceEntitlementID string `json:"serviceEntitlementId"`
		Costs                []struct {
			ComputedQuantity float64 `json:"computedQuantity"`
			ComputedAmount   float64 `json:"computedAmount"`
			UnitPrice        int     `json:"unitPrice"`
			OveragesFlag     string  `json:"overagesFlag"`
		} `json:"costs"`
		StartTimeUtc string `json:"startTimeUtc"`
		EndTimeUtc   string `json:"endTimeUtc"`
	} `json:"items"`
	CanonicalLink string `json:"canonicalLink"`
}

// GetOracleBillingData GetOracleBillingData
func GetOracleBillingData() (OracleBillingDataList, error) {
	oracleBillingDataList := OracleBillingDataList{}

	client, err := getOracleBillingClient()
	if err != nil {
		return oracleBillingDataList, err
	}

	oracleBillingDataList, err = getOracleBillingList(client)
	oracleBillingDataList.Currency = "JPY"
	oracleBillingDataList.YesterdayStartDate = client.YesterdayStartDate
	oracleBillingDataList.YesterdayEndDate = client.YesterdayEndDate
	oracleBillingDataList.MonthStartDate = client.MonthStartDate
	oracleBillingDataList.MonthEndDate = client.MonthEndDate

	return oracleBillingDataList, nil
}

// getOracleBillingClient getOracleBillingClient
func getOracleBillingClient() (OracleBillingClient, error) {
	client := OracleBillingClient{}

	oraclecloudAccountID, ok := os.LookupEnv(envOraclecloudAccountID)
	if !ok {
		err := fmt.Errorf("can not read envOraclecloudAccountID from environment variable %s", envOraclecloudAccountID)
		loglib.Sugar.Error(err)
		return client, err
	}
	client.AccountID = oraclecloudAccountID

	oraclecloudIDSCID, ok := os.LookupEnv(envOraclecloudIDCSID)
	if !ok {
		err := fmt.Errorf("can not read envOraclecloudIDCSID from environment variable %s", envOraclecloudIDCSID)
		loglib.Sugar.Error(err)
		return client, err
	}
	client.IDCSID = oraclecloudIDSCID

	oraclecloudUsername, ok := os.LookupEnv(envOraclecloudUsername)
	if !ok {
		err := fmt.Errorf("can not read envOraclecloudUsername from environment variable %s", envOraclecloudUsername)
		loglib.Sugar.Error(err)
		return client, err
	}
	client.UserName = oraclecloudUsername

	oraclecloudPassword, ok := os.LookupEnv(envOraclecloudPassword)
	if !ok {
		err := fmt.Errorf("can not read envOraclecloudPassword from environment variable %s", envOraclecloudPassword)
		loglib.Sugar.Error(err)
		return client, err
	}
	client.Password = oraclecloudPassword

	client.YesterdayStartDate = getYesterdayStartDate()
	client.YesterdayEndDate = getYesterdayEndDate()
	client.MonthStartDate = getMonthStartDate()
	client.MonthEndDate = getMonthEndDate()

	baseURL := "https://itra.oraclecloud.com/metering/api/v1/usagecost/"

	client.YesterdayURL = baseURL + oraclecloudIDSCID + getURLDateParameter(client.YesterdayStartDate, client.YesterdayEndDate)
	client.MonthURL = baseURL + oraclecloudIDSCID + getURLDateParameter(client.MonthStartDate, client.MonthEndDate)

	return client, nil
}

func getYesterdayStartDate() time.Time {
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayStartDate := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	return yesterdayStartDate
}

func getYesterdayEndDate() time.Time {
	now := time.Now()
	yesterdayEndDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return yesterdayEndDate
}

func getMonthStartDate() time.Time {
	now := time.Now()
	monthStartDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	return monthStartDate
}

func getMonthEndDate() time.Time {
	now := time.Now()
	monthStartDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEndDate := monthStartDate.AddDate(0, 1, -1) //来月の初日から1日を引くと、今月の月末
	return monthEndDate
}

func getURLDateParameter(start time.Time, end time.Time) string {
	parameter := "?startTime=" + start.Format("2006-01-02T15:04:05.000Z") + "&endTime=" + end.Format("2006-01-02T15:04:05.000Z") + "&timeZone=UTC&usageType=TOTAL"
	return parameter
}

func getOracleBillingList(client OracleBillingClient) (OracleBillingDataList, error) {
	oracleBillingDataList := OracleBillingDataList{}
	oracleResponseJSON, err := getOracleBillingData(client, "YESTERDAY")
	if err != nil {
		return oracleBillingDataList, err
	}

	for _, billingItem := range oracleResponseJSON.Items {
		oracleBillingData := OracleBillingData{
			ServiceName:  billingItem.ServiceName,
			ResourceName: billingItem.ResourceName,
			Billing:      billingItem.Costs[0].ComputedAmount,
			GsiProductID: billingItem.GsiProductID,
			Currency:     "JPY",
		}
		oracleBillingDataList.YesterdayBilling = append(oracleBillingDataList.YesterdayBilling, oracleBillingData)
		oracleBillingDataList.YesterdayTotal = oracleBillingDataList.YesterdayTotal + oracleBillingData.Billing
	}

	oracleResponseJSON, err = getOracleBillingData(client, "MONTH")
	if err != nil {
		return oracleBillingDataList, err
	}

	for _, billingItem := range oracleResponseJSON.Items {
		oracleBillingData := OracleBillingData{
			ServiceName:  billingItem.ServiceName,
			ResourceName: billingItem.ResourceName,
			Billing:      billingItem.Costs[0].ComputedAmount,
			GsiProductID: billingItem.GsiProductID,
			Currency:     "JPY",
		}
		oracleBillingDataList.MonthBilling = append(oracleBillingDataList.MonthBilling, oracleBillingData)
		oracleBillingDataList.MonthTotal = oracleBillingDataList.MonthTotal + oracleBillingData.Billing
	}

	oracleBillingDataList.BillingPageURL = generateBillingPageURL(client)

	return oracleBillingDataList, nil
}

func getOracleBillingData(oracleBillingClient OracleBillingClient, duration string) (OracleResponseJSON, error) {
	var billingResponse OracleResponseJSON

	var url string
	switch duration {
	case "YESTERDAY":
		url = oracleBillingClient.YesterdayURL
	case "MONTH":
		url = oracleBillingClient.MonthURL
	default:
		err := fmt.Errorf("duration is invalid. please select YESTERDAY or MONTH")
		return billingResponse, err
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-ID-TENANT-NAME", oracleBillingClient.AccountID)
	req.SetBasicAuth(oracleBillingClient.UserName, oracleBillingClient.Password)
	client := new(http.Client)
	res, err := client.Do(req)

	if err != nil {
		return billingResponse, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return billingResponse, err
	}

	err = json.Unmarshal(body, &billingResponse)

	return billingResponse, nil
}

func generateBillingPageURL(oracleBillingClient OracleBillingClient) string {
	return "https://myservices-" + oracleBillingClient.IDCSID + ".console.oraclecloud.com/mycloud/cloudportal/accountDetail"
}

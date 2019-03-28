package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"unicode/utf8"

	"github.com/Sugi275/oraclecloud-billing-slack/src/loglib"
	"github.com/Sugi275/oraclecloud-billing-slack/src/oraclecloud"
)

const (
	envSlackIncomingURL = "SLACK_INCOMING_URL"
)

// slackMessage slackMessage
type slackMessage struct {
	Attachments []Attachments `json:"attachments"`
}

// Attachments Attachments
type Attachments struct {
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	Pretext    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	Fields     []Fields `json:"fields"`
	ImageURL   string   `json:"image_url"`
	ThumbURL   string   `json:"thumb_url"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
	Ts         int      `json:"ts"`
}

// Fields Fields
type Fields struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// PostBilling PostBilling
func PostBilling(billing oraclecloud.OracleBillingDataList) error {
	slackIncomingURL, ok := os.LookupEnv(envSlackIncomingURL)
	if !ok {
		err := fmt.Errorf("can not read envSlackIncomingURL from environment variable %s", envSlackIncomingURL)
		loglib.Sugar.Error(err)
		return err
	}

	attachments := generateAttachments(billing)
	slackMessage := slackMessage{
		Attachments: attachments,
	}

	slackJSON, err := json.Marshal(&slackMessage)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", slackIncomingURL, bytes.NewBuffer([]byte(slackJSON)))
	req.Header.Set("Content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func generateAttachments(billing oraclecloud.OracleBillingDataList) []Attachments {
	dateFormat := "2006-01-02 15:04:05"

	attachments := []Attachments{
		Attachments{
			Pretext: "<" + billing.BillingPageURL + "|OracleCloud BillingDetailPage>\n*OneDay*\n" +
				billing.YesterdayStartDate.Format(dateFormat) + "(UTC) ~ " + billing.YesterdayEndDate.Format(dateFormat) + "(UTC)",
			Color:  "#7CD197",
			Fields: generateFields(billing, "YESTERDAY"),
		}, {
			Pretext: "*MonthDay*\n" +
				billing.MonthStartDate.Format(dateFormat) + "(UTC) ~ " + billing.MonthEndDate.Format(dateFormat) + "(UTC)",
			Color:  "#7CD197",
			Fields: generateFields(billing, "MONTH"),
		},
	}

	return attachments
}

func generateFields(billing oraclecloud.OracleBillingDataList, duration string) []Fields {
	resourceValue := "Total\n"
	var billingValue string
	var billingSlice []oraclecloud.OracleBillingData

	switch duration {
	case "YESTERDAY":
		billingValue = strconv.FormatFloat(billing.YesterdayTotal, 'f', 1, 64) + " " + billing.Currency + "\n"
		billingSlice = billing.YesterdayBilling
	case "MONTH":
		billingValue = strconv.FormatFloat(billing.MonthTotal, 'f', 1, 64) + " " + billing.Currency + "\n"
		billingSlice = billing.MonthBilling
	}

	for _, o := range billingSlice {
		if utf8.RuneCountInString(o.ResourceName) > 28 {
			o.ResourceName = o.ResourceName[:28]
		}
		resourceValue = resourceValue + o.ResourceName + "\n" // slack上の表示欄が小さいため、28文字を上限にしている

		billingValue = billingValue + strconv.FormatFloat(o.Billing, 'f', 1, 64) +
			" " + o.Currency + " (GsiProductId: " + o.GsiProductID + ")" + "\n"
	}

	fields := []Fields{
		Fields{
			Title: "Resource",
			Value: resourceValue,
			Short: true,
		}, {
			Title: "Billing",
			Value: billingValue,
			Short: true,
		},
	}

	return fields
}

package main

import (
	"fmt"

	"github.com/Sugi275/oraclecloud-billing-slack/src/loglib"
	"github.com/Sugi275/oraclecloud-billing-slack/src/oraclecloud"
	"github.com/Sugi275/oraclecloud-billing-slack/src/slack"
)

func main() {
	oracleBillingData, err := oraclecloud.GetOracleBillingData()
	if err != nil {
		loglib.Sugar.Error(err)
		return
	}

	err = slack.PostBilling(oracleBillingData)
	if err != nil {
		loglib.Sugar.Error(err)
		return
	}

	for _, item := range oracleBillingData.YesterdayBilling {
		println(item.ServiceName)
		println(item.ResourceName)
		fmt.Printf("%.4f\n", item.Billing)
		println(item.Billing)
		println(item.Currency)
		println()
	}
}

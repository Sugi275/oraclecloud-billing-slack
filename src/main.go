package main

import (
	"github.com/Sugi275/oraclecloud-billing-slack/src/loglib"
	"github.com/Sugi275/oraclecloud-billing-slack/src/oraclecloud"
	"github.com/Sugi275/oraclecloud-billing-slack/src/slack"
)

func main() {
	loglib.InitSugar()
	defer loglib.Sugar.Sync()

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
}

package main

import (
	"context"
	"io"

	"github.com/Sugi275/oraclecloud-billing-slack/src/loglib"
	"github.com/Sugi275/oraclecloud-billing-slack/src/oraclecloud"
	"github.com/Sugi275/oraclecloud-billing-slack/src/slack"
	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(oraclecloudBillingSlackHangler))
}

func oraclecloudBillingSlackHangler(ctx context.Context, in io.Reader, out io.Writer) {
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

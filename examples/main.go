package main

import (
	"context"
	"time"

	"github.com/jony4/ezpay"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := &ezpay.Config{
		ProjectID:     "1",
		ProjectSecret: "2",
		PaywallID:     "3",
	}
	client, err := ezpay.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	customerInfo, err := client.CustomerInfo(ctx, &ezpay.EZCustomerInfoRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(customerInfo)
}

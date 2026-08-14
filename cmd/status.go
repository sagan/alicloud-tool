package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display current Alibaba Cloud account status info (CDT traffic, Billing)",
	Long:  `Queries and outputs the current month's internet CDT traffic used and the billing invoice amounts for the last 3 months.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching Alibaba Cloud status...")

		// CDT Traffic
		totalGB, err := getTotalTrafficGB()
		if err != nil {
			fmt.Printf("❌ Failed to fetch CDT traffic: %v\n", err)
		} else {
			fmt.Printf("✅ Current Total Internet Traffic: %.2f GB\n", totalGB)
		}

		// Billing Info
		fmt.Println("\nFetching Billing Invoices...")
		bssClient, err := getBssClient()
		if err != nil {
			fmt.Printf("❌ Failed to initialize BSS client: %v\n", err)
			return
		}

		now := time.Now()
		for i := range 3 {
			targetDate := now.AddDate(0, -i, 0)
			billingCycle := targetDate.Format("2006-01")

			amount, err := fetchMonthlyInvoiceAmount(bssClient, billingCycle)
			if err != nil {
				fmt.Printf("  - %s: ❌ Error retrieving invoice: %v\n", billingCycle, err)
			} else {
				fmt.Printf("  - %s: $%.2f\n", billingCycle, amount)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getBssClient() (*openapi.Client, error) {
	config, err := createConfig()
	if err != nil {
		return nil, err
	}
	if viper.GetBool("ALIBABA_CLOUD_CN") {
		config.Endpoint = tea.String("business.aliyuncs.com")
	} else {
		config.Endpoint = tea.String("business.ap-southeast-1.aliyuncs.com")
	}
	return openapi.NewClient(config)
}

func fetchMonthlyInvoiceAmount(client *openapi.Client, billingCycle string) (float64, error) {
	params := &openapi.Params{
		Action:      tea.String("QueryBillOverview"),
		Version:     tea.String("2017-12-14"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	request := &openapi.OpenApiRequest{
		Query: map[string]*string{
			"BillingCycle": tea.String(billingCycle),
		},
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.CallApi(params, request, runtime)
	if err != nil {
		return 0, err
	}

	bodyMap, ok := resp["body"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected response body format")
	}

	dataMap, ok := bodyMap["Data"].(map[string]interface{})
	if !ok || dataMap == nil {
		return 0, nil // No data / Not billed yet
	}

	itemsMap, ok := dataMap["Items"].(map[string]interface{})
	if !ok || itemsMap == nil {
		return 0, nil
	}

	itemList, ok := itemsMap["Item"].([]interface{})
	if !ok || len(itemList) == 0 {
		return 0, nil
	}

	var total float64
	for _, rawItem := range itemList {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}

		pretaxRaw, valOk := item["PretaxAmount"]
		if !valOk {
			continue
		}

		switch v := pretaxRaw.(type) {
		case float64:
			total += v
		case json.Number:
			if amount, err := v.Float64(); err == nil {
				total += amount
			}
		}
	}

	return total, nil
}

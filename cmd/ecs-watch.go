package cmd

import (
	"fmt"
	"log"
	"math"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ecsWatchCmd = &cobra.Command{
	Use:   "ecs-watch",
	Short: "Watch CDT traffic and start/stop ECS",
	Long:  `Monitors CDT traffic and starts or stops the target ECS instance based on the traffic threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting traffic monitor...")

		threshold := viper.GetFloat64("ALIBABA_CLOUD_THRESHOLD")
		instanceID := viper.GetString("ALIBABA_CLOUD_INSTANCE_ID")

		if instanceID == "" {
			log.Fatalf("Error: Instance ID must be specified (via flag --instance-id or env ALIBABA_CLOUD_INSTANCE_ID).")
		}

		totalGB, err := getTotalTrafficGB()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		ecsClient, err := getEcsClient()
		if err != nil {
			log.Fatalf("Error initializing ECS client: %v", err)
		}

		status, err := getEcsStatus(ecsClient, instanceID)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("ECS instance %s current status: %s", instanceID, status)

		if totalGB < threshold {
			log.Printf("Traffic %.2f GB < Threshold %.2f GB, attempting to start ECS...", totalGB, threshold)
			ecsStart(ecsClient, instanceID, status)
		} else {
			log.Printf("Traffic %.2f GB >= Threshold %.2f GB, attempting to stop ECS...", totalGB, threshold)
			ecsStop(ecsClient, instanceID, status)
		}

		log.Println("Script execution completed.")
	},
}

func init() {
	rootCmd.AddCommand(ecsWatchCmd)
	ecsWatchCmd.Flags().Float64("threshold", 195.0, "Traffic threshold in GB (Env: ALIBABA_CLOUD_THRESHOLD)")
	viper.BindPFlag("ALIBABA_CLOUD_THRESHOLD", ecsWatchCmd.Flags().Lookup("threshold"))
}

func getTotalTrafficGB() (float64, error) {
	config, err := createConfig()
	if err != nil {
		return 0, err
	}
	config.Endpoint = tea.String("cdt.aliyuncs.com")

	client, err := openapi.NewClient(config)
	if err != nil {
		return 0, fmt.Errorf("failed to init CDT client: %v", err)
	}

	params := &openapi.Params{
		Action:      tea.String("ListCdtInternetTraffic"),
		Version:     tea.String("2021-08-13"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	request := &openapi.OpenApiRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.CallApi(params, request, runtime)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch CDT traffic: %v", err)
	}

	bodyMap, ok := resp["body"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("unexpected response body format")
	}

	trafficDetails, ok := bodyMap["TrafficDetails"].([]interface{})
	if !ok {
		return 0, nil
	}

	var totalBytes float64
	for _, detail := range trafficDetails {
		detailMap, ok := detail.(map[string]interface{})
		if !ok {
			continue
		}
		if traffic, valOk := detailMap["Traffic"].(float64); valOk {
			totalBytes += traffic
		}
	}

	totalGB := totalBytes / math.Pow(1024, 3)
	log.Printf("Current Total Internet Traffic: %.2f GB", totalGB)
	return totalGB, nil
}

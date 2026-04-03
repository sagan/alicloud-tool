package cmd

import (
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs "github.com/alibabacloud-go/ecs-20140526/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "alicloud-tool",
	Short: "A tool to manage and monitor Alibaba Cloud ECS and CDT",
	Long:  `A command line tool specifically designed to start/stop ECS instances based on CDT traffic thresholds.

Note: All settings can also be configured using environment variables.
The corresponding environment variables are prefixed with "ALIBABA_CLOUD_"
and use underscores instead of hyphens (e.g., ALIBABA_CLOUD_ACCESS_KEY_ID).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("access-key-id", "", "Alibaba Cloud Access Key ID (Env: ALIBABA_CLOUD_ACCESS_KEY_ID)")
	rootCmd.PersistentFlags().String("access-key-secret", "", "Alibaba Cloud Access Key Secret (Env: ALIBABA_CLOUD_ACCESS_KEY_SECRET)")
	rootCmd.PersistentFlags().String("region", "cn-hongkong", "Alibaba Cloud Region (Env: ALIBABA_CLOUD_REGION)")
	rootCmd.PersistentFlags().String("instance-id", "", "Alibaba Cloud ECS Instance ID (Env: ALIBABA_CLOUD_INSTANCE_ID)")

	viper.BindPFlag("ALIBABA_CLOUD_ACCESS_KEY_ID", rootCmd.PersistentFlags().Lookup("access-key-id"))
	viper.BindPFlag("ALIBABA_CLOUD_ACCESS_KEY_SECRET", rootCmd.PersistentFlags().Lookup("access-key-secret"))
	viper.BindPFlag("ALIBABA_CLOUD_REGION", rootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("ALIBABA_CLOUD_INSTANCE_ID", rootCmd.PersistentFlags().Lookup("instance-id"))

	viper.AutomaticEnv() // Read from environment variables
}

func initConfig() {
}

func createConfig() (*openapi.Config, error) {
	accessKeyID := viper.GetString("ALIBABA_CLOUD_ACCESS_KEY_ID")
	accessKeySecret := viper.GetString("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	region := viper.GetString("ALIBABA_CLOUD_REGION")

	if accessKeyID == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("please set ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET environment variables or flags")
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		RegionId:        tea.String(region),
	}
	return config, nil
}

func getEcsClient() (*ecs.Client, error) {
	config, err := createConfig()
	if err != nil {
		return nil, err
	}
	region := viper.GetString("ALIBABA_CLOUD_REGION")
	config.Endpoint = tea.String(fmt.Sprintf("ecs.%s.aliyuncs.com", region))
	return ecs.NewClient(config)
}

func getEcsStatus(client *ecs.Client, instanceID string) (string, error) {
	region := viper.GetString("ALIBABA_CLOUD_REGION")
	req := &ecs.DescribeInstancesRequest{
		InstanceIds: tea.String(fmt.Sprintf(`["%s"]`, instanceID)),
		RegionId:    tea.String(region),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.DescribeInstancesWithOptions(req, runtime)
	if err != nil {
		return "", fmt.Errorf("failed to get ECS status: %v", err)
	}

	if resp.Body == nil || resp.Body.Instances == nil || len(resp.Body.Instances.Instance) == 0 {
		return "", fmt.Errorf("ECS instance not found")
	}

	status := tea.StringValue(resp.Body.Instances.Instance[0].Status)
	return status, nil
}

func ecsStart(client *ecs.Client, instanceID string, currentStatus string) {
	if currentStatus == "Running" {
		fmt.Printf("ECS instance %s is already running. No action needed.\n", instanceID)
		return
	}

	req := &ecs.StartInstanceRequest{
		InstanceId: tea.String(instanceID),
	}
	runtime := &util.RuntimeOptions{}

	_, err := client.StartInstanceWithOptions(req, runtime)
	if err != nil {
		fmt.Printf("Failed to start ECS instance: %v\n", err)
		return
	}
	fmt.Println("ECS start request triggered successfully.")
}

func ecsStop(client *ecs.Client, instanceID string, currentStatus string) {
	if currentStatus == "Stopped" {
		fmt.Printf("ECS instance %s is already stopped. No action needed.\n", instanceID)
		return
	}

	req := &ecs.StopInstanceRequest{
		InstanceId: tea.String(instanceID),
		ForceStop:  tea.Bool(false),
	}
	runtime := &util.RuntimeOptions{}

	_, err := client.StopInstanceWithOptions(req, runtime)
	if err != nil {
		fmt.Printf("Failed to stop ECS instance: %v\n", err)
		return
	}
	fmt.Println("ECS stop request triggered successfully.")
}

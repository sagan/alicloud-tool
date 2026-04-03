package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ecsStopCmd = &cobra.Command{
	Use:   "ecs-stop",
	Short: "Stop ECS instance manually",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID := viper.GetString("ALIBABA_CLOUD_INSTANCE_ID")
		if instanceID == "" {
			log.Fatalf("Error: Instance ID must be specified (via flag --instance-id or env ALIBABA_CLOUD_INSTANCE_ID).")
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
		ecsStop(ecsClient, instanceID, status)
	},
}

func init() {
	rootCmd.AddCommand(ecsStopCmd)
}

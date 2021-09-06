package cmd

import (
	"context"

	"github.com/matt-potter/scaling-demo-app/populator"
	"github.com/spf13/cobra"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate SQS queue with fake data",
	Run:   runPopulate,
}

var (
	populateRegion string
	populateQueue  string
)

func init() {
	rootCmd.AddCommand(populateCmd)

	populateCmd.Flags().StringVarP(&populateRegion, "region", "r", "", "AWS region")

	populateCmd.MarkFlagRequired("region")

	populateCmd.Flags().StringVarP(&populateQueue, "queue-name", "q", "", "SQS queue name")

	populateCmd.MarkFlagRequired("queue-name")
}

func runPopulate(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	p := populator.NewPopulatorClient(ctx, populateRegion)

	p.PopulateQueue(ctx, populateQueue)

}

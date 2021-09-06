package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matt-potter/scaling-demo-app/logger"
	"github.com/matt-potter/scaling-demo-app/sqspoller"
	"github.com/spf13/cobra"
)

var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Starts a SQS poller",
	Run:   runSQSHandler,
}

var (
	sqsQueue      string
	sqsRegion     string
	sqsDownstream string
)

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringVarP(&sqsQueue, "queue-name", "q", "", "SQS queue name")
	sqsCmd.Flags().StringVarP(&sqsRegion, "region", "r", "", "AWS region")
	sqsCmd.Flags().StringVarP(&sqsDownstream, "downstream", "d", "", "Downstream URL")

	sqsCmd.MarkFlagRequired("region")
	sqsCmd.MarkFlagRequired("queue-name")
	sqsCmd.MarkFlagRequired("downstream")

}

func runSQSHandler(cmd *cobra.Command, args []string) {

	ctx := context.Background()

	poller := sqspoller.NewSQSPoller(ctx, sqsRegion)

	h := &http.Client{
		Timeout: time.Second * 10,
	}

	for {
		logger.Log.Info("listening for messages...")
		poller.ProcessMessages(ctx, sqsQueue, func(body string) error {

			resp, err := h.Post(sqsDownstream, "text/plain", bytes.NewBuffer([]byte(body)))

			if err != nil {
				logger.Log.Fatal(err.Error())
			}

			rBody, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				logger.Log.Fatal(err.Error())
			}

			logger.Log.Info(fmt.Sprintf("status: %d, msg: %s", resp.StatusCode, rBody))

			return err
		})
	}
}

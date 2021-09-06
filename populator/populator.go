package populator

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/matt-potter/scaling-demo-app/logger"
	"github.com/matt-potter/scaling-demo-app/queue"
)

type Client struct {
	sqs *sqs.Client
}

func NewPopulatorClient(ctx context.Context, region string) *Client {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := sqs.NewFromConfig(cfg)

	return &Client{
		sqs: svc,
	}
}

func (c *Client) PopulateQueue(ctx context.Context, queueName string) error {

	q, err := c.sqs.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})

	if err != nil {
		logger.Log.Fatalf("unable to retrieve queue URL: %s", err.Error())
	}

	var wg sync.WaitGroup

	for i := 0; i <= 100; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()
			msg := &queue.Message{
				Body: gofakeit.Email(),
			}

			byte, err := json.Marshal(msg)

			if err != nil {
				logger.Log.Error(err)
			}

			logger.Log.Infof("message: %+v", string(byte))

			out, err := c.sqs.SendMessage(ctx, &sqs.SendMessageInput{MessageBody: aws.String(string(byte)), QueueUrl: q.QueueUrl})

			if err != nil {
				logger.Log.Error(err)
			}

			logger.Log.Infof("added message %s", *out.MessageId)
		}()

		wg.Wait()

	}

	return nil

}

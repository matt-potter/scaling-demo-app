package sqspoller

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/matt-potter/scaling-demo-app/logger"
	"github.com/matt-potter/scaling-demo-app/queue"
)

type Poller struct {
	sqs *sqs.Client
}

func NewSQSPoller(ctx context.Context, region string) *Poller {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := sqs.NewFromConfig(cfg)

	return &Poller{
		sqs: svc,
	}
}

func (p *Poller) ProcessMessages(ctx context.Context, queueName string, handle func(body string) error) error {

	q, err := p.sqs.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})

	if err != nil {
		logger.Log.Fatalf("unable to retrieve queue URL: %s", err.Error())
	}

	out, err := p.sqs.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{QueueUrl: q.QueueUrl, MaxNumberOfMessages: 10, WaitTimeSeconds: 20})

	if err != nil {
		return err
	}

	logger.Log.Infof("received %d messages", len(out.Messages))

	for _, m := range out.Messages {

		msg := queue.Message{}

		err := json.Unmarshal([]byte(*m.Body), &msg)

		if err != nil {
			logger.Log.Errorf("incorrect type for message: %s", err.Error())
		}

		err = handle(msg.Body)

		if err != nil {
			logger.Log.Errorf("error from downstream: %s", err.Error())
			continue
		}

		_, err = p.sqs.DeleteMessage(ctx, &sqs.DeleteMessageInput{QueueUrl: q.QueueUrl, ReceiptHandle: m.ReceiptHandle})

		if err != nil {
			logger.Log.Errorf("error deleting message: %s", err.Error())
		}

		logger.Log.Infof("message successfully processed: %s", *m.MessageId)

	}

	return nil

}

package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"go.uber.org/zap"
)

// SNSNotifier is an adapter that publishes events to AWS SNS
type SNSNotifier struct {
	client   *sns.Client
	topicARN string
	logger   *zap.Logger
}

// NewSNSNotifier creates a new SNSNotifier instance
func NewSNSNotifier(ctx context.Context, region, topicARN string, logger *zap.Logger) (*SNSNotifier, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create an SNS client
	client := sns.NewFromConfig(cfg)

	logger.Info("SNS notifier initialized",
		zap.String("region", region),
		zap.String("topicARN", topicARN),
	)

	return &SNSNotifier{
		client:   client,
		topicARN: topicARN,
		logger:   logger,
	}, nil
}

func (n *SNSNotifier) SendNotification(ctx context.Context, message map[string]any, eventType string) error {
	// Marshal message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		n.logger.Error("Failed to marshal message to JSON",
			zap.Error(err),
			zap.String("eventType", eventType),
		)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(n.topicARN),
		Subject:  aws.String(fmt.Sprintf("Blockchain Event: %s", eventType)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"eventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(eventType),
			},
		},
	}

	output, err := n.client.Publish(ctx, input)
	if err != nil {
		n.logger.Error("Failed to publish message to SNS",
			zap.Error(err),
			zap.String("eventType", eventType),
		)
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	n.logger.Info("Successfully published to SNS",
		zap.String("messageId", *output.MessageId),
		zap.String("eventType", eventType),
	)

	return nil
}

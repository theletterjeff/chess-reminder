package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type smsNotifier struct {
	snsClient *sns.Client
	ssmClient *ssm.Client
}

func newSMSNotifier(snsClient *sns.Client, ssmClient *ssm.Client) *smsNotifier {
	return &smsNotifier{
		snsClient: snsClient,
		ssmClient: ssmClient,
	}
}

func (n *smsNotifier) sendSMS(ctx context.Context, message string) error {
	phoneNumber, err := n.phoneNumber(ctx)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: phoneNumber,
	}

	_, err = n.snsClient.Publish(ctx, input)
	return err
}

func (n *smsNotifier) phoneNumber(ctx context.Context) (*string, error) {
	getParameterInput := &ssm.GetParameterInput{
		Name:           aws.String("chess_phone_number"),
		WithDecryption: aws.Bool(false),
	}

	parameterOutput, err := n.ssmClient.GetParameter(ctx, getParameterInput)
	if err != nil {
		return nil, err
	}

	return parameterOutput.Parameter.Value, nil
}

func message(mins int, game *game) string {
	return fmt.Sprintf(
		"You have less than %d minutes remaining in one of your games. "+
			"Move by %s in game %s",
		mins,
		time.Unix(int64(game.moveBy), 0).Format("15:04"),
		game.url,
	)
}

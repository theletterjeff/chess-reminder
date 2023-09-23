package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func handler(ctx context.Context) error {
	log.SetFlags(log.Lshortfile)

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	ssmClient := ssm.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)

	apiClient := newApiClient(ssmClient)

	resp, err := apiClient.fetchApiData(ctx)
	if err != nil {
		log.Println("error fetching chess.com API data")
		return err
	}

	notifier := newSMSNotifier(snsClient, ssmClient)
	chessReminderCfg := newChessReminderCfg()

	for _, game := range resp.games {
		log.Printf(game.url, game.moveBy)

		for _, mins := range chessReminderCfg.reminderMins {
			if game.moveBy-(mins*60) > int(time.Now().Unix()) {
				continue
			}
			if err := notifier.sendSMS(ctx, message(mins, game)); err != nil {
				log.Printf("error sending SMS for game %s", game.url)
				return err
			}
		}
	}

	return nil
}

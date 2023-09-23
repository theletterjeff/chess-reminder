package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// List of games for a player where it's their move
type apiResponse struct {
	games []*game
}

type game struct {
	url    string
	moveBy int
}

type apiClient struct {
	ssmClient *ssm.Client
}

func newApiClient(ssmClient *ssm.Client) *apiClient {
	return &apiClient{
		ssmClient: ssmClient,
	}
}

func (c *apiClient) fetchApiData(ctx context.Context) (*apiResponse, error) {
	url, err := c.toMoveURL(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data *apiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Get the API endpoint for the player's games where it's their turn
func (c *apiClient) toMoveURL(ctx context.Context) (string, error) {
	getParameterInput := &ssm.GetParameterInput{
		Name:           aws.String("chess_username"),
		WithDecryption: aws.Bool(false),
	}

	output, err := c.ssmClient.GetParameter(ctx, getParameterInput)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.chess.com/pub/player/%s/games/to-move",
		*output.Parameter.Value)

	return url, nil
}

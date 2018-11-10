
package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"net/http"
	"os"
)

type movie struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func findAll(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while retrieving AWS credentials",
		}, nil
	}

	svc := dynamodb.New(cfg)
	req := svc.ScanRequest(&dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	})

	res, err := req.Send()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while scanning DynamoDB",
		}, nil
	}

	movies := make([]movie, 0)
	for _, item := range res.Items {
		movies = append(movies, movie{
			ID:   *item["ID"].S,
			Name: *item["Name"].S,
		})
	}

	response, err := json.Marshal(movies)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error while decoding to string value",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "http://localhost:4200",
		},
		Body: string(response),
	}, nil
}

func main() {
	lambda.Start(findAll)
}
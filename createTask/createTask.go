package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Task struct {
	UserName    string `json:"User"`
	DateCreated string `json:"DateCreated"`
	TaskName    string `json:"taskName"`
	Description string `json:"description"`
	TaskRunTime string `json:"taskRunTime"`
}

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Creating Task...")
	var task Task
	TaskBytes := []byte(string(request.Body))
	err := json.Unmarshal(TaskBytes, &task)
	if err != nil {
		fmt.Println(err)
	}

	task.DateCreated = time.Now().String()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	av, err := dynamodbattribute.MarshalMap(task)
	if err != nil {
		fmt.Println(err)
	}

	tableName := "Tasks"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error in Put item")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       string("Put Item Error"),
			StatusCode: 502,
		}, nil
	}

	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}
	TaskJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error")
	}
	return events.APIGatewayProxyResponse{
		Body:       string(TaskJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

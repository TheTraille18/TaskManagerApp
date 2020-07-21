package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//Task object
type Task struct {
	User        string `json:"User"`
	DateCreated string `json:"DateCreated"`
	TaskName    string `json:"taskName"`
	Description string `json:"description"`
	TaskRunTime string `json:"taskRunTime"`
}

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
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

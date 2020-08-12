package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
)

type Task struct {
	User     string `json:"User"`
	TaskName string `json:"TaskName"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Deleting Task")
	//Get Request Body
	fmt.Println(request.Body)
	TaskBytes := []byte(request.Body)

	//Unmarshal to Task struct
	var task Task
	err := json.Unmarshal(TaskBytes, &task)
	if err != nil {
		fmt.Println("Error Marshalling")
		fmt.Println(err)
	}

	//Delete Input
	tableName := "TaskManagerApp-Tasks"
	fmt.Println(task)
	deleteInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"User": {
				S: aws.String(task.User),
			},
			"TaskName": {
				S: aws.String(task.TaskName),
			},
		},
		TableName: aws.String(tableName),
	}

	//Delete From Dynamodb
	_, err = svc.DeleteItem(deleteInput)
	if err != nil {
		fmt.Println("Delete item error")
		fmt.Println(err)
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error Marshaling")
		fmt.Println(err)
	}

	//Headers
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(jsonTask),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

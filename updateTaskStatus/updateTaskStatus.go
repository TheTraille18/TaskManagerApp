package main

import (
	"fmt"

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

type Event struct {
	Payload  string `json:"payload"`
	User     string `json:'User'`
	TaskName string `json:'TaskName'`
}

type Task struct {
	User     string `json:'User'`
	TaskName string `json:'TaskName'`
}

func HandlerRequest(e Event) {

	//TAsk to be Update to Inactive
	task := Task{
		User:     e.User,
		TaskName: e.TaskName,
	}

	tableName := "TaskManagerApp-Tasks"
	newStatus := "Inactive"

	//Input for UpdateItemFunction
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":taskStatus": {
				S: aws.String(newStatus),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"User": {
				S: aws.String(task.User),
			},
			"TaskName": {
				S: aws.String(task.TaskName),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set TaskStatus = :taskStatus"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated Status for " + task.TaskName + " for user " + task.User + " to Inactive")

	fmt.Println(task.User)
	fmt.Println(task.TaskName)
	//return Event{fmt.Sprintf("%s is handled by 1st function", e)}, nil
}

func main() {
	lambda.Start(HandlerRequest)
}

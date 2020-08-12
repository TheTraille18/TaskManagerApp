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
	"github.com/aws/aws-sdk-go/service/sfn"
)

//Task object
type Task struct {
	User        string `json:"User"`
	TaskName    string `json:"TaskName"`
	DateCreated string `json:"DateCreated"`
	Description string `json:"description"`
	TaskRunTime string `json:"taskRunTime"`
	Status      string `json:"TaskStatus"`
}

type RunningTask struct {
	User        string `json:"User"`
	TaskName    string `json:"TaskName"`
	TaskRunTime string `json:"TaskRunTime"`
}

var (
	// Create DynamoDB client
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = dynamodb.New(sess)

	//Create Step Function client
	sessStep = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	stepSvc = sfn.New(sessStep)
)

//RunTask  State Step Function
func RunTask(t Task) {
	runTask := RunningTask{
		User:        t.User,
		TaskRunTime: t.TaskRunTime,
		TaskName:    t.TaskName,
	}
	TaskJson, err := json.Marshal(runTask)
	if err != nil {
		fmt.Println("Error Marshaling")
		fmt.Println(err)
	}
	sfnInput := &sfn.StartExecutionInput{
		Input:           aws.String(string(TaskJson)),
		StateMachineArn: aws.String("arn:aws:states:us-east-1:398080922284:stateMachine:TaskManager-StateMachine"),
	}
	fmt.Println("Start Execution")
	_, err = stepSvc.StartExecution(sfnInput)
	if err != nil {
		fmt.Println("Error staring execution")
		fmt.Println(err)
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Task...")
	var task Task
	fmt.Println("Request: ")
	fmt.Println(request.Body)

	TaskBytes := []byte(string(request.Body))
	err := json.Unmarshal(TaskBytes, &task)
	if err != nil {
		fmt.Println("Error in Unmarshal")
		fmt.Println(err)
	}
	RunTask(task)
	task.DateCreated = time.Now().String()
	task.Status = "Active"

	av, err := dynamodbattribute.MarshalMap(task)
	if err != nil {
		fmt.Println(err)
	}

	tableName := "TaskManagerApp-Tasks"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Headers"] = "*"
	headers["Access-Control-Allow-Credentials"] = "true"

	fmt.Println("Writting Item to Table")
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error in Put item")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       string("Put Item Error"),
			StatusCode: 502,
		}, nil
	}

	TaskJSON, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println("Returning Response")
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(TaskJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

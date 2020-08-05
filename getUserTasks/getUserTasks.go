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

//Task Item
type Task struct {
	User        string `json:'User'`
	Description string `json:'description'`
	DateCreated string `json:'DateCreated'`
	TaskName    string `json:'taskName'`
	TaskRunTime string `json:'taskRunTime'`
}

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Getting All Task...")
	user := request.QueryStringParameters["user"]

	tableName := "TaskManagerApp-Tasks"
	fmt.Println("User: " + user)
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"User": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(user),
					},
				},
			},
		},
	}

	taskQuery, err := svc.Query(queryInput)
	if err != nil {
		fmt.Println("Got error in getting User Tasks")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       string("Put Error"),
			StatusCode: 502,
		}, nil
	}

	//Create Task Lists
	Tasks := make([]Task, 0)
	for _, taskItem := range taskQuery.Items {
		user := taskItem["User"]
		dateCreated := taskItem["DateCreated"]
		description := taskItem["description"]
		taskName := taskItem["taskName"]
		taskRunTime := taskItem["taskRunTime"]

		task := Task{
			User:        *user.S,
			DateCreated: *dateCreated.S,
			Description: *description.S,
			TaskName:    *taskName.S,
			TaskRunTime: *taskRunTime.S,
		}
		Tasks = append(Tasks, task)
	}

	//Convert list to Json for response output
	jsonTasks, err := json.Marshal(Tasks)
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
		Body:       string(jsonTasks),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

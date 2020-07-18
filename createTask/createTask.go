package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Task struct {
	TaskID      string
	TaskName    string `json:"name"`
	Description string `json:"description"`
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
	var task Task
	fmt.Println("jfdsajfaoijiosjoiajfoaijjiodsjfiojeoij")
	TaskBytes := []byte(string(request.Body))
	err := json.Unmarshal(TaskBytes, &task)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! " + task.TaskName)
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Task Created!, %v", string(task.TaskName)),
		StatusCode: 200,
	}, nil
}

func main() {
	fmt.Println("Creating Task.....")
	lambda.Start(handler)
}

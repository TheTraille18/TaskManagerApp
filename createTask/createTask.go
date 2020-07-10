package main

import (
	"fmt"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

//Initial Struct
type MyEvent struct {
	Name string `json:"name"`
}

//Handles the Request
func HandleRequest(ctx context.Context, name MyEvent) (string, error){
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main(){
	lambda.Start(HandleRequest)
}
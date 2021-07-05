package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	pkg "testlocalstack/sqs_lambda/handler"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age int     `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	fmt.Println("start recieving message result")
	qName :="http://localhost:4566/000000000000/myQueue"
	//qArn:="arn:aws:sqs:us-east-1:000000000000:"+qName
	result,err:= pkg.RecieveMessage( qName,15)
	if err != nil {
		fmt.Println(err.Error())
		return MyResponse{Message: err.Error()}, nil
	}
	fmt.Println("Recieve message result")
	fmt.Println(result)



	return MyResponse{Message: result}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}


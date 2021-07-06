package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age  int    `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Println("000000000000000000000000")
	for _, r := range sqsEvent.Records {
		fmt.Println("reccccccccc")
		fmt.Println(r.MessageId)
		fmt.Println(r.Body)
		fmt.Println(r.Attributes)
		fmt.Println(r.MessageAttributes)
	}
	fmt.Println("start recieving message result")
	//qName := "http://localhost:4566/000000000000/myQueue"
	//qArn:="arn:aws:sqs:us-east-1:000000000000:"+qName
	/*result, err := pkg.RecieveMessage(qName, 15)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println("Recieve message result")
	fmt.Println(result)*/

	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

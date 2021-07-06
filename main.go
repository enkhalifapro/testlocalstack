package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"testlocalstack/pkg"
	"time"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age  int    `json:"How old are you?"`
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	/*topic,err:= pkg.CreateSNSTopic("topicSLS")
	if err != nil {
		fmt.Println(err.Error())
		return MyResponse{Message: err.Error()}, nil
	}
	fmt.Println(topic)

	/*qName := "myQueue"

	qUrl,err:=pkg.CreateSQSQueue(qName)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(qUrl)*/
	qArn := "arn:aws:sqs:us-east-1:000000000000:myQueue"
	topic := "arn:aws:sns:us-east-1:000000000000:topicSLS"
	sArn, err := pkg.SubscribeToSNSTopic(qArn, topic)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(sArn)

	for i := 0; i < 10; i++ {
		messageId, err := pkg.PublishToSNSTopic(topic, fmt.Sprintf("message in time cccc %v", time.Now().Unix()))
		if err != nil {
			fmt.Println("errrrrr")
			fmt.Println(err)
		}
		fmt.Println("sennnnnnnndddddd")
		fmt.Println(messageId)
		time.Sleep(1 * time.Second)
	}

	return MyResponse{Message: topic}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}

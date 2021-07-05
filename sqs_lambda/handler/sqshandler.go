package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
	"os"
)

var (
	awsRegion   string
	awsEndpoint string

	sqsClient *sqs.Client
)


type SQSReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

// GetMessages gets the most recent message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func init() {
	awsRegion ="us-east-1"// os.Getenv("AWS_REGION")
	awsEndpoint =os.Getenv("LOCALSTACK_HOSTNAME")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				URL:           fmt.Sprintf("http://%s:4566", awsEndpoint),
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

fmt.Println("==========================")
	fmt.Println("initialize sqs .....")


	sqsClient = sqs.NewFromConfig(awsCfg)
	fmt.Println("==========================")
}


func RecieveMessage(queue string,timeout int) (string,error) {
	if queue == "" {
		return "",errors.New("You must supply the name of a queue")
	}
	if timeout < 0 {
		timeout = 0
	}

	if timeout > 12*60*60 {
		timeout = 12 * 60 * 60
	}

	// Get URL of queue
	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            &queue,
		MaxNumberOfMessages: 5,
		VisibilityTimeout:   int32(timeout),
	}
	fmt.Println("\nstart reading message from sqs\n")
	msgResult, err := GetMessages(context.TODO(), sqsClient, gMInput)

	if err != nil {
		fmt.Println("Got an error receiving messages:")
		fmt.Println(err)
		return "",err
	}
	fmt.Println("\n *** result *** \n")
	fmt.Println(msgResult)
	fmt.Println("\n==== Messages count====\n")

	fmt.Println(len(msgResult.Messages))
	fmt.Println("\n==== End Messages====\n")
	return *msgResult.Messages[0].MessageId,nil
	//fmt.Println("Message ID:     " + *msgResult.Messages[0].MessageId)
	//fmt.Println("Message Handle: " + *msgResult.Messages[0].ReceiptHandle)
}


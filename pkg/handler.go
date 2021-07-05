package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"os"
	"strings"
)

// CustomEvent for lambda
type CustomEvent struct {
	ID   string
	Name string
}

var (
	awsRegion   string
	awsEndpoint string
	bucketName  string

	s3svc *s3.Client
	snsClient *sns.Client
	sqsClient *sqs.Client
)
type SNSCreateTopicAPI interface {
	CreateTopic(ctx context.Context,
		params *sns.CreateTopicInput,
		optFns ...func(*sns.Options)) (*sns.CreateTopicOutput, error)
}
type SNSPublishAPI interface {
	Publish(ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
func MakeTopic(c context.Context, api SNSCreateTopicAPI, input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error) {
	return api.CreateTopic(c, input)
}
func PublishMessage(c context.Context, api SNSPublishAPI, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return api.Publish(c, input)
}
type SNSSubscribeAPI interface {
	Subscribe(ctx context.Context,
		params *sns.SubscribeInput,
		optFns ...func(*sns.Options)) (*sns.SubscribeOutput, error)
}
type SQSCreateQueueAPI interface {
	CreateQueue(ctx context.Context,
		params *sqs.CreateQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
}

// CreateQueue creates an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateQueueOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateQueue.
func CreateQueue(c context.Context, api SQSCreateQueueAPI, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return api.CreateQueue(c, input)
}
func SubscribeTopic(c context.Context, api SNSSubscribeAPI, input *sns.SubscribeInput) (*sns.SubscribeOutput, error) {
	return api.Subscribe(c, input)
}
func init() {
	awsRegion ="us-east-1"// os.Getenv("AWS_REGION")
	awsEndpoint =os.Getenv("LOCALSTACK_HOSTNAME")// "http://localhost:4566/"//os.Getenv("AWS_ENDPOINT")
	bucketName = "ayman"//os.Getenv("S3_BUCKET")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:            fmt.Sprintf("http://%s:4566", awsEndpoint),
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

	s3svc = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	snsClient = sns.NewFromConfig(awsCfg)
	sqsClient = sqs.NewFromConfig(awsCfg)
}
//Add text file to local aws bucket
func addTextToBucket(fileName string,textBody string) error {
	s3Key := fmt.Sprintf("%s.txt",fileName)
	body := []byte(fmt.Sprintf(textBody))
	resp, err := s3svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(s3Key),
		Body:               bytes.NewReader(body),
		ContentLength:      int64(len(body)),
		ContentType:        aws.String("application/text"),
		ContentDisposition: aws.String("attachment"),
	})
	log.Printf("S3 PutObject response: %+v", resp)

	return  err

}

func CreateSNSTopic(topic string) (string,error)  {

	input := &sns.CreateTopicInput{
		Name: &topic,
	}
	//results, err :=snsClient.CreateTopic(context.TODO(),  input)
	results, err := MakeTopic(context.TODO(), snsClient, input)
	if err != nil {
		return "",err
	}
	return  *results.TopicArn,nil

}

func PublishToSNSTopic(topicArn string,msg string) (string,error)  {

	input := &sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicArn,
	}

	result, err := PublishMessage(context.TODO(), snsClient, input)
	if err != nil {
		return "",err
	}
	return  *result.MessageId,nil

}

func CreateSQSQueue(queue string) (string,error)  {

	listQueuesRequest := sqs.ListQueuesInput{}

	listQueueResults, _ := sqsClient.ListQueues( context.TODO(),&listQueuesRequest)

	for _, t := range listQueueResults.QueueUrls {
		// If one of the returned queue URL's contains the required name we need then break the loop
		if strings.Contains(t, queue) {
			return t,nil
		}
	}
	input := &sqs.CreateQueueInput{
		QueueName: &queue,
		Attributes: map[string]string{
			"DelaySeconds":           "60",
			"MessageRetentionPeriod": "86400",
		},
	}

	result, err := CreateQueue(context.TODO(), sqsClient, input)
	if err != nil {
		fmt.Println("Got an error creating the queue:")
		fmt.Println(err)
		return "",err
	}
	return  *result.QueueUrl,nil

}
func SubscribeToSNSTopic(qUrl string,topicARN string) (string,error) {
	input := &sns.SubscribeInput{
		Endpoint:              &qUrl,
		Protocol:              aws.String("sqs"),
		ReturnSubscriptionArn: true, // Return the ARN, even if user has yet to confirm
		TopicArn:              &topicARN,
	}

	result, err := SubscribeTopic(context.TODO(), snsClient, input)
	if err != nil {
		return "",err
	}
	return  *result.SubscriptionArn,nil

}
/*
func main() {

	topic,err:=createSNSTopic("myTopic")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//topic:="arn:aws:sns:us-east-1:000000000000:myTopic"
	qName := "myQueue"
	qArn:="arn:aws:sqs:us-east-1:000000000000:"+qName
	qUrl,err:=createQueue(qName)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(qUrl)
	sArn,err:=subscribeToSNSTopic(qArn,topic)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(sArn)
	messageId,err:=publishToSNSTopic(topic,"Hello to all")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(messageId)




}
*/



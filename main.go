package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
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
func init() {
	awsRegion ="us-east-1"// os.Getenv("AWS_REGION")
	awsEndpoint = "http://localhost:4566/"//os.Getenv("AWS_ENDPOINT")
	bucketName = "ayman"//os.Getenv("S3_BUCKET")

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
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

func createSNSTopic(topic string) (string,error)  {

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

func publishToSNSTopic(topicArn string,msg string) (string,error)  {

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
func main() {


	/*err:=addTextToBucket("ayman2","hello ayman")
	if err != nil {
		fmt.Println(err.Error())
	}*/


	/*arn,err:=createSNSTopic("myTopic")
	if err != nil {
		fmt.Println(err.Error())
	}*/
	//print(arn)

	messageId,err:=publishToSNSTopic("arn:aws:sns:us-east-1:000000000000:myTopic","Hello to all")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(messageId)
}


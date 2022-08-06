# SQS ROUTER


## Install
```
go get github.com/sofyan48/mq-router
```
## Example
``` golang
package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/sofyan48/mq-router"
)

func main() {
	client := newClient()
	topic1Url := "https://sqs.ap-southeast-1.amazonaws.com/AWS_ID/sqs-example-topic"

	router := mq.NewRouter()

	router.Handle("/push/wa", mq.HandlerFunc(func(m *mq.Message) error {
		fmt.Println("WA:> ", aws.StringValue(m.SQSMessage.Body))
		return nil
	})).Method("POST")
	router.Handle("/push/pu", mq.HandlerFunc(func(m *mq.Message) error {
		fmt.Println("PU:> ", aws.StringValue(m.SQSMessage.Body))
		return nil
	})).Method("POST")

	topic1 := mq.NewServer(topic1Url, router, mq.WithClient(client))

	for {
		topic1.Start()
		time.Sleep(1 * time.Second)
	}
}

func newClient() sqsiface.SQSAPI {

	creds := credentials.NewStaticCredentials(
		"ACCESS_KEY",
		"SECRET_KEY", "")

	credential := &aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-southeast-1"),
	}
	return sqs.New(session.New(), credential)
}
```
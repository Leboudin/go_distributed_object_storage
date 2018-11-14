package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
)

func NoErr(err error) {
	if err != nil {
		log.Println(err)
		log.Fatal("Exiting...")
		return
	}
}

func newServiceClient() *sqs.SQS {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	NoErr(err)

	// check session
	cre, err := ssn.Config.Credentials.Get()
	log.Printf("Credentials: %v", cre)
	svc := sqs.New(ssn)
	return svc
}

func listQueues(svc *sqs.SQS) {
	result, err := svc.ListQueues(nil)
	NoErr(err)
	log.Printf("Listed queues: %v", result.QueueUrls)
}

func createQueue(svc *sqs.SQS) {
	queueInput := sqs.CreateQueueInput{
		QueueName: aws.String("test-queue"),
	}
	result, err := svc.CreateQueue(&queueInput)
	NoErr(err)
	log.Printf("Created queue: %v", result.QueueUrl)
}

func getQueueUrl(svc *sqs.SQS) *string {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("godos-test"),
	})
	NoErr(err)
	log.Printf("Got queue url for queue godos-test: %v", *result.QueueUrl)
	return result.QueueUrl
}

func sendMessage(svc *sqs.SQS) {
	url := getQueueUrl(svc)
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(*url),
		MessageBody: aws.String("This is a test message"),
	})

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Successfully sent a message, ID: %v, NO.: %v",
			result.MessageId, result.SequenceNumber)
	}
}

func consumeMessage(svc *sqs.SQS) {
	url := getQueueUrl(svc)
	for {
		time.Sleep(1)
		log.Println("Consumer event loop end")
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(*url),
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(0),
		})

		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		if len(result.Messages) == 0 {
			log.Println("Received no message, continue")
			continue
		}

		// we got a message, let's consume it
		msg := result.Messages[0]
		log.Printf("Consumed message, ID: %v, body: %v",
			*msg.MessageId, *msg.Body)
		//log.Println(msg.String())

		// then delete the message
		_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      url,
			ReceiptHandle: msg.ReceiptHandle,
		})

		if err != nil {
			log.Printf("Failed to delete message %v", *msg.MessageId)
			continue
		}

		log.Printf("Message %v deleted", *msg.MessageId)
	}
}

func main() {
	log.Println("Start testing...")
	svc := newServiceClient()
	log.Printf("New service created: %v", svc.ServiceID)

	go func(svc *sqs.SQS) {
		consumeMessage(svc)
	}(svc)

	//listQueues(svc)
	//createQueue(svc)
	//getQueueUrl(svc)
	//sendMessage(svc)
	//runtime.Gosched()
	for {
		time.Sleep(1)
	}
}

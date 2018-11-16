package sqs

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type SQS struct {
	// Queue Url
	Url string

	// reply Url
	ReplyUrl string

	// AWS SQS service client
	svc *sqs.SQS

	// Message channel
	msgC chan sqs.Message

	// Close channel, used to close connection
	closeC chan bool
}

func newServiceClient() (*sqs.SQS, error) {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})
	if err != nil {
		return nil, err
	}

	// check session
	svc := sqs.New(ssn)
	return svc, nil
}

func getQueueUrl(svc *sqs.SQS, name string) (*string, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	//log.Printf("Got queue Url for queue godos-test: %v", *result.QueueUrl)
	return result.QueueUrl, nil
}

func NewSQS() *SQS {
	// Create session
	svc, err := newServiceClient()
	if err != nil {
		log.Printf("Failed to create session, error: %s", err)
		log.Fatal("Now exiting...")
	}

	// Get object location query queue Url
	url, err := getQueueUrl(svc, "godos-test")
	if err != nil {
		log.Printf("Failed to get queue Url, error: %s", err)
		log.Fatal("Now exiting...")
	}

	// Get object located reply queue Url
	ReplyUrl, err := getQueueUrl(svc, "godos-test-located")
	if err != nil {
		log.Printf("Failed to get reply queue Url, error: %s", err)
		log.Fatal("Now exiting...")
	}

	return &SQS{
		svc:      svc,
		Url:      *url,
		ReplyUrl: *ReplyUrl,
		msgC:     make(chan sqs.Message),
		closeC:   make(chan bool),
	}
}

func NewSQSFromUrl(url string, replyUrl string) *SQS {
	// Create session
	svc, err := newServiceClient()
	if err != nil {
		log.Printf("Failed to create session, error: %s", err)
		log.Fatal("Now exiting...")
	}

	return &SQS{
		svc:      svc,
		Url:      url,
		ReplyUrl: replyUrl,
		msgC:     make(chan sqs.Message),
		closeC:   make(chan bool),
	}
}

func (s *SQS) consume(url string) {
	result, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(url),
		MaxNumberOfMessages: aws.Int64(1), // receive 1 message per time, for test purpose?
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		log.Printf("An error occurred while receiving message, error: %s", err)
		return
	}

	if len(result.Messages) == 0 {
		return
	}

	// We successfully got a message, let's put it into the msgC channel
	msg := result.Messages[0]
	s.msgC <- *msg
}

func (s *SQS) Consume(url string) <-chan sqs.Message {
	go func(url string) {
		for {
			select {
			case <-s.closeC:
				close(s.msgC)
				log.Println("Closing message channel")
				return
			default:
				s.consume(url)
			}
		}
	}(url)

	return s.msgC
}

// Send message, used by data provider server to tell API server that it holds current object
func (s *SQS) SendMessage(msg map[string]string, url string) (sqs.SendMessageOutput, error) {
	msgStr, err := json.Marshal(msg)

	if err != nil {
		return sqs.SendMessageOutput{}, err
	}
	result, err := s.svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(url),
		MessageBody: aws.String(string(msgStr)),
	})
	return *result, err
}

// Delete specified message
func (s *SQS) DeleteMessage(msg sqs.Message, url string) error {
	_, err := s.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: msg.ReceiptHandle,
	})

	return err
}

func (s *SQS) Close() {
	s.closeC <- true
}

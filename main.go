package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	QueueUrl    = "https://sqs.ap-southeast-2.amazonaws.com/456511713453/SQS-Jeff-test-1.fifo"
	Region      = "ap-southeast-2"
	CredPath    = "/home/jeff/.aws/credentials"
	CredProfile = "aws-cred-profile"
)

func main() {

	sess := session.New(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewSharedCredentials(CredPath, CredProfile),
		MaxRetries:  aws.Int(5),
	})

	svc := sqs.New(sess)
	begin := time.Now()

	content, err := ioutil.ReadFile("./tt")
	if err != nil {
		log.Println("open failed")
		return
	}
	// Send message
	for i := 0; i < 10; i++ {
		send_params := &sqs.SendMessageInput{
			MessageBody:            aws.String("message body new " + strconv.Itoa(i)), // Required
			QueueUrl:               aws.String(QueueUrl),                              // Required
			DelaySeconds:           aws.Int64(0),                                      // (optional)
			MessageGroupId:         aws.String("group1"),
			MessageDeduplicationId: aws.String(time.Now().Format("2006-01-02_15:04:05.999999-07:00") + strconv.Itoa(i)),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"payload": &sqs.MessageAttributeValue{
					DataType:    aws.String("Binary"),
					BinaryValue: content,
					//BinaryValue: []byte("SGVsbG8gQmluYXJ5"),
				},
				"Author": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("test"),
				},
			},
		}
		send_resp, err := svc.SendMessage(send_params)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[Send message] \n%v \n\n", send_resp)
	}
	fmt.Println("duration for send", time.Since(begin))
	/*
		// Receive message
		receive_params := &sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            aws.String(QueueUrl),
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(30),
			WaitTimeSeconds:     aws.Int64(20), // long polling
		}
		receive_resp, err := svc.ReceiveMessage(receive_params)
		if err != nil {
			log.Println(err)
			return
		}
		for _, msg := range receive_resp.Messages {
			fmt.Printf("[Receive message] \n%v \n\n", msg)
			for key, attr := range msg.MessageAttributes {
				fmt.Println("attr: ", key)
				switch *attr.DataType {
				case "Binary":
					fmt.Println("   with binary value:", string(attr.BinaryValue[0:10]))
				case "String":
					fmt.Println("   with string value:", *attr.StringValue)
				default:
					fmt.Println("   with unknowb data type ", attr.DataType)
				}
			}
		}

		fmt.Println("duration for send + receive", time.Since(begin))

		for _, message := range receive_resp.Messages {
			delete_params := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(QueueUrl),  // Required
				ReceiptHandle: message.ReceiptHandle, // Required
			}
			_, err := svc.DeleteMessage(delete_params) // No response returned when successed.
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
		}
		fmt.Println("duration for send+rec+del", time.Since(begin))

		return
	*/
}

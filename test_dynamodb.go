package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type DB struct {
	// session
	Svc *dynamodb.DynamoDB
}

func NewDB() *DB {
	ssn, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	})

	if err != nil {
		log.Printf("Failed to create session, error: %s", err)
	}

	svc := dynamodb.New(ssn)
	return &DB{
		Svc: svc,
	}
}

func (db *DB) TestListTables() {
	result, err := db.Svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("Failed to list tables, error: %s", err)
	} else {
		log.Println("Tables: ")
		for _, n := range result.TableNames {
			log.Println(*n)
		}
	}

}

func (db *DB) TestCreateTable() {
	result, err := db.Svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String("users"),

		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("uid"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			//{
			//	AttributeName: aws.String("addr"),
			//	AttributeType: aws.String("S"),
			//},
		},

		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("uid"),
				KeyType:       aws.String("HASH"),	// uid is used to shard
			},
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("RANGE"),	// name is used as sort key
			},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	if err != nil {
		log.Printf("Failed to create table, error: %s", err)
	} else {
		log.Printf("Created table users, table status: ")
		log.Println(*result.TableDescription.TableName)
		log.Println(*result.TableDescription.TableSizeBytes)
		log.Println(*result.TableDescription.TableStatus)
		log.Println(*result.TableDescription.ItemCount)
	}
}

type UserDetail struct {
	Addr  string `json:"addr"`
	Phone string `json:"phone"`
}

type User struct {
	Uid    int64      `json:"uid"`
	Name   string     `json:"name"`
	Detail UserDetail `json:"detail"`
}

func (db *DB) TestCreateItem() {
	// create sample data
	detail := UserDetail{
		Addr:  "some address NO.14",
		Phone: "1231231234",
	}

	user := User{
		Uid:    111,
		Name:   "TestUser",
		Detail: detail,
	}

	val, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("Failed to marshal data")
		return
	}

	_, err = db.Svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      val,
	})

	if err != nil {
		log.Printf("Failed to put item, error: %s", err)
		return
	}

}

func main() {
	db := NewDB()
	//db.TestCreateTable()
	db.TestListTables()
	//db.TestCreateItem()
}

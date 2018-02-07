package user

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws"
	"os"
)

type Store interface {
	saveKey(user, key string) error
	getDetails(user string) (*Multipooluser, error)
	Iterate(fn func(*dynamodb.ScanOutput, bool) bool) error
}

type dynamoStore struct {
	db   *dynamodb.DynamoDB
	More bool
}

type Multipooluser struct {
	Userid string  `json:"userid"`
	Apikey string
}

func NewDynamoStore(db *dynamodb.DynamoDB) Store {
	return &dynamoStore{db: db}
}

func (s *dynamoStore) saveKey(user, key string) error {
	r := Multipooluser{Userid: user, Apikey: key}
	av, err := dynamodbattribute.MarshalMap(r)

	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(getMultiPoolTable()),
	}

	_, err = s.db.PutItem(input)
	return err
}

func (s dynamoStore) getDetails(user string) (*Multipooluser, error) {
	get := &dynamodb.GetItemInput{
		TableName: aws.String(getMultiPoolTable()),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				S: aws.String(user),
			},
		},
	}
	result, err := s.db.GetItem(get)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	r := Multipooluser{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	return &r, err
}

func (s dynamoStore) Iterate(fn func(*dynamodb.ScanOutput, bool) bool) error {
	input := &dynamodb.ScanInput{
		TableName:aws.String(getMultiPoolTable()),
	}
		return s.db.ScanPages(input,fn)

}

func getMultiPoolTable() string {
	return os.Getenv("MULTIPOOL_TABLE")
}

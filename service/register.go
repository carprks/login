package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/satori/go.uuid"
	"os"
)

func RegisterService(body string) (string, error) {
	r := Register{}

	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("register unmarshall: %v", err))
		return "", err
	}

	// check the passwords are the same
	if r.Password != r.Verify {
		return "", fmt.Errorf("passwords don't match")
	}

	alreadyExists, err := r.CheckEmail()
	if alreadyExists {
		return "", fmt.Errorf("email already exists")
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("check email: %v", err))
		return "", err
	}

	crypt, err := HashPassword(r.Password)
	if err != nil {
		fmt.Println(fmt.Sprintf("crypt password: %v", err))
		return "", err
	}
	r.Crypt = crypt

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return "", err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(r.createIdentifier()),
			},
			"email": {
				S: aws.String(r.Email),
			},
			"phone": {
				S: aws.String(r.Phone),
			},
			"password": {
				S: aws.String(r.Crypt),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(#IDENTIFIER)"),
		ExpressionAttributeNames: map[string]*string{
			"#IDENTIFIER": aws.String("identifier"),
		},
	}
	_, err = svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return "", fmt.Errorf("ErrCodeConditionalCheckFailedException: %v", aerr)
			case "ValidationException":
				fmt.Println(fmt.Sprintf("validation err reason: %v", input))
				return "", fmt.Errorf("validation error: %v", aerr)
			default:
				fmt.Println(fmt.Sprintf("unknown code err reason: %v", input))
				return "", fmt.Errorf("unknown code err: %v", aerr)
			}
		}
	}


	return "registered", nil
}

func (r Register)createIdentifier() string {
	u := uuid.NewV5(uuid.NamespaceURL, fmt.Sprintf("https://identity.carprk.com/user/%s:%s", r.Email, r.Phone))
	return u.String()
}

func (r Register)CheckEmail() (bool, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return false, err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DB_TABLE")),
		FilterExpression: aws.String("Email = :email"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(r.Email),
			},
		},
	}
	result, err := svc.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(fmt.Sprintf("ErrCodeProvisionedThroughputExceededException: $%v", aerr))
				return false, aerr
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(fmt.Sprintf("ErrCodeResourceNotFoundException: %v", aerr))
				return false, aerr
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(fmt.Sprintf("ErrCodeRequestLimitExceeded: %v", aerr))
				return false, aerr
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(fmt.Sprintf("ErrCodeInternalServerError: %v", aerr))
				return false, aerr
			default:
				fmt.Println(fmt.Sprintf("unknown: %v", aerr))
				return false, aerr
			}
		}
		return false, err
	}

	if len(result.Items) >= 1 {
		return true, nil
	}

	return false, nil
}
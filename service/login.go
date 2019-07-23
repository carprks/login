package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// LoginService ...
func LoginService(body string) (string, error) {
	r := Login{}

	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("login unmarshall: %v", err))
		return "", err
	}

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(s)
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#EMAIL": aws.String("email"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(r.Email),
			},
		},
		FilterExpression: aws.String("#EMAIL = :email"),
		TableName: aws.String(os.Getenv("DB_TABLE")),
	}
	result, err := svc.Scan(input)
	if err != nil {
		return "", err
	}
	if len(result.Items) != 1 {
		return "", fmt.Errorf("no user, or more than 1")
	}

	crypt := ""
	for key, value := range result.Items[0] {
		switch key {
		case "password":
			crypt = *value.S
		}
	}

	valid := CheckPassword(crypt, r.Password)
	if !valid {
		return "", fmt.Errorf("invalid password")
	}

	return "success", nil
}

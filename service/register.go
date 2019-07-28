package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"os"
)

func register(body string) (string, error) {
	r := RegisterRequest{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		return "", err
	}

  if os.Getenv("DEVELOPMENT") == "" {
    if r.Email == "tester@carpark.ninja" {
      return "", fmt.Errorf("tester account not allowed in production")
    }
  }

	resp, err := r.Register()
	if err != nil {
		fmt.Println(fmt.Sprintf("register service err: %v", err))
		return "", err
	}

	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(fmt.Sprintf("register marshall: %v", err))
		return "", err
	}

	return string(res), nil
}

// Register ...
func (r RegisterRequest) Register() (Register, error) {
	if len(os.Args) >= 1 {
		if os.Args[2] == "localDev" {
			err := godotenv.Load()
			if err != nil {
				fmt.Println(fmt.Sprintf("godotenv err: %v", err))
			}
		}
	}

	// check the passwords are the same
	if r.Password != r.Verify {
		fmt.Println("Passwords dont match")
		return Register{}, fmt.Errorf("passwords don't match")
	}

	alreadyExists, err := r.CheckEmail()
	if alreadyExists {
		fmt.Println("account already exists")
		return Register{}, fmt.Errorf("email already exists")
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("check email: %v", err))
		return Register{}, err
	}

	crypt, err := HashPassword(r.Password)
	if err != nil {
		fmt.Println(fmt.Sprintf("crypt password: %v", err))
		return Register{}, err
	}
	r.Crypt = crypt

	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return Register{}, err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(GenerateIdent(r.Email)),
			},
			"email": {
				S: aws.String(r.Email),
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
				return Register{}, fmt.Errorf("ErrCodeConditionalCheckFailedException: %v", aerr)
			case "ValidationException":
				fmt.Println(fmt.Sprintf("validation err reason: %v", input))
				return Register{}, fmt.Errorf("validation error: %v", aerr)
			default:
				fmt.Println(fmt.Sprintf("unknown code err reason: %v", input))
				return Register{}, fmt.Errorf("unknown code err: %v", aerr)
			}
		}
	}

	return Register{
		ID:    GenerateIdent(r.Email),
		Email: r.Email,
	}, nil
}

// CheckEmail ...
func (r RegisterRequest) CheckEmail() (bool, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return false, err
	}
	svc := dynamodb.New(s)
	input := &dynamodb.ScanInput{
		TableName:        aws.String(os.Getenv("DB_TABLE")),
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
		fmt.Println(fmt.Sprintf("really unknown: %v", err))
		return false, err
	}

	if len(result.Items) >= 1 {
		fmt.Println("something found")
		return true, nil
	}

	return false, nil
}

package service

import (
  "encoding/json"
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/badoux/checkmail"
  "os"
)

func register(body string) (string, error) {
	r := RegisterRequest{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
    res, err := json.Marshal(Register{
      Error: err.Error(),
    })
    if err != nil {
      return "", err
    }
    return string(res), nil
	}

  if os.Getenv("DEVELOPMENT") == "" {
    if r.Email == "tester@carpark.ninja" {
      res, err := json.Marshal(Register{
        Error:   fmt.Errorf("tester account not allowed in production").Error(),
      })
      if err != nil {
        return "", err
      }
      return string(res), nil
    }
  }

	resp, err := r.Register()
  if err != nil {
    res, err := json.Marshal(Register{
      Error: err.Error(),
    })
    if err != nil {
      return "", err
    }
    return string(res), nil
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
	// check the passwords are the same
	if r.Password != r.Verify {
		return Register{}, fmt.Errorf("passwords don't match")
	}

	alreadyExists, err := r.EmailExists()
	if alreadyExists {
		return Register{}, fmt.Errorf("email already exists")
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("check email: %v", err))
		return Register{}, err
	}

	emailErr := r.emailTest()
	if emailErr != nil {
	  return Register{}, emailErr
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
func (r RegisterRequest) EmailExists() (bool, error) {
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
				return false, aerr
			case dynamodb.ErrCodeResourceNotFoundException:
				return false, aerr
			case dynamodb.ErrCodeRequestLimitExceeded:
				return false, aerr
			case dynamodb.ErrCodeInternalServerError:
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
		return true, nil
	}

	return false, nil
}

func (r RegisterRequest)emailTest() error {
  err := checkmail.ValidateFormat(r.Email)
  if err != nil {
    return err
  }

  if os.Getenv("DEVELOPMENT") != "" {
    if r.Email == "tester@carpark.ninja" || r.Email == "testfail-login@carpark.ninja"{
      return nil
    }
  }
  err = checkmail.ValidateHost(r.Email)
  if serr, ok := err.(checkmail.SmtpError); ok && err != nil {
    if serr.Code() == "550" {
      return fmt.Errorf("invalid email")
    }
    return serr
  }

  return nil
}

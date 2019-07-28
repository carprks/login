package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func login(body string) (string, error) {
	r := LoginRequest{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
	  fmt.Println(fmt.Sprintf("login request err: %v", err))
    res, err := json.Marshal(Login{
      Success: false,
      Error:   err,
    })
    if err != nil {
      return "", err
    }
    return string(res), nil
	}

	if os.Getenv("DEVELOPMENT") == "" {
	  if r.Email == "tester@carpark.ninja" {
	    fmt.Println(fmt.Sprintf("tester account not allowed in production"))
	    res, err := json.Marshal(Login{
        Success: false,
        Error:   fmt.Errorf("tester account not allowed in production"),
      })
	    if err != nil {
	      return "", err
      }
	    return string(res), nil
    }
  }

	resp, err := r.Login()
	if err != nil {
	  fmt.Println(fmt.Sprintf("login service err: %v", err))
    res, err := json.Marshal(Login{
      Success: false,
      Error:   err,
    })
    if err != nil {
      return "", err
    }
    return string(res), nil
	}

	res, err := json.Marshal(resp)
	if err != nil {
	  fmt.Println(fmt.Sprintf("login marshall err: %v", err))
		return "", err
	}
	return string(res), err
}

// Login ...
func (r LoginRequest) Login() (Login, error) {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return Login{}, err
	}
	svc := dynamodb.New(s)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(GenerateIdent(r.Email)),
			},
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	})
	if err != nil {
		return Login{}, err
	}
	if len(result.Item) == 0 {
		return Login{}, fmt.Errorf("no identity")
	}

	valid := CheckPassword(*result.Item["password"].S, r.Password)
	if !valid {
		return Login{}, fmt.Errorf("invalid password")
	}

	return Login{
		Success: true,
		ID:      *result.Item["identifier"].S,
	}, nil
}

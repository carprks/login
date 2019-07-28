package service_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func CleanTest(ident string) error {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return err
	}
	svc := dynamodb.New(s)
	_, err = svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(ident),
			},
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	})
	if err != nil {
		return err
	}

	return nil
}

func TestRegister(t *testing.T) {
	if len(os.Args) >= 1 {
		for _, env := range os.Args {
			if env == "localDev" {
				err := godotenv.Load()
				if err != nil {
					fmt.Println(fmt.Sprintf("godotenv err: %v", err))
				}
			}
		}
	}

	tests := []struct {
		request service.RegisterRequest
		expect  service.Register
		err     error
	}{
		{
			request: service.RegisterRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
				Phone:    "0123456789",
			},
			expect: service.Register{
				ID:    "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Email: "tester@carpark.ninja",
			},
		},
		{
		  request: service.RegisterRequest{
        Email:    "@carpark.ninja",
        Phone:    "123456789",
        Password: "tester",
        Verify:   "tester",
      },
      expect: service.Register{},
      err: fmt.Errorf("invalid format"),
    },
	}

	for _, test := range tests {
		response, err := test.request.Register()
		if err != nil {
		}
		passed := assert.IsType(t, test.err, err)
		if !passed {
      fmt.Println(fmt.Sprintf("register test err: %v", err))
    }
		assert.Equal(t, test.expect, response)

		CleanTest(test.expect.ID)
	}
}

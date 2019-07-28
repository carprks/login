package service_test

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type resp struct {
	ID string `json:"id"`
}

func TestHandler(t *testing.T) {
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
		request events.APIGatewayProxyRequest
		expect  events.APIGatewayProxyResponse
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{
				Resource:   "/register",
				Body:       `{"email":"tester@carpark.ninja","password":"tester","verify":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"id":"5f46cf19-5399-55e3-aa62-0e7c19382250","email":"tester@carpark.ninja"}`,
			},
		},
		{
		  request: events.APIGatewayProxyRequest{
		    Resource: "/register",
		    Body: `{"email":"@carpark.ninja","password":"tester","verify":"tester"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode:        200,
        Body:              `{"error":"invalid format"}`,
      },
    },
    {
      request: events.APIGatewayProxyRequest{
        Resource: "/register",
        Body: `{"email":"tester@carpark.ninja","password":"tester","verify":"test123"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body: `{"error":"passwords don't match"}`,
      },
    },
    {
      request: events.APIGatewayProxyRequest{
        Resource: "/register",
        Body: `{"email":"testfail-register@carpark.ninja","password":"tester","verify":"tester"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body: `{"error":"invalid email"}`,
      },
    },
    {
      request: events.APIGatewayProxyRequest{
        Resource:   "/login",
        Body:       `{"email":"tester@carpark.ninja","password":"tester"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       `{"id":"5f46cf19-5399-55e3-aa62-0e7c19382250"}`,
      },
    },
    {
      request: events.APIGatewayProxyRequest{
        Resource: "/login",
        Body: `{"email":"tester@carpark.ninja","password":"test123"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body: `{"error":"invalid password"}`,
      },
    },
    {
      request: events.APIGatewayProxyRequest{
        Resource: "/register",
        Body: `{"email":"tester-fail@carpark.ninja","password":"tester"}`,
      },
      expect: events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body: `{"error":"passwords don't match"}`,
      },
    },
	}

	for _, test := range tests {
		response, err := service.Handler(test.request)
		passed := assert.IsType(t, test.err, err)
		if !passed {
      fmt.Println(fmt.Sprintf("service test err: %v, request: %v", err, test.request))
    }
		assert.Equal(t, test.expect, response)

		s := resp{}
		err = json.Unmarshal([]byte(response.Body), &s)
		if err != nil {
			t.Fail()
		}
	}

	CleanTest("5f46cf19-5399-55e3-aa62-0e7c19382250")
}

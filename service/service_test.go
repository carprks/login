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
				Path:       "/register",
				HTTPMethod: "POST",
				Body:       `{"email":"tester@carpark.ninja","password":"tester","verify":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"id":"5f46cf19-5399-55e3-aa62-0e7c19382250","email":"tester@carpark.ninja"}`,
			},
		},
		{
			request: events.APIGatewayProxyRequest{
				Resource:   "/login",
				Path:       "/login",
				HTTPMethod: "POST",
				Body:       `{"email":"tester@carpark.ninja","password":"tester"}`,
			},
			expect: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       `{"id":"5f46cf19-5399-55e3-aa62-0e7c19382250","success":true}`,
			},
		},
	}

	for _, test := range tests {
		response, err := service.Handler(test.request)
		if err != nil {
			fmt.Println(fmt.Sprintf("service test err: %v, request: %v", err, test.request))
		}
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)

		s := resp{}
		err = json.Unmarshal([]byte(response.Body), &s)
		if err != nil {
			t.Fail()
		}
	}

	CleanTest("5f46cf19-5399-55e3-aa62-0e7c19382250")
}

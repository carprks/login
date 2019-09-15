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
	Identifier string `json:"identifier"`
}

var testsService = []struct {
	name    string
	request events.APIGatewayProxyRequest
	expect  events.APIGatewayProxyResponse
	err     error
}{
	// Register Tests
	{
		name: "register - success",
		request: events.APIGatewayProxyRequest{
			Resource: "/register",
			Body:     `{"email":"test-service@carpark.ninja","password":"tester","verify":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3","email":"test-service@carpark.ninja"}`,
		},
	},
	{
		name: "register - invalid format",
		request: events.APIGatewayProxyRequest{
			Resource: "/register",
			Body:     `{"email":"@carpark.ninja","password":"tester","verify":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"invalid format"}`,
		},
	},
	{
		name: "register - passwords don't match",
		request: events.APIGatewayProxyRequest{
			Resource: "/register",
			Body:     `{"email":"test-register-service-failure@carpark.ninja","password":"tester","verify":"test123"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"passwords don't match"}`,
		},
	},
	{
		name: "register - invalid email",
		request: events.APIGatewayProxyRequest{
			Resource: "/register",
			Body:     `{"email":"failed@carpark.ninja","password":"tester","verify":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"invalid email"}`,
		},
	},
	{
		name: "register - missing 2nd password",
		request: events.APIGatewayProxyRequest{
			Resource: "/register",
			Body:     `{"email":"failed@carpark.ninja","password":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"passwords don't match"}`,
		},
	},

	// Login Tests
	{
		name: "login - success",
		request: events.APIGatewayProxyRequest{
			Resource: "/login",
			Body:     `{"email":"test-service@carpark.ninja","password":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3"}`,
		},
	},
	{
		name: "login - invalid password",
		request: events.APIGatewayProxyRequest{
			Resource: "/login",
			Body:     `{"email":"test-service@carpark.ninja","password":"test123"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"invalid password"}`,
		},
	},
	{
		name: "login - no identity",
		request: events.APIGatewayProxyRequest{
			Resource: "/login",
			Body:     `{"email":"test-service-failure@carpark.ninja", "password":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"no identity"}`,
		},
	},
	{
		name: "login - no identity, invalid format",
		request: events.APIGatewayProxyRequest{
			Resource: "/login",
			Body:     `{"email":"@carpark.ninja","password":"tester"}`,
		},
		expect: events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       `{"error":"no identity"}`,
		},
	},
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

	for _, test := range testsService {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.Handler(test.request)
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("service type err: %w", err)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("service equal %v, %v", test.expect, response)
			}

			s := resp{}
			err = json.Unmarshal([]byte(response.Body), &s)
			if err != nil {
				t.Errorf("service unmarshal: %w", err)
			}
		})
	}

	req := events.APIGatewayProxyRequest{
		Resource: "/delete",
		Body:     `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3"}`,
	}
	resp, err := service.Handler(req)
	passed := assert.IsType(t, nil, err)
	if !passed {
		t.Errorf("service delete test err: %w, request: %v", err, req)
	}
	assert.Equal(t, events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3","status":"Deleted"}`,
	}, resp)
}

func BenchmarkHandler(b *testing.B) {
	b.ReportAllocs()

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

	b.ResetTimer()
	for _, test := range testsService {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()

			response, err := service.Handler(test.request)
			passed := assert.IsType(b, test.err, err)
			if !passed {
				b.Errorf("service test err: %w", err)
			}
			passed = assert.Equal(b, test.expect, response)
			if !passed {
				b.Errorf("service equal: %v, %v", test.expect, response)
			}

			s := resp{}
			err = json.Unmarshal([]byte(response.Body), &s)
			if err != nil {
				b.Errorf("service unmarshal: %w", err)
			}

			b.StartTimer()
		})
	}

	req := events.APIGatewayProxyRequest{
		Resource: "/delete",
		Body:     `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3"}`,
	}
	resp, err := service.Handler(req)
	passed := assert.IsType(b, nil, err)
	if !passed {
		b.Errorf("service delete test err: %w", err)
	}
	assert.Equal(b, events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"identifier":"25842e21-60a7-5bb1-8af2-5a66fb6acac3","status":"Deleted"}`,
	}, resp)
}

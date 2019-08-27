package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDelete_Delete(t *testing.T) {
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
		create  service.RegisterRequest
		request service.Delete
		expect  error
	}{
		{
			create: service.RegisterRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
			},
			request: service.Delete{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
			},
			expect: nil,
		},
	}

	for _, test := range tests {
		test.create.Register()

		response := test.request.Delete()
		passed := assert.IsType(t, test.expect, response)
		if !passed {
			fmt.Println(fmt.Sprintf("register test err: %v", response))
		}
	}
}

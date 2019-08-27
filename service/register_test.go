package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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
			},
			expect: service.Register{
				Identifier: "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Email:      "tester@carpark.ninja",
			},
		},
		{
			request: service.RegisterRequest{
				Email:    "@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
			},
			expect: service.Register{},
			err:    fmt.Errorf("invalid format"),
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

		d := service.Delete{
			Identifier: test.expect.Identifier,
		}
		d.Delete()
	}
}

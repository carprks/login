package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
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
		register service.RegisterRequest
		request  service.LoginRequest
		expect   service.Login
		err      error
	}{
		{
			register: service.RegisterRequest{
				Email:    "tester@carpark.ninja",
				Phone:    "0123456789",
				Password: "tester",
				Verify:   "tester",
			},
			request: service.LoginRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
			},
			expect: service.Login{
				ID:      "5f46cf19-5399-55e3-aa62-0e7c19382250",
				Success: true,
			},
		},
	}

	for _, test := range tests {
		reg, err := test.register.Register()
		if err != nil {
			t.Fail()
		}

		response, err := test.request.Login()
		if err != nil {
			fmt.Println(fmt.Sprintf("login test err: %v", err))
		}
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)

		CleanTest(reg.ID)
	}
}

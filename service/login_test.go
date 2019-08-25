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
				Password: "tester",
				Verify:   "tester",
			},
			request: service.LoginRequest{
				Email:    "tester@carpark.ninja",
				Password: "tester",
			},
			expect: service.Login{
				ID: "5f46cf19-5399-55e3-aa62-0e7c19382250",
			},
		},
		{
			register: service.RegisterRequest{
				Email:    "testfail-login@carpark.ninja",
				Password: "testfail",
				Verify:   "testfail",
			},
			request: service.LoginRequest{
				Email:    "testfail-login@carpark.ninja",
				Password: "tester",
			},
			expect: service.Login{},
			err:    fmt.Errorf("invalid password"),
		},
		{
			register: service.RegisterRequest{
				Email:    "testfail-login@carpark.ninja",
				Password: "tester",
				Verify:   "tester",
			},
			request: service.LoginRequest{
				Email:    "testpass@carpark.ninja",
				Password: "tester",
			},
			expect: service.Login{},
			err:    fmt.Errorf("no identity"),
		},
		{
			request: service.LoginRequest{
				Email:    "@carpark.ninja",
				Password: "tester",
			},
			expect: service.Login{},
			err:    fmt.Errorf("invalid format"),
		},
	}

	for _, test := range tests {
		reg := service.Register{}
		if test.register.Email != "" {
			reg, _ = test.register.Register()
		}

		response, err := test.request.Login()
		passed := assert.IsType(t, test.err, err)
		if !passed {
			fmt.Println(fmt.Sprintf("login test err: %v", err))
		}
		assert.Equal(t, test.expect, response)

		if reg.ID != "" {
			d := service.Delete{
				ID: reg.ID,
			}
			d.Delete()
		}
	}
}

package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testsLogin = []struct {
	name     string
	register service.RegisterRequest
	request  service.LoginRequest
	expect   service.Login
	err      error
}{
	{
		name: "success",
		register: service.RegisterRequest{
			Email:    "test-login-success@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		request: service.LoginRequest{
			Email:    "test-login-success@carpark.ninja",
			Password: "tester",
		},
		expect: service.Login{
			Identifier: "f913598f-2017-5cb7-a84e-d22008de7889",
		},
	},
	{
		name: "fail",
		register: service.RegisterRequest{
			Email:    "test-login-failure@carpark.ninja",
			Password: "testfail",
			Verify:   "testfail",
		},
		request: service.LoginRequest{
			Email:    "test-login-failure@carpark.ninja",
			Password: "test-this",
		},
		expect: service.Login{},
		err:    fmt.Errorf("invalid password"),
	},
	{
		name: "no identity",
		register: service.RegisterRequest{
			Email:    "test-login-noidentity@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		request: service.LoginRequest{
			Email:    "test-login-noidentity-failyre@carpark.ninja",
			Password: "tester",
		},
		expect: service.Login{},
		err:    fmt.Errorf("no identity"),
	},
	{
		name: "invalid format",
		request: service.LoginRequest{
			Email:    "@carpark.ninja",
			Password: "tester",
		},
		expect: service.Login{},
		err:    fmt.Errorf("invalid format"),
	},
}

func TestLoginRequest_Login(t *testing.T) {
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

	for _, test := range testsLogin {
		t.Run(test.name, func(t *testing.T) {
			reg := service.Register{}
			if test.register.Email != "" {
				reg, _ = test.register.Register()
			}

			response, err := test.request.Login()
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("login type test err: %w, %w", err, test.err)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("login equal test err: %v, %w", err, test.err)
			}

			if reg.Identifier != "" {
				d := service.Delete{
					Identifier: reg.Identifier,
				}
				d.Delete()
			}
		})
	}
}

func BenchmarkLoginRequest_Login(b *testing.B) {
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
	for _, test := range testsLogin {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()

			reg := service.Register{}
			if test.register.Email != "" {
				reg, _ = test.register.Register()
			}

			response, err := test.request.Login()
			passed := assert.IsType(b, test.err, err)
			if !passed {
				b.Errorf("login type test err: %w, %w", err, test.err)
			}
			passed = assert.Equal(b, test.expect, response)
			if !passed {
				b.Errorf("login equal test err: %v, %w", err, test.err)
			}

			if reg.Identifier != "" {
				d := service.Delete{
					Identifier: reg.Identifier,
				}
				d.Delete()
			}

			b.StartTimer()
		})
	}
}

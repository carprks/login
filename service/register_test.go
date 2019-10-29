package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testsRegister = []struct {
	name    string
	request service.RegisterRequest
	expect  service.Register
	err     error
}{
	{
		name: "success",
		request: service.RegisterRequest{
			Email:    "test-register-success@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		expect: service.Register{
			Identifier: "5358c467-3930-5d9a-a818-126c5a5c007f",
			Email:      "test-register-success@carpark.ninja",
		},
	},
	{
		name: "invalid format",
		request: service.RegisterRequest{
			Email:    "@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		expect: service.Register{},
		err:    fmt.Errorf("invalid format"),
	},
}

func TestRegisterRequest_Register(t *testing.T) {
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

	for _, test := range testsRegister {
		t.Run(test.name, func(t *testing.T) {
			response, err := test.request.Register()
			passed := assert.IsType(t, test.err, err)
			if !passed {
				t.Errorf("register type err: %w", err)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("register equal: %v, %v", test.expect, response)
			}

			d := service.Delete{
				Identifier: test.expect.Identifier,
			}
			err = d.Delete()
			if err != nil {
				t.Errorf("register delete: %w", err)
			}
		})
	}
}

func BenchmarkRegisterRequest_Register(b *testing.B) {
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
	for _, test := range testsRegister {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()

			response, err := test.request.Register()
			passed := assert.IsType(b, test.err, err)
			if !passed {
				b.Errorf("register type test err: %w", err)
			}
			passed = assert.Equal(b, test.expect, response)
			if !passed {
				b.Errorf("register equal test: %v, %v", test.expect, resp{})
			}

			d := service.Delete{
				Identifier: test.expect.Identifier,
			}
			err = d.Delete()
			if err != nil {
				b.Errorf("register delete: %w", err)
			}

			b.StartTimer()
		})
	}
}

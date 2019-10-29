package service_test

import (
	"fmt"
	"github.com/carprks/login/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testsDelete = []struct {
	name    string
	create  service.RegisterRequest
	request service.Delete
	expect  error
}{
	{
		name: "success",
		create: service.RegisterRequest{
			Email:    "test-delete-success@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		request: service.Delete{
			Identifier: "8340d46a-b9ba-5dff-9e21-b39a942a9f98",
		},
	},
	{
		name: "failed",
		create: service.RegisterRequest{
			Email:    "test-delete-failure@carpark.ninja",
			Password: "tester",
			Verify:   "tester",
		},
		request: service.Delete{
			Identifier: "df00ec88-8868-5633-90d3-dff567a12929",
		},
	},
}

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

	for _, test := range testsDelete {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.create.Register()
			if err != nil {
				t.Errorf("register create: %w", err)
			}

			response := test.request.Delete()
			passed := assert.IsType(t, test.expect, response)
			if !passed {
				t.Errorf("register type test err: %w", response)
			}
			passed = assert.Equal(t, test.expect, response)
			if !passed {
				t.Errorf("register equal test err: %v", response)
			}
		})
	}
}

func BenchmarkDelete_Delete(b *testing.B) {
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
	for _, test := range testsDelete {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()

			_, err := test.create.Register()
			if err != nil {
				b.Errorf("register create: %w", err)
			}

			response := test.request.Delete()
			passed := assert.IsType(b, test.expect, response)
			if !passed {
				b.Errorf("register type test err: %w", response)
			}
			passed = assert.Equal(b, test.expect, response)
			if !passed {
				b.Errorf("register equal test err: %v", response)
			}

			b.StartTimer()
		})
	}
}

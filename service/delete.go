package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// Delete ...
type Delete struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

func delete(body string) (string, error) {
	r := Delete{}
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		res, err := json.Marshal(Delete{
			Error: err.Error(),
		})
		if err != nil {
			return "", err
		}
		return string(res), nil
	}

	err = r.Delete()
	if err != nil {
		return "", err
	}

	return "Deleted", nil
}

// Delete ...
func (d Delete) Delete() error {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("DB_REGION")),
		Endpoint: aws.String(os.Getenv("DB_ENDPOINT")),
	})
	if err != nil {
		return err
	}
	svc := dynamodb.New(s)
	_, err = svc.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"identifier": {
				S: aws.String(d.ID),
			},
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	})
	if err != nil {
		return err
	}

	return nil
}

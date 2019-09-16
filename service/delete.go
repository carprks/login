package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

// Delete ...
type Delete struct {
	Identifier string `json:"identifier"`
	Error      string `json:"error,omitempty"`
	Status     string `json:"status,omitempty"`
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

	r.Status = "Deleted"
	res, err := json.Marshal(r)
	if err != nil {
		fmt.Println(fmt.Sprintf("delete marshall err: %v", err))
		return "", err
	}
	return string(res), err
}

// Delete ...
func (d Delete) Delete() error {
	if len(d.Identifier) <= 2 {
		return fmt.Errorf("delete invalid identifier")
	}

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
				S: aws.String(d.Identifier),
			},
		},
		TableName: aws.String(os.Getenv("DB_TABLE")),
	})
	if err != nil {
		return err
	}

	return nil
}

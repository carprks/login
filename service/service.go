package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Password string `json:"password"`
	Verify string `json:"verify"`
	Crypt string
}

type Message struct {
	Message string `json:"message"`
}

type Identity struct {
	Ident struct {
		ID string `json:"id"`
		Registrations []struct {
			Plate string `json:"plate"`
		} `json:"registrations"`
	} `json:"identity"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	message := Message{}

	// Login
	if request.Resource == "/login" {
		fmt.Println("Start Login")
		resp, err := LoginService(request.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("login service err: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}

		r, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(fmt.Sprintf("login marshall: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}

		message.Message = string(r)
		fmt.Println("End Login")
	}

	// Register
	if request.Resource == "/register" {
		fmt.Println("Start Register")
		resp, err := RegisterService(request.Body)
		if err != nil {
			fmt.Println(fmt.Sprintf("register service err: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}

		r, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(fmt.Sprintf("register marshall: %v", err))
			return events.APIGatewayProxyResponse{}, err
		}

		message.Message = string(r)
		fmt.Println("End Register")
	}

	// Marshall the message
	m, err := json.Marshal(message)
	if err != nil {
		fmt.Println(fmt.Sprintf("message marshall: %v", err))
		return events.APIGatewayProxyResponse{}, err
	}

	// Return the message
	return events.APIGatewayProxyResponse{
		Body: string(m),
		StatusCode: 200,
	}, nil
}

func HashPassword(p string) (string, error) {
	fmt.Println(fmt.Sprintf("Hash Password: %v", p))
	r, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	fmt.Println(fmt.Sprintf("Hash-- R: %v, RS: %v, Err: %v", r, string(r), err))
	return string(r), err
}

func CheckPassword(p string, o string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(o))
	return err == nil
}
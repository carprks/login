package service

import (
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  uuid "github.com/satori/go.uuid"
  "golang.org/x/crypto/bcrypt"
)

// Login ...
type Login struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Error   error  `json:"error,omitempty"`
}

// LoginRequest ...
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest ...
type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Verify   string `json:"verify"`
	Crypt    string `json:"crypt,-"`
}

// Register ...
type Register struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Error error  `json:"error,omitempty"`
}

// Identity ...
type Identity struct {
	Ident struct {
		ID            string `json:"id"`
		Registrations []struct {
			Plate string `json:"plate"`
		} `json:"registrations"`
	} `json:"identity"`
}

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Login
	if request.Resource == "/login" {
		resp, err := login(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	// Register
	if request.Resource == "/register" {
		resp, err := register(request.Body)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       resp,
		}, nil
	}

	// Return the message
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
	}, nil
}

// HashPassword ...
func HashPassword(p string) (string, error) {
	r, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		fmt.Println(fmt.Sprintf("Hash err: %v", err))
		return "", err
	}
	return string(r), err
}

// CheckPassword ...
func CheckPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// GenerateIdent ...
func GenerateIdent(email string) string {
	u := uuid.NewV5(uuid.NamespaceURL, fmt.Sprintf("https://identity.carprk.com/user/%s", email))
	return u.String()
}

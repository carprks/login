package service

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/badoux/checkmail"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

func rest() (string, error) {
	return "", nil
}

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := rest()
	fmt.Println(fmt.Sprintf("Request Resource: %v", request.Resource))

	switch request.Resource {
	case "/login":
		resp, err = login(request.Body)

	case "/register":
		resp, err = register(request.Body)

	case "/remove":
		resp, err = delete(request.Body)
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("%v Err: %v", request.Resource, err))
		return events.APIGatewayProxyResponse{}, err
	}

	// Default
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       resp,
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

// CheckEmail ...
func CheckEmail(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}

	if os.Getenv("DEVELOPMENT") != "" {
		if email == "tester@carpark.ninja" || email == "testfail-login@carpark.ninja" {
			return nil
		}
	}
	err = checkmail.ValidateHost(email)
	if serr, ok := err.(checkmail.SmtpError); ok && err != nil {
		fmt.Println(fmt.Sprintf("Code: %v, Err: %v", serr.Code(), serr))

		if strings.Contains(serr.Error(), "Blocked") {
			fmt.Println(fmt.Sprintf("Email maybe fake but cant check, not users fault"))
			return nil
		}

		switch serr.Code() {
		case "550":
			return fmt.Errorf("invalid email")
		case "dia":
			return fmt.Errorf("invalid email")
		}

		fmt.Println(fmt.Sprintf("Unknown Code: %v", serr.Code()))
		return serr
	}

	return nil
}

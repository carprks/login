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

// EncPass Encoded Password
var EncPass = make(chan string)

func rest() (string, error) {
	return "", nil
}

// Handler ...
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := rest()

	switch request.Resource {
	case "/login":
		resp, err = login(request.Body)

	case "/register":
		resp, err = register(request.Body)

	case "/delete":
		resp, err = delete(request.Body)
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("%v Err: %v", request.Resource, err))
		return events.APIGatewayProxyResponse{}, err
	}

	if len(resp) == 0 {
		fmt.Println(fmt.Sprintf("Request Resource: %v", request.Resource))
		fmt.Println(fmt.Sprintf("Request Body %v", request.Body))
	}

	// Default
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       resp,
	}, nil
}

// HashPassword ...
func HashPassword(p string) {
	r, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		fmt.Println(fmt.Sprintf("Hash err: %v", err))
		return
	}

	EncPass <- string(r)
}

// CheckPassword ...
func CheckPassword(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	// 4fmt.Println(fmt.Sprintf("CheckPassword -- Password: %s, Err: %v", plainPassword, err))
	return err == nil
}

// GenerateIdent ...
func GenerateIdent(email string) string {
	u := uuid.NewV5(uuid.NamespaceURL, fmt.Sprintf("https://identity.carprk.com/user/%s", email))
	return u.String()
}

// CheckEmail ...
func CheckEmail(email string) error {
	if len(email) <= 2 {
		return fmt.Errorf("check email: email invalid")
	}

	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}

	if os.Getenv("DEVELOPMENT") != "" {
		if strings.Contains(email, "test") {
			return nil
		}
	}
	err = checkmail.ValidateHost(email)
	if serr, ok := err.(checkmail.SmtpError); ok && err != nil {
		fmt.Println(fmt.Sprintf("Code: %v, Err: %v", serr.Code(), serr))
		switch serr.Code() {
		case "550":
			return blockedCheck(serr.Error())
		case "dia":
			return fmt.Errorf("timeout err")
		}

		fmt.Println(fmt.Sprintf("Unknown Code: %v", serr.Code()))
		return serr
	}

	return nil
}

func blockedCheck(err string) error {
	if strings.Contains(err, "Blocked") {
		fmt.Println(fmt.Sprintf("Email maybe fake but cant check blocked, not users fault"))
		return nil
	} else if strings.Contains(err, "Spamhaus") {
		fmt.Println(fmt.Sprintf("Email maybe fake but cant check due to spamhaus, not users fault"))
		return nil
	}

	return fmt.Errorf("invalid email")
}

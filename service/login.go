package service

import (
	"encoding/json"
	"fmt"
)

// LoginService ...
func LoginService(body string) (string, error) {
	r := Login{}

	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Println(fmt.Sprintf("login unmarshall: %v", err))
		return "", err
	}

	fmt.Println("Create Crypt")
	crypt, err := HashPassword(r.Password)
	if err != nil {
		fmt.Println(fmt.Sprintf("crypt password: %v", err))
		return "", err
	}
	fmt.Println("Created Crypt")

	valid := CheckPassword(crypt, r.Password)
	if valid {
		return "success", nil
	}


	return "", nil
}

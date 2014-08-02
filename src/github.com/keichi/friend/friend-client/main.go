package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/keichi/friend/common"
	"net/http"
	"os"
	"strings"
)

func main() {
	createUserPrompt()
}

func createUserPrompt() {
	// if alreadyUser {
	// 	confirm
	// }

	// scan
	var name, password string
	s := bufio.NewScanner(os.Stdin)

	fmt.Print("New User's Name: ")
	s.Scan()
	name = strings.ToLower(s.Text())

	fmt.Print("Password: ")
	s.Scan()
	password = s.Text()

	user, err := CreateUser(name, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// StoreUser(user)
}

func CreateUser(name, password string) (*common.User, error) {
	// create user struct
	user := new(common.User)
	user.Name = name
	user.Password = password
	user.PublicKey = ""

	// convert to json
	data, err := json.Marshal(user)
	if err != nil {
		return user, err
	}

	// request to server
	response, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewReader(data))
	if err != nil {
		return user, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if response.StatusCode != 200 {
		body := buf.String()
		return user, errors.New(body)
	}

	err = json.Unmarshal(buf.Bytes(), user)
	if err != nil {
		return user, err
	}

	return user, nil
}

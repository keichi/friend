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
	osuser "os/user"
	"strings"
)

func main() {
	createUserPrompt()
}

func createUserPrompt() {
	// scan
	var name, password string
	s := bufio.NewScanner(os.Stdin)

	if _, err := os.Stat(getHomePath() + "/.friend/user.json"); !os.IsNotExist(err) {
		fmt.Print("Already user data exists. Create new user? [y/N]: ")
		s.Scan()
		if s.Text() != "y" {
			return
		}
	}

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

	StoreUser(user)
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

func StoreUser(user *common.User) {
	if err := os.MkdirAll(getHomePath()+"/.friend", 0700); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	file, err := os.OpenFile(getHomePath()+"/.friend/user.json", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err = json.NewEncoder(file).Encode(user); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getHomePath() string {
	user, err := osuser.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return user.HomeDir
}

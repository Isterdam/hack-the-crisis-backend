package api

import (
	"io/ioutil"
	"strings"
)

var JWTkey string
var PhoneUser string
var PhonePass string

// initialize very secret constants from local environment
func Initialize_constants() {
	JWTkeyTemp, _ := ioutil.ReadFile("secretJWTKey.txt")
	JWTkey = string(JWTkeyTemp)

	phoneNumTemp, _ := ioutil.ReadFile("phoneNum.txt")
	temp := strings.Split(string(phoneNumTemp), "-")
	PhoneUser = temp[0]
	PhonePass = temp[1]
}
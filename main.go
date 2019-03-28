package main

import (
"fmt"
"github.com/dgrijalva/jwt-go"
"net/http"
"encoding/json"
"bytes"
)

type MyCustomClaims struct {
	Scope string `json:"scope,omitempty"`
	jwt.StandardClaims
}

type Toke struct {
	Access string `json:"access_token,omitempty"`
	Type string `json:"token_type,omitempty"`
	Expire string `json:"expires_in,omitempty"`
}

func main() {

	key := []byte("<your private key>")

	key1, _ := jwt.ParseRSAPrivateKeyFromPEM(key)

	claims := MyCustomClaims{
		"https://www.googleapis.com/auth/cloudprint",
		jwt.StandardClaims{
			IssuedAt: <currrent-epoch-time>,         // eg 1234566000
			ExpiresAt: <currrent-epoch-time + 3600>, // 3600 secs = 1hour, so expires in 1 hour, eg 1234569600
			Issuer:    "<your service account email>",
			Audience: "https://www.googleapis.com/oauth2/v4/token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ss)
	url := "https://www.googleapis.com/oauth2/v4/token"
	any := "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer&assertion=" + ss
	a := []byte(any)
	b := bytes.NewBuffer(a)
	var tok Toke

	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	} else {
		json.NewDecoder(resp.Body).Decode(&tok)
	}
	fmt.Println("----------- Access Token -----------------")
	fmt.Println("Access: ", tok.Access)
}


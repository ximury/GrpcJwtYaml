package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)
/*
jwt测试
 */
const (
	SecretKey = "welcome here!"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {

	http.HandleFunc("/login", LoginHandler)
/*
   negroni 是为了方便使用 net/http 而设计的一个库
 */
	http.Handle("/resource", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}

func main() {
	//str:="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	//	"eyJleHAiOjE1OTIzNzYxMjcsImlhdCI6MTU5MjM3MjUyN30." +
	//	"bsQeAHV9Xy1h487iz1LIbaq6HXxejYzd_OgtbQbuI0Q"
	//base:="{\"alg\":\"HS256\",\"typ\":\"JWT\"}{\"exp\":1592376127,\"iat\":1592372527}"
	//fmt.Println(base64.StdEncoding.DecodeString(str))
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(base)))
	StartServer()
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "someone" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in, username is wrong")
		fmt.Fprint(w, "Invalid credentials")
		return
	}
	if user.Password != "password" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in, password is wrong")
		fmt.Fprint(w, "Invalid credentials")
		return
	}
	fmt.Println("Success!")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

/**
POST:
http://127.0.0.1:8080/login
GET:
http://127.0.0.1:8080/resource
Bearer +token
 */
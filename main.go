package main

import (
	"fmt"
	"net/http"

	"io/ioutil"

	"bytes"

	"encoding/json"

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

var token string

func userHandler(c *gin.Context) {
	url := "https://api.instagram.com/v1/users/self/?access_token=" + token
	fmt.Printf("Requesting %s", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		c.Abort()
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		c.Abort()
		return
	}
	c.String(http.StatusOK, string(body))
}

//PostResponse its going to change
type PostResponse struct {
	AccessToken string `json:"access_token"`
}

func instagramAuthHandler(c *gin.Context) {
	clientID := "a9e0f62b9ab844cda85194bfcaa649df"
	clientSecret := "2b4cd2664a2e405faec355a8e9f6a3fc"
	grantType := "authorization_code"
	redirectURI := "http://localhost:8080/instagram_redirect"
	code := c.Query("code")
	postTokenURL := "https://api.instagram.com/oauth/access_token"

	var postData bytes.Buffer
	w := multipart.NewWriter(&postData)

	fw, _ := w.CreateFormField("client_id")
	fw.Write([]byte(clientID))

	fw, _ = w.CreateFormField("client_secret")
	fw.Write([]byte(clientSecret))

	fw, _ = w.CreateFormField("grant_type")
	fw.Write([]byte(grantType))

	fw, _ = w.CreateFormField("redirect_uri")
	fw.Write([]byte(redirectURI))

	fw, _ = w.CreateFormField("code")
	fw.Write([]byte(code))

	w.Close()

	postResponse, _ := http.Post(postTokenURL, w.FormDataContentType(), &postData)
	postResponseData := &PostResponse{}
	postBody, _ := ioutil.ReadAll(postResponse.Body)
	err := json.Unmarshal(postBody, &postResponseData)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, postResponseData.AccessToken)
}

func getTokenHandler(c *gin.Context) {
	clientID := "a9e0f62b9ab844cda85194bfcaa649df"
	redirect := "http://localhost:8080/instagram_redirect"
	responseType := "code"
	getCodeURL := "https://api.instagram.com/oauth/authorize/?client_id=" + clientID + "&redirect_uri=" + redirect + "&response_type=" + responseType
	//response, _ := http.Get(getCodeURL)
	c.Redirect(301, getCodeURL)
	//GET -> https://api.instagram.com/oauth/authorize/?client_id=CLIENT-ID&redirect_uri=REDIRECT-URI&response_type=code
	//CON RESPUESTA
	//POST -> https://api.instagram.com/oauth/access_token
	//CON
	/*
			'client_id=CLIENT_ID'
		    'client_secret=CLIENT_SECRET'
			'grant_type=authorization_code'
		    'redirect_uri=AUTHORIZATION_REDIRECT_URI'
		    'code=CODE'
	*/
	//data, _ := ioutil.ReadAll(response.Body)
	//c.Data(http.StatusOK, "text/html", data)
}

func postTokenHandler(c *gin.Context) {
	newToken := c.PostForm("auth_token")
	if newToken != "" {
		token = newToken
		c.Request.Method = "GET"
		c.Redirect(301, "/home")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe ingresar un Token valido"})
	}
}

func setUpServer() *gin.Engine {
	r := gin.Default()
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/home", "./frontend/home.html")
	r.StaticFile("/frontapp", "./frontend/js/app.js")
	r.POST("/ingresar_token", postTokenHandler)
	r.GET("/token", getTokenHandler)
	r.GET("/user", userHandler)
	r.GET("/instagram_redirect", instagramAuthHandler)
	return r
}
func main() {
	r := setUpServer()
	r.Run()
}

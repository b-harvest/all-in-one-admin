package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/b-harvest/all-in-one-admin/config"
	"github.com/b-harvest/all-in-one-admin/structs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
)

var cred Credentials
var conf *oauth2.Config
var conn *grpc.ClientConn
var conn2 *grpc.ClientConn

// Credentials which stores google ids.
type Credentials struct {
	Cid         string `json:"cid"`
	Csecret     string `json:"csecret"`
	RedirectURL string `json:"url"`
}

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal data")
		return
	}

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  cred.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
	grpcclient()
	grpcclient_v2()
}

// AuthHandler handles authentication of a user and initiates a session.
func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}
	code := c.Request.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := structs.User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}
	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	c.HTML(http.StatusOK, "validator.tmpl", gin.H{"link": "/"})
}

// LoginHandler handles the login procedure.
func LoginHandler(c *gin.Context) {
	state, err := RandToken(32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}
	link := getLoginURL(state)
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link})
}

func grpcclient() {
	conn_v, err := grpc.Dial("localhost:8088", grpc.WithInsecure(), grpc.WithBlock())
	conn = conn_v
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func grpcclient_v2() {
	conn_v2, err := grpc.Dial("localhost:8089", grpc.WithInsecure(), grpc.WithBlock())
	conn2 = conn_v2
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

//You need to change it to be a little more communication efficient.
func GetnodeStatusHandler(c *gin.Context) {
	if conn == nil {
		grpcclient()
	}
	log.Println(conn)
	connect := pb.NewMonitoringClient(conn)
	nodeuri := c.Query("nodeuri")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := connect.GetnodeStatus(ctx, &pb.StatusRequest{NodeURI: nodeuri})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": r.Status})
	}
	log.Printf("Config: %v", r)
}

func GetvalidatorSignInfo(c *gin.Context) {
	if conn == nil {
		grpcclient()
	}
	log.Println(conn)
	connect := pb.NewMonitoringClient(conn)
	nodeuri := c.Query("nodeuri")
	validator := c.Query("validatoraddress")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q, err := connect.GetvalidatorSignInfo(ctx, &pb.SignInfoRequest{NodeURI: nodeuri, ValidatorAddress: validator})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": q.Status})
	}
	log.Printf("Config: %v", q)
}

func GetnodeStatusHandler_v2(c *gin.Context) {
	if conn2 == nil {
		grpcclient_v2()
	}
	log.Println(conn2)
	connect := pb.NewMonitoringClient(conn2)
	nodeuri := c.Query("nodeuri")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := connect.GetnodeStatus(ctx, &pb.StatusRequest{NodeURI: nodeuri})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": r.Status})
	}
	log.Printf("Config: %v", r)
}

func GetvalidatorSignInfo_v2(c *gin.Context) {
	if conn2 == nil {
		grpcclient_v2()
	}
	log.Println(conn2)
	connect := pb.NewMonitoringClient(conn2)
	nodeuri := c.Query("nodeuri")
	validator := c.Query("validatoraddress")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	q, err := connect.GetvalidatorSignInfo(ctx, &pb.SignInfoRequest{NodeURI: nodeuri, ValidatorAddress: validator})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": q.Status})
	}
	log.Printf("Config: %v", q)
}

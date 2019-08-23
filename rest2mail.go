package main

import (
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"reflect"
	"strconv"
	"strings"
)



type MailRequest struct {
	//Mail Server and Port address
	Server string   `binding:"required" json:"server"`
	Port int      `binding:"required,gte=0,lte=65535" json:"port"`
	From string   `binding:"required,email" json:"from"`
	To []string `binding:"required,emails,gte=1,lte=100" json:"to"`
	Subject string   `binding:"required" json:"subject"`
	Message string   `binding:"required" json:"message"`
	Hostname string   `binding:"required" json:"hostname"`
	Username string   `binding:"required" json:"username"`
	Password string   `binding:"required" json:"password"`
	Identity string   `json:"identity"`
}

func main() {
	r := StartServer()
	var serverDomain = os.Getenv("DOMAIN")
	if serverDomain != "" { //If Domain is setted run under http using let's encrypt
		println("Running HTTPS server under domain ", serverDomain)
		log.Fatal(autotls.Run(r,serverDomain))
	}else {
		r.Run(":8080")
	}
}

func StartServer() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("emails", ValidateMailList)
	}
	r.POST("/send", send)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{"status":"OK"})
	})
	return r
}

func send(c *gin.Context) {
	var mailReq MailRequest
	if err := c.ShouldBindJSON(&mailReq); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Actually send the mail using go smtp client

	//Get the auth
	auth := smtp.PlainAuth(mailReq.Identity, mailReq.Username, mailReq.Password, mailReq.Hostname)
	//Send the mail generating a message body with a RFC 822-style email with headers first
	err := smtp.SendMail(mailReq.Server+":"+strconv.Itoa(mailReq.Port),auth, mailReq.From, mailReq.To,
		[]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",mailReq.From,strings.Join(mailReq.To,","),mailReq.Subject,mailReq.Message)))
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func ValidateMailList(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
	mails := field.Interface().([]string)
	for _,m := range mails{
		if err := v.Field(m, "required,email"); err != nil{
			return false
		}
	}
	return true
}


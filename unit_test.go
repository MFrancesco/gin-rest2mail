package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func Test(t *testing.T){
	router := StartServer()
	invalidMailRequest := MailRequest{
		Server:   "smtp.sendgrid.net",
		Port:     587,
		From:     "sender@test.com",
		To:       []string{"invalid","mail@gmail.com"},
		Subject:  "Test",
		Message:  "Test message",
		Hostname: "hostname",
		Username: "username",
		Password: "pwd",
		Identity: "",
	}
	resp := performRequest(invalidMailRequest,router)
	if resp.Code != 400{
		t.Error()
	}

	invalidMailRequest = MailRequest{
		Server:   "smtp.sendgrid.net",
		Port:     587,
		From:     "sender@test.com",
		To:       []string{"valid@gmail.com","mail@gmail.com"},
		Subject:  "Test",
		Message:  "Test message",
		Hostname: "hostname",
		Username: "username",
		Password: "",//empy passwd so will crash
		Identity: "",
	}
	resp = performRequest(invalidMailRequest,router)
	if resp.Code != 400{
		t.Error()
	}

	res, _ := ioutil.ReadAll(resp.Body); fmt.Println(string(res))

	validRequest := MailRequest{
		Server:   "smtp.sendgrid.net",
		Port:     587,
		From: "sender@test.com",
		To: []string{"test1@gmail.com","test2@gmail.com"},
		Subject: "Subject",
		Message: "Message",
		Hostname: "hostname",
		Username: "username",
		Password: "test",
		Identity: "",
	}
	resp = performRequest(validRequest,router)
	res, _ = ioutil.ReadAll(resp.Body); fmt.Println(string(res))

	if resp.Code != 400{
		t.Error()
	}
	if strings.Contains(string(res),"wrong host name") == false{
		t.Error()
	}
	/*
	This has been used with sendgrid to actually check if mail are being sent, and they are
	validRequest := MailRequest{
		Server:   "smtp.sendgrid.net",
		Port:     587,
		From: "sender@test.com",
		To: []string{"my_personal_email@gmail.com"},
		Subject: "Subject",
		Message: "Message",
		Hostname: "smtp.sendgrid.net",
		Username: "apikey",
		Password: "myapikey",
		Identity: "",
	}
	jsn,err = json.Marshal(validRequest)
	resp = performRequest(router,"POST","/send", bytes.NewBuffer(jsn))
	res, _ = ioutil.ReadAll(resp.Body)
	if resp.Code != 200{
		t.Error()
	}
	*/
}


func performRequest(mr MailRequest , r http.Handler) *httptest.ResponseRecorder {
	j,_ := json.Marshal(mr)
	req, _ := http.NewRequest("POST","/send", bytes.NewBuffer(j))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}


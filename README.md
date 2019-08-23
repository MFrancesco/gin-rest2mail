# gin-rest2mail

Play project used to write some go.
 
It expose an health endpoint and a POST endpoint that can be used to send email
if the MailRequest model is valid and auth credentials are correct

- Run it

`go run rest2mail.go`

- Check health

`curl localhost:8080/health`

- Send mail using post 
```bash
curl --header "Content-Type: application/json" \
       --request POST \
       --data '{"Server": "smtp.sendgrid.net", "Port": 587, "From": "sender@test.com", "To": ["my_personal_email@gmail.com"], "Subject": "Subject", "Message": "Message", "Hostname": "smtp.sendgrid.net","Username": "apikey", "Password": "myapikey", "Identity": ""}' \
    localhost:8080/send
``` 
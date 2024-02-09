package util

import (
	"encoding/json"
	"log"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSMS(toNum string, message string) (string, error) {
	config, err := LoadConfig("../")
	if err != nil {
		log.Fatal("[sendSMS] cannot load config: ", err)
	}

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.TwiolioAccountSid,
		Password: config.TwiolioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(toNum)
	params.SetFrom(config.TwiolioFrom)
	params.SetBody(message)

	resp, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		log.Fatal("Error sending SMS message: " + err.Error())
		return "", err
	} else {
		response, _ := json.Marshal(*resp)
		log.Println("Response: " + string(response))
		return string(response), nil
	}
}

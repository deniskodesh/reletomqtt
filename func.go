package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var logs = logrus.New()

func getURLData(url string, token string) (*http.Response, []byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getURLData",
			"operation": "Create NewRequest",
		}).Errorln(err)
		return nil, nil, err
	}
	req.Header.Add("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getURLData",
			"operation": "client.Do(req)",
		}).Errorln(err)
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getURLData",
			"operation": "ioutil.ReadAll(resp.Body)",
		}).Errorln(err)
		return nil, nil, err
	}

	return resp, body, nil
}

func getURLDataWithRetries(respB Bakler, operation string) (*http.Response, []byte, error, string) {

	var body []byte
	var err error
	var resp *http.Response
	var url string

	var backoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}

	if operation == "voltage" {
		url = "https://buckler.pw/api/relay/" + respB.Relays[0].SerialNumber + "/readings?from=" + getTimeStamps("from") + "&to=" + getTimeStamps("to") + "&type=voltage"
	}

	if operation == "power" {
		url = "https://buckler.pw/api/relay/" + respB.Relays[0].SerialNumber + "/readings?from=" + getTimeStamps("from") + "&to=" + getTimeStamps("to") + "&type=power"
	}

	if operation == "current" {
		url = "https://buckler.pw/api/relay/" + respB.Relays[0].SerialNumber + "/readings?from=" + getTimeStamps("from") + "&to=" + getTimeStamps("to") + "&type=current"
	}

	for i, backoff := range backoffSchedule {

		resp, body, err = getURLData(url, respB.Token)
		logs.WithFields(logrus.Fields{
			"func":      "getURLDataWithRetries",
			"operation": "Get len of body byte slice",
			"iteration": i,
		}).Info("Slice size is " + strconv.Itoa(len(body)))
		if err == nil {
			logs.WithFields(logrus.Fields{
				"func":      "getURLDataWithRetries",
				"operation": "Check if no error and body slice size",
				"iteration": i,
			}).Info("Slice size is " + strconv.Itoa(len(body)))
			defer resp.Body.Close()
			break
		}
		logs.WithFields(logrus.Fields{
			"func":      "getURLDataWithRetries",
			"operation": "All retries failed",
		}).Info("Request error: %+v\n", err)

		logs.WithFields(logrus.Fields{
			"func":      "getURLDataWithRetries",
			"operation": "All retries failed",
		}).Info("Retrying in %v\n", backoff)

		time.Sleep(backoff)
	}

	// All retries failed
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getURLDataWithRetries",
			"operation": "All retries failed",
		}).Error(err)
		return nil, nil, err, operation
	}

	return resp, body, nil, operation
}

//Get JWT token from Bakler portal
func getJWTtoken(cr Creds, url string) Bakler {

	//Marshall struct

	sJson, err := json.Marshal(cr)

	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getJWTtoken",
			"operation": "json.Marshal(cr)",
		}).Error(err)

	}

	var jsonStr = []byte(sJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getJWTtoken",
			"operation": "http.NewRequest()",
		}).Error(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getJWTtoken",
			"operation": "client.Do(req)",
		}).Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "getJWTtoken",
			"operation": "ioutil.ReadAll(resp.Body)",
		}).Error(err)
	}

	req.Header.Set("Content-Type", "application/json")

	var tmpResp Bakler
	json.Unmarshal([]byte(body), &tmpResp)

	return tmpResp

}

// Get actual time stamp for the end of the day
func getStartandEndoftheDay(optionD string) (day string) {

	if optionD == "start" {
		now := time.Now()
		tmp := Bod(now)
		dayStart := strconv.FormatInt(tmp.UnixNano()/1000000, 10)

		return dayStart
	}

	if optionD == "end" {
		now := time.Now()
		tmp := Bod(now)
		dayEnd := strconv.FormatInt((tmp.UnixNano()/1000000)+86399999, 10)
		return dayEnd
	}

	return "error"

}

//Get time stamp for some day in miliseconds
func getTimeStamps(optionD string) (day string) {

	if optionD == "to" {
		now := time.Now()
		to := strconv.FormatInt(now.UnixNano()/1000000, 10)

		return to
	}

	if optionD == "from" {
		now := time.Now()
		from := strconv.FormatInt(now.UnixNano()/1000000-60000, 10)
		return from
	}

	return "error"

}

//Get stamp for the day start
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

//Get voltage from bakler
func getVoltage(body []byte) []Voltage {

	var volts []Voltage

	var respVolts BaklerVoltage

	json.Unmarshal([]byte(body), &respVolts)

	if len(respVolts.Data) > 0 {

		for _, tmp := range respVolts.Data {

			volts = append(volts, tmp)

		}

	} else {

		volts = nil
	}

	return volts
}

//Get amperage from bakler
func getAmperage(body []byte) []Amperage {

	var amperage []Amperage

	var respAmperage BaklerAmperage

	json.Unmarshal([]byte(body), &respAmperage)

	if len(respAmperage.Data) > 0 {

		for _, tmp := range respAmperage.Data {

			amperage = append(amperage, tmp)

		}
	} else {

		amperage = nil
	}

	return amperage
}

//Get power from bakler
func getPower(body []byte) []Power {

	var power []Power

	var respPower BaklerPower

	json.Unmarshal([]byte(body), &respPower)

	if len(respPower.Data) > 0 {

		for _, tmp := range respPower.Data {

			power = append(power, tmp)

		}
	} else {

		power = nil
	}
	return power
}

//Preparing struct to send in mqtt functions
func createStructToSend(voltage []Voltage, power []Power, amperage []Amperage) TooSend {

	var readyTosend TooSend

	logs.WithFields(logrus.Fields{
		"func":      "createStructToSend",
		"operation": "Len of voltage struct",
	}).Info(strconv.Itoa(len(voltage)))

	logs.WithFields(logrus.Fields{
		"func":      "createStructToSend",
		"operation": "Len of power struct",
	}).Info(strconv.Itoa(len(power)))

	logs.WithFields(logrus.Fields{
		"func":      "createStructToSend",
		"operation": "Len of amperage struct",
	}).Info(strconv.Itoa(len(amperage)))

	for _, tmp := range voltage {

		readyTosend.Voltage = append(readyTosend.Voltage, tmp)

	}

	for _, tmp := range power {

		readyTosend.Power = append(readyTosend.Power, tmp)

	}

	for _, tmp := range amperage {

		readyTosend.Amperage = append(readyTosend.Amperage, tmp)

	}

	return readyTosend

}

//This func publish data on some topic
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	logs.WithFields(logrus.Fields{
		"func":      "messagePubHandler",
		"operation": "Publish",
	}).Info("Received message: " + string(msg.Payload()) + "from topic: " + msg.Topic())
}

//This func return string if all ok
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	logs.WithFields(logrus.Fields{
		"func":      "connectHandler",
		"operation": "Connection",
	}).Info("Connected")
}

//This func return string if connection broken
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {

	logs.WithFields(logrus.Fields{
		"func":      "connectLostHandler",
		"operation": "lost connection ",
	}).Error(err)
}

//This func publush data to some topic
func publish(client mqtt.Client, metrics string, topic string) {

	token := client.Publish(topic, 0, false, metrics)
	token.Wait()
	time.Sleep(time.Second)
}

//This func create mqtt client
func sentToMqtt(tosend TooSend, mqttSettings MqttSetting, id int, wg *sync.WaitGroup) {
	logs.WithFields(logrus.Fields{
		"func":      "sentToMqtt",
		"operation": "go routine",
	}).Info("Worker starting id=", id)

	sJson, err := json.Marshal(tosend)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"func":      "sentToMqtt",
			"operation": "json.Marshal(tosend)",
		}).Error(err)
		return
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttSettings.ConnectionString)
	opts.SetUsername(mqttSettings.Username)
	opts.SetPassword(mqttSettings.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	clientL := mqtt.NewClient(opts)
	if token := clientL.Connect(); token.Wait() && token.Error() != nil {

		logs.WithFields(logrus.Fields{
			"func":      "sentToMqtt",
			"operation": "json.Marshal(tosend)",
		}).Error(err)
	}

	if mqttSettings.ConnectionString == "tcp://mqtt.lynus.io:1883" {

		var toSend Lynustype

		if len(tosend.Voltage) > 0 {
			voltage := Lynustypeone{"Voltage", tosend.Voltage[0].Voltage.Actual}
			toSend = append(toSend, voltage)
		}

		if len(tosend.Power) > 0 {
			power := Lynustypeone{"Power", tosend.Power[0].Power.Actual}
			toSend = append(toSend, power)
		}

		if len(tosend.Voltage) > 0 {
			amperage := Lynustypeone{"Amperage", tosend.Amperage[0].Current.Actual}
			toSend = append(toSend, amperage)
		}

		sJson, err = json.Marshal(toSend)
		if err != nil {
			logs.WithFields(logrus.Fields{
				"func":      "sentToMqtt",
				"operation": "json.Marshal(tosend)",
			}).Error(err)
			return
		}

	}

	publish(clientL, string(sJson), mqttSettings.Topic)

	clientL.Disconnect(3)
	defer wg.Done()

	logs.WithFields(logrus.Fields{
		"func":      "sentToMqtt",
		"operation": "go routine ",
	}).Info("Worker done id=", id)

}

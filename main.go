package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {

	var volts []Voltage
	var power []Power
	var amperage []Amperage

	operations := [...]string{"voltage", "current", "power"}

	var wg sync.WaitGroup

	var body []byte
	log.WithFields(log.Fields{
		"Starting ": " App",
	}).Info("*****************")

	//reading config from file
	var settings Config

	data, err := ioutil.ReadFile("/etc/hometolynus/settings.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(data), &settings)

	respB := getJWTtoken(settings.Credentials, settings.Authurl)

	for {

		now := time.Now()

		if respB.Token == "" || respB.ExpireAt <= int(now.UnixNano()/1000000) {

			fmt.Println("Send request to get token it is from condition ")

			s := strconv.Itoa(respB.ExpireAt)
			fmt.Println("From responce: " + s)

			s = strconv.Itoa(int(now.UnixNano() / 1000000))
			fmt.Println("Time stamp: " + s)
			respB = getJWTtoken(settings.Credentials, settings.Authurl)
		}
		for i, operation := range operations {

			_, body, _, operation = getURLDataWithRetries(respB, operation)

			logs.WithFields(logrus.Fields{
				"func":          "main",
				"operation":     "Range of operations",
				"iteration":     i,
				"operation opt": operation,
			}).Info("Slice size is " + strconv.Itoa(len(body)))

			if operation == "voltage" {
				volts = getVoltage(body)

				log.WithFields(log.Fields{
					"Struct slice size ": " Voltage",
				}).Info("struct size is " + strconv.Itoa(len(volts)))

			}

			if operation == "power" {
				power = getPower(body)

				log.WithFields(log.Fields{
					"Struct slice size ": " Power",
				}).Info("struct size is " + strconv.Itoa(len(volts)))

			}

			if operation == "current" {
				amperage = getAmperage(body)

				log.WithFields(log.Fields{
					"Struct slice size ": " Voltage",
				}).Info("struct size is " + strconv.Itoa(len(volts)))

			}
		}

		toSend := createStructToSend(volts, power, amperage)

		for i, tmp := range settings.MqttSettings {

			wg.Add(1)

			go sentToMqtt(toSend, tmp, i, &wg)

		}

		time.Sleep(60 * time.Second)

	}

}

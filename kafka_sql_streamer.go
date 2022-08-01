package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"../../pkg_kafka"
	"../../pkg_mysql"
)

// Running Variables
var RunningProcessMsg = false
var Terminate = false

func main() {
	// ConnectDB(blah)
	err := sql.Connect("root", "password", "mydb", "127.0.0.1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sql.Close()

	log("pass db")

	// ConnectKafka(blah)
	err = kafka.Connect("localhost:9092")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log("pass kafka")
	// read kafka starting point from the user
	// read skip point from user
	// confirm starting and skip points to the user
	// read settings
	// setup msgbus
	// setup log files

	// Setup the Kafka listening queue
	receiveChannel := make(chan []byte)
	go func() {
		err = kafka.Receive("mysqlUpdateQueue", receiveChannel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	log("pass kafka receive")

	// Process the files
	var msg []byte
	go func() {
	Loop:
		for {
			log("starting loop")
			switch Terminate {
			case false:
				msg = <-receiveChannel
				RunningProcessMsg = true

				// Execute Query
				log("Processing Msg " + string(msg))
				myval := make(map[int]map[string]string)
				myval, err := sql.FetchAny(string(msg))
				// check for errors
				if err != nil {
					log(err.Error())
				}
				// display any returned results, you should be having any
				for i := range myval {
					for k, v := range myval[i] {
						log("Returned " + k + " : " + v)
					}
				}
				// Query Ended
				RunningProcessMsg = false
			case true:
				break Loop
			}
		}
	}()

	//
	// Termination Loop
	//

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	// Tell all loops to terminate
	Terminate = true
	crt("Termination Received, Terminating Nicely - Max 120s")

	for x := 0; x < 1; x++ {
		time.Sleep(time.Second * 1)
		if RunningProcessMsg == false {
			log("Clean Exit")
			// output last kafka msg id.
			// log last kafka msg id
			os.Exit(0)
		}
	}
	fmt.Println("Timeout Reached Dirty Exit")
	// output last kafka msg id.
	// log last kafka msg id
	crt("Timeout Reached Dirty Exit")
	os.Exit(1)
}

func log(s string) {
	fmt.Println(s)
}

func war(string) {

}

func crt(string) {

}

package main

import (
	//"io"
	"fmt"
	"time"
	"os"
	"bufio"

	"github.com/tarm/serial"
	"go.bug.st/serial.v1/enumerator"
)
var next = time.Now().Add(30 * time.Second)
var stomp = make(chan bool)
var datastream = make(chan string)

func timestamp(eventchan chan bool) {
	for true {
		time.Sleep(5 * time.Second)
		eventchan <- true
	}
}

func serialreceive(receivechan chan string) {
	serialPort := ""
	if 2 == len(os.Args) {
		serialPort = os.Args[1]
		
	} else {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			fmt.Printf("Can't get serial port list:", err)
			os.Exit(-1)
		}
		if 0 == len(ports) {
			fmt.Println("No serial ports available")
			os.Exit(-1)
		} else { // take the first available port in the list
			serialPort = ports[0].Name
		}
	}
	fmt.Println("opening port:", serialPort)
	conf := &serial.Config{Name: serialPort, Baud: 115200}
	sp, err := serial.OpenPort(conf)
	if err != nil {
		fmt.Println("Unable to open serial port", serialPort)
		os.Exit(-1)
	}
	defer sp.Close()
	scanner := bufio.NewScanner(sp)
	for true {
		if scanner.Scan() {
			receivechan <- scanner.Text()
		}
	}
}

func main() {
	go timestamp(stomp)
	go serialreceive(datastream)
	for true {
		select {
			case rxmessage := <- datastream:
				fmt.Println(rxmessage)
			case <- stomp:
				currenttime := time.Now()
				fmt.Println(currenttime.Format("========== TIMESTAMP ========== 2006-01-02 15:04:05 =========="))
		}
	}
}
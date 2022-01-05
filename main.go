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

func main() {
	next := time.Now().Add(5 * time.Second)
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
	conf := &serial.Config{Name: serialPort, Baud: 115200}
	sp, err := serial.OpenPort(conf)
	if err != nil {
		fmt.Sprintf("Unable to open serial port %s", serialPort)
		os.Exit(-1)
	}
	defer sp.Close()
	scanner := bufio.NewScanner(sp)
	for true {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		} else {
			currenttime := time.Now()
			if currenttime.After(next) {
				next = time.Now().Add(5 * time.Second)
				fmt.Println(currenttime.Format("2006-01-02 15:04:05"))
			}
		}		
	}
}
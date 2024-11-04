package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Shtokarev/raspberry-telegram-bot/config"
	"github.com/Shtokarev/raspberry-telegram-bot/sim900"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	conf := config.New()

	fmt.Println("SERIAL_PORT:", conf.SerialPort)
	fmt.Println("DEBUG_MODE:", conf.DebugMode)
	fmt.Println("BAUD_RATE:", conf.BaudRate)
	fmt.Println("AUTODETECT_BAUD_RATE:", conf.AutodetectBaudRate)

	gsm := sim900.New(conf.SerialPort, conf.BaudRate, conf.DebugMode)

	err := gsm.Connect()
	if err != nil {
		panic(err)
	}

	if conf.AutodetectBaudRate {
		maxBaudRate, err := gsm.GetMaxBaudRate()
		if err == nil && maxBaudRate > config.MIN_BAUD_RATE {
			if maxBaudRate > config.MAX_BAUD_RATE {
				maxBaudRate = config.MAX_BAUD_RATE
			}

			fmt.Printf("Detected max baud rate speed %d, reconnecting...\n", maxBaudRate)
			err := gsm.SetMaxBaudRate(maxBaudRate)

			if err == nil {
				gsm.Disconnect()

				gsm = sim900.New(conf.SerialPort, maxBaudRate, conf.DebugMode)
				err := gsm.Connect()
				if err != nil {
					panic(err)
				}
				fmt.Printf("Switched to %d successfully\n", maxBaudRate)
			}
		}
	}

	defer gsm.Disconnect()

	time.Sleep(time.Millisecond * 500)

	// gsm.SendSMS("+79185564752", "Test message 8", "")
	data, err := gsm.GET("catfact.ninja")
	if err != nil {
		panic(err)
	}

	fmt.Println("---------------->")
	fmt.Println(data)

	// return

	// fmt.Println("-------- GET SMS List")

	// gsm.GetSMSList(sim900.SMS_UNREADED)
	// list, err := gsm.GetSMSList(sim900.SMS_ALL)
	// fmt.Println("-------- List msg:", list)

	// msg, err := gsm.ReadSMS("2")
	// err = gsm.DeleteSMS("5")
	// if err != nil {
	// 	fmt.Printf("Error happens %s", err.Error())
	// }

	// fmt.Printf("Message: %s", msg)

	return

	fmt.Println("-------- SMS")
	// gsm.SendSMS("+79185564752", "Test message 3 (русский текст)", sim900.CHSET_UCS2)
	gsm.SendSMS("+79185564752", "Test message 6 (english)", "")

	// time.Sleep(time.Second * 1)

	// phoneNumber := "XXXXXXXXXX" // The number to send the SMS
	// gsm.SendSMS(phoneNumber, "Hello World!")
	fmt.Println("-------- FINISH!!")

}

// func main2() {
// 	config := &serial.Config{
// 		Name:        "/dev/ttyAMA0",
// 		Baud:        9600,
// 		ReadTimeout: 0,
// 		Size:        8,
// 	}

// 	fmt.Println(config)

// 	port, err := serial.OpenPort(config)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	n, err := port.Write([]byte("AT\r"))
// 	// n, err := port.Print("AT\r\n")
// 	fmt.Println("Written", n)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	buf := make([]byte, 256)
// 	//n, err = port.Read(buf)
// 	//if err != nil {
// 	//   log.Fatal(err)
// 	// }
// 	// fmt.Println("Readen", n)
// 	// fmt.Println(string(buf))

// 	// Send the binary integer 4 to the serial port
// 	// buf := make([]byte, "AT\r")
// 	// fmt.printLn(buf)
// 	// res, err = connection.Write(buf)
// 	// fmt.printLn(res)
// 	// if err != nil {
// 	//   log.Fatal(err)
// 	// }

// 	// buf := make([]byte, 1024)

// 	for {
// 		n, err := port.Read(buf)

// 		if err != nil {
// 			fmt.Println("ERROR:")
// 			fmt.Println(err)
// 			log.Fatal(err)
// 		}

// 		fmt.Println("No error after reading")
// 		s := string(buf[:n])
// 		fmt.Println("Readen n:", n)
// 		fmt.Println("String(s):", string(s))
// 		fmt.Println("s:", s)
// 	}
// }

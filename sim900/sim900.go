package sim900

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type SIM900 struct {
	port      *serial.Port
	config    *serial.Config
	logger    *log.Logger
	APN       string
	debugMode bool
}

type SMS struct {
	Id       int
	Status   string
	Operator string
	Datetime string
	Message  []byte
}

func New(name string, baud int, debugMode bool) *SIM900 {
	return &SIM900{
		config: &serial.Config{
			Name:        name,
			Baud:        baud,
			ReadTimeout: 0,
			Size:        8,
		},
		debugMode: debugMode,
		logger:    log.New(os.Stdout, "[SIM900] ", log.LstdFlags),
	}
}

func (sim *SIM900) Connect() error {
	var err error

	sim.port, err = serial.OpenPort(sim.config)
	if err != nil {
		sim.logger.Fatal("Not possible to open port", err)
		return err
	}

	// check connection with modem
	err = sim.At()
	if err != nil {
		return err
	}

	if sim.debugMode {
		sim.EnableEcho()
	} else {
		sim.DisableEcho()
	}

	atDebugCommand := CMD_DEBUG_LOG_OFF

	if sim.debugMode {
		atDebugCommand = CMD_DEBUG_LOG_ON
	}

	if _, err := sim.SendCommand(atDebugCommand); err != nil {
		return err
	}

	err = sim.DisplayModemInfo()

	return err
}

func (sim *SIM900) Disconnect() error {
	err := sim.port.Close()
	fmt.Println("Modem disconnected.")

	return err

}

// AT command (modem should return OK)
func (sim *SIM900) At() error {
	_, err := sim.SendCommand("")

	return err
}

// get last element in modem response (string array)
func GetLastLine(resultLines []string) string {
	length := len(resultLines)

	if length > 0 {
		return resultLines[length-1]
	}

	return ""
}

func ParseResponse(response []string, cmd string) []string {
	if strings.Index(cmd, "?") == len(cmd)-1 {
		cmd = strings.Replace(cmd, "?", ":", 1)
	}

	if strings.Index(cmd, "=%s") == len(cmd)-3 {
		cmd = strings.Replace(cmd, "=%s", ":", 1)
	}

	if strings.Index(cmd, ":") != len(cmd)-1 {
		cmd = cmd + ":"
	}

	// fmt.Println("cmd:", cmd)

	for _, line := range response {
		// fmt.Println("line:", line)

		if strings.Index(line, cmd) == 0 {
			cmdResponse := strings.Replace(line, cmd, "", 1)
			responses := strings.Split(cmdResponse, ",")

			for index := range responses {
				responses[index] = strings.TrimSpace(responses[index])
			}

			// fmt.Println("strs:", responses)

			return responses
		}
	}

	return nil
}

// Modem info
func (sim *SIM900) DisplayModemInfo() error {
	vendor, err := sim.SendCommand(CMD_VENDOR)
	if err != nil {
		return err
	}

	model, err := sim.SendCommand(CMD_MODEL)
	if err != nil {
		return err
	}

	revision, err := sim.SendCommand(CMD_REVISION)
	if err != nil {
		return err
	}

	serial, err := sim.SendCommand(CMD_SERIAL)
	if err != nil {
		return err
	}

	fmt.Println("Modem connected.")

	fmt.Println("Vendor:", GetLastLine(vendor))
	fmt.Println("Model:", GetLastLine(model))
	fmt.Println("Revision:", GetLastLine(revision))
	fmt.Println("Serial:", GetLastLine(serial))

	return err
}

// Disable Echo
func (sim *SIM900) DisableEcho() error {
	_, err := sim.SendCommand(CMD_DISABLE_ECHO)

	return err
}

// Enable Echo
func (sim *SIM900) EnableEcho() error {
	_, err := sim.SendCommand(CMD_ENABLE_ECHO)

	return err
}

// Get maximum supported baud rate
func (sim *SIM900) GetMaxBaudRate() (int, error) {
	data, err := sim.SendCommand(CMD_GET_SUPPORTED_RATES)
	if err != nil {
		return 0, err
	}

	values := ParseResponse(data, CMD_GET_RATES_RESPONSE)
	maxBaudRateStr := strings.TrimSuffix(values[len(values)-1], ")")

	fmt.Printf("variable maxBaudRateStr=%v is of type %T size %d\n", maxBaudRateStr, maxBaudRateStr, len(maxBaudRateStr))

	maxBaudRate, err := strconv.Atoi(maxBaudRateStr)
	if err != nil {
		return 0, err
	}

	if maxBaudRate < sim.config.Baud {
		return 0, errors.New("band rate lower then existing")
	}

	return maxBaudRate, nil
}

// Enable Echo
func (sim *SIM900) SetMaxBaudRate(baudRate int) error {
	cmd := fmt.Sprintf(CMD_SET_BAUD_RATES, baudRate)
	_, err := sim.SendCommand(cmd)

	return err
}

// send command and wait for OK or ERROR from modem, return array on strings
func (sim *SIM900) SendCommand(cmd string, args ...time.Duration) ([]string, error) {
	timeout := time.Second * 10

	if len(args) > 0 {
		timeout = args[0]
	}

	sim.logger.Printf("CMD >> AT%s", cmd)

	if _, err := sim.port.Write([]byte(CMD_AT + cmd + CMD_CR)); err != nil {
		return nil, err
	}

	return sim.AwaitCommand(CMD_FINISH_OK, false, timeout)
}

// wait for specific response from modem, return array on strings
func (sim *SIM900) AwaitCommand(okLine string, prefix bool, timeout time.Duration) ([]string, error) {
	timeExpired := false
	resultChan := make(chan []string, 1024)
	inputScanner := bufio.NewScanner(sim.port)

	// if CMD_DEBUG_LOG_ON
	if sim.debugMode {

	}

	go func() {
		result := make([]string, 0)

		for !timeExpired {
			for inputScanner.Scan() {
				inputLine := inputScanner.Text()

				if !prefix && inputLine == okLine {
					// if sim.debugMode {
					sim.logger.Printf("GET >> OK")
					// }
					resultChan <- result
					return
				}

				// in case of prefix = true (http response) don't modify strings
				if prefix && strings.HasPrefix(inputLine, okLine) {
					if sim.debugMode {
						sim.logger.Printf("FOUND PREFIX")
						sim.logger.Printf(inputLine)
					}
					result = append(result, inputLine)
					resultChan <- result
					return
				}

				if inputLine == CMD_FINISH_ERROR {
					// if sim.debugMode {
					sim.logger.Printf("GET >> ERROR")
					// }
					resultChan <- nil
					return
				}

				if sim.debugMode && strings.HasPrefix(inputLine, CMD_VERBOSE_ERROR) {
					sim.logger.Printf("FOUND ERROR")
					resultChan <- nil
					return
				}

				if len(inputLine) > 0 {
					result = append(result, inputLine)
					sim.logger.Printf("GET >> %s", inputLine)
				}
			}
		}
	}()

	select {
	case data := <-resultChan:
		if sim.debugMode {
			sim.logger.Printf("INF >> Command result: %s", data)
			for _, str := range data {
				fmt.Printf("* %s\n", str)
			}
		}

		if data == nil {
			return nil, errors.New("modem return ERROR")
		}

		return data, nil
	case <-time.After(timeout):
		timeExpired = true
		return nil, errors.New("waiting timeout expired")
	}
}

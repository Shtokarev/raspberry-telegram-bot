package sim900

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type SIM900 struct {
	port   *serial.Port
	config *serial.Config
	logger *log.Logger
	APN    string
}

type SMS struct {
	Id       int
	Status   string
	Operator string
	Datetime string
	Message  []byte
}

func New(name string, baud int) *SIM900 {
	return &SIM900{
		config: &serial.Config{
			Name:        name,
			Baud:        baud,
			ReadTimeout: 0,
			Size:        8,
		},
		logger: log.New(os.Stdout, "[SIM900] ", log.LstdFlags),
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

	//???
	err = sim.EnableEcho()

	return err
}

func (sim *SIM900) Disconnect() error {
	return sim.port.Close()
}

// AT command (modem should return OK)
func (sim *SIM900) At() error {
	_, err := sim.wait4response(CMD_AT, CMD_OK, time.Second*1)
	return err
}

// Disable Echo
func (sim *SIM900) DisableEcho() error {
	_, err := sim.wait4response(CMD_DISABLE_ECHO, CMD_OK, time.Second*1)
	return err
}

// Enable Echo
func (sim *SIM900) EnableEcho() error {
	_, err := sim.wait4response(CMD_ENABLE_ECHO, CMD_OK, time.Second*1)
	return err
}

func (sim *SIM900) wait4response(cmd string, expected string, timeout time.Duration) (string, error) {
	_, err := sim.port.Write([]byte(cmd + CMD_CR))

	if err != nil {
		return "", err
	}

	regexp := expected + "|" + CMD_ERROR
	response, err := sim.waitForRegexTimeout(regexp, timeout)

	if err != nil {
		return "", err
	}

	if strings.Contains(response, "ERROR") {
		return response, errors.New("errors found on command response")
	}

	return response, nil
}

func (sim *SIM900) waitForRegexTimeout(exp string, timeout time.Duration) (string, error) {
	timeExpired := false
	regex := regexp.MustCompile(`(?m)` + exp)
	found := make(chan string, 1)

	go func() {
		sim.logger.Printf("INF >> Waiting for RegExp: \"%s\"", exp)

		inputScanner := bufio.NewScanner(sim.port)

		for !timeExpired {
			for inputScanner.Scan() {
				line := inputScanner.Text()
				sim.logger.Printf("CMD >> %s", line)

				result := regex.FindAllString(line, -1)

				if len(result) > 0 {
					found <- result[0]
					return
				}
			}
		}
	}()

	select {
	case data := <-found:
		sim.logger.Printf("INF >> The RegExp: \"%s\" Has been matched: \"%s\"", exp, data)
		return data, nil
	case <-time.After(timeout):
		timeExpired = true
		sim.logger.Printf("INF >> Unable to match RegExp: \"%s\"", exp)
		return "", errors.New("waiting timeout expired")
	}
}

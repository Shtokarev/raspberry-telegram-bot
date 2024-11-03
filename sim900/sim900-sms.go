package sim900

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Send a SMS
func (sim *SIM900) SendSMS(number, msg string, chset string) error {
	// Set message format
	err := sim.SetSMSMode(TEXT_MODE)

	if err != nil {
		return err
	}

	cmd := fmt.Sprintf(CMD_CMGS, number)

	_, err = sim.port.Write([]byte(cmd + CMD_CR))

	if err != nil {
		return err
	}

	// Wait modem to be ready
	time.Sleep(time.Millisecond * 500)

	_, err = sim.SendCommand(msg + CMD_CTRL_Z)

	if err != nil {
		return err
	}

	return nil
}

// SetSMSMode selects SMS Message Format (0 = PDU mode, 1 = Text mode)
func (sim *SIM900) SetSMSMode(mode string) error {
	cmd := fmt.Sprintf(CMD_CMGF_SET, mode)
	_, err := sim.SendCommand(cmd)

	return err
}

// SMSMode reads SMS Message Format (0 = PDU mode, 1 = Text mode)
// func (sim *SIM900) SMSMode() (mode string, err error) {
// 	mode, err = sim.wait4response(CMD_CMGF, CMD_CMGF_REGEXP, time.Second*1)
// 	if err != nil {
// 		return
// 	}
// 	if len(mode) >= len(CMD_CMGF_RX) {
// 		mode = mode[len(CMD_CMGF_RX):]
// 	}
// 	return
// }

// GetSMSList retrieves unreaded SMS list from inbox
func (sim *SIM900) GetSMSList(status string) (list []string, err error) {
	err = sim.SetSMSMode(TEXT_MODE)
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf(CMD_CMGL, status)
	list, err = sim.SendCommand(cmd)

	if err != nil {
		return nil, err
	}

	return list, err
}

// ReadSMS retrieves SMS text from inbox memory by ID
func (sim *SIM900) ReadSMS(id string) (msg string, err error) {
	// Set message format
	err = sim.SetSMSMode(PDU_MODE)
	// err = sim.SetSMSMode(TEXT_MODE)
	if err != nil {
		return "", err
	}

	inputScanner := bufio.NewScanner(sim.port)
	regexOk := regexp.MustCompile(CMD_OK)

	// Send command
	cmd := fmt.Sprintf(CMD_AT+CMD_CMGR, id)

	_, err = sim.port.Write([]byte(cmd + CMD_CR))

	if err != nil {
		return "", err
	}

	// if _, err := sim.wait4response(cmd, CMD_CMGR_REGEXP, time.Second*5); err != nil {
	// 	return "", err
	// }

	// time.Sleep(time.Millisecond * 1000)
	// Reading succesful get message data
	// return s.port.ReadLine()
	// inputScanner := bufio.NewScanner(sim.port)
	var sb strings.Builder
	modemOutput := make([]string, 0)

	scan := inputScanner.Scan()
	fmt.Println("Scan:", scan)

	for inputScanner.Scan() {
		// line := inputScanner.Bytes()
		line := inputScanner.Text()
		// sim.logger.Printf("CMD >> %s", line)
		//  .Text()
		// sim.logger.Printf("CMD >> %s", line)
		fmt.Println("line:", line)
		modemOutput = append(modemOutput, line)

		sb.WriteString(string(line))

		// if (regexOk === )
		result := regexOk.FindAllString(line, -1)
		fmt.Println("result:", result)
		if len(result) > 0 {

			fmt.Println("modemOutput:", modemOutput)
			r := ParseResponse(modemOutput, CMD_CMGR)
			fmt.Println("r===>:", r)

			return sb.String(), err
		}

	}

	return sb.String(), err
}

// ReadSMS deletes SMS from inbox memory by ID
func (sim *SIM900) DeleteSMS(id string) error {
	cmd := fmt.Sprintf(CMD_CMGD, id)
	_, err := sim.SendCommand(cmd)

	return err
}

package sim900

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (sim *SIM900) GET(url string) ([]byte, error) {
	if err := sim.CheckGprsReady(); err != nil {
		return nil, err
	}

	if err := sim.InitGprsSession(); err != nil {
		return nil, err
	}

	if err := sim.InitHTTP(); err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf(CMD_HTTP_PARAMETERS, HTTP_CID)
	if _, err := sim.SendCommand(cmd); err != nil {
		return nil, err
	}

	cmd = fmt.Sprintf(CMD_HTTP_PARAMETERS, fmt.Sprintf(HTTP_URL, url))
	if _, err := sim.SendCommand(cmd); err != nil {
		return nil, err
	}

	// FIXME: check content-type?
	status, size, err := sim.HTTPRequest(0)

	if err != nil {
		sim.logger.Println("HTTP REQUEST ERROR:", err)
		return nil, err
	}

	sim.logger.Printf("File obtained with status code:%d, size:%d bytes\n", status, size)

	if size > 0 {
		httpData, err := sim.ReadHTTPResponse()
		fmt.Printf("httpData!!!  %#v\n", httpData)

		if err != nil {
			return nil, err
		}

		return httpData, nil
	}

	defer sim.CloseGprsSession()

	return nil, nil
}

func (sim *SIM900) CheckGprsReady() error {
	if err := sim.CheckSimCard(); err != nil {
		return err
	}

	if err := sim.CheckSignalLevel(); err != nil {
		return err
	}

	if err := sim.CheckRegistered(); err != nil {
		return err
	}

	if err := sim.CheckGPRSEnabled(); err != nil {
		return err
	}

	return nil
}

func (sim *SIM900) InitGprsSession() error {
	// set GPRS context
	cmd := fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_TYPE)
	if _, err := sim.SendCommand(cmd); err != nil {
		sim.logger.Println("SAPBR ERROR:", err)
		return err
	}

	// Set the Access Point Name (APN) for the network provider
	cmd = fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_APN)
	if _, err := sim.SendCommand(cmd); err != nil {
		sim.logger.Println("SAPBR ERROR:", err)
		return err
	}

	// Activate GPRS context
	if err := sim.ActivateGPRSContext(); err != nil {
		return err
	}

	// Open GPRS connection bearer
	ip, err := sim.OpenGPRSContextBearer()
	if err != nil {
		return err
	}

	sim.logger.Printf("GRPC connected, obtained IP: %s\n", ip)

	return nil
}

func (sim *SIM900) CloseGprsSession() error {
	if err := sim.TerminateHTTP(); err != nil {
		return err
	}

	if err := sim.DeactivateGPRSContext(); err != nil {
		return err
	}

	return nil
}

// func (sim *SIM900) InitGprs() error {
// 	if _, err := sim.SendCommand(CMD_DEBUG_LOG_ON); err != nil {
// 		return err
// 	}

// 	if err := sim.CheckSimCard(); err != nil {
// 		return err
// 	}

// 	if err := sim.CheckSignalLevel(); err != nil {
// 		return err
// 	}

// 	if err := sim.CheckRegistered(); err != nil {
// 		return err
// 	}

// 	if err := sim.CheckGPRSEnabled(); err != nil {
// 		return err
// 	}

// 	// set GPRS context
// 	cmd := fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_TYPE)
// 	if _, err := sim.SendCommand(cmd); err != nil {
// 		sim.logger.Println("SAPBR ERROR:", err)
// 		return err
// 	}

// 	// Set the Access Point Name (APN) for the network provider
// 	cmd = fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_APN)
// 	if _, err := sim.SendCommand(cmd); err != nil {
// 		sim.logger.Println("SAPBR ERROR:", err)
// 		return err
// 	}

// 	// Activate GPRS context
// 	if err := sim.ActivateGPRSContext(); err != nil {
// 		return err
// 	}

// 	// Open GPRS connection bearer
// 	ip, err := sim.OpenGPRSContextBearer()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("GRPC connected, obtained IP: %s\n", ip)

// 	if err = sim.InitHTTP(); err != nil {
// 		return err
// 	}

// 	cmd = fmt.Sprintf(CMD_HTTP_PARAMETERS, HTTP_CID)
// 	if _, err := sim.SendCommand(cmd); err != nil {
// 		return err
// 	}

// 	cmd = fmt.Sprintf(CMD_HTTP_PARAMETERS, HTTP_URL)
// 	if _, err = sim.SendCommand(cmd); err != nil {
// 		return err
// 	}

// 	status, size, err := sim.HTTPRequest(0)
// 	if err != nil {
// 		sim.logger.Println("HTTP REQUEST ERROR:", err)
// 		return err
// 	}

// 	fmt.Printf("size:  %d\n", size)
// 	fmt.Printf("status:  %d\n", status)

// 	if size > 0 {
// 		httpData, err := sim.ReadHTTPResponse()
// 		fmt.Printf("httpData!!!  %#v\n", httpData)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if err = sim.TerminateHTTP(); err != nil {
// 		return err
// 	}

// 	if err := sim.DeactivateGPRSContext(); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (sim *SIM900) CheckSimCard() error {
	simStatus, err := sim.SendCommand(CMD_CHECK_SIM)
	if err != nil {
		return err
	}

	values := ParseResponse(simStatus, CMD_CHECK_SIM)
	if values[0] != CMD_FINISH_READY {
		return errors.New("SIM card is not ready")
	}

	return err
}

func (sim *SIM900) CheckSignalLevel() error {
	signalLevel, err := sim.SendCommand(CMD_SIGNAL_LVL)
	if err != nil {
		return err
	}

	signal := ParseResponse(signalLevel, CMD_SIGNAL_LVL)

	values, err := strconv.Atoi(signal[0])
	if err != nil {
		return err
	}

	if values < 5 {
		msg := fmt.Sprintf("signal is %s - too low to connect", signal[0])
		return errors.New(msg)
	}

	return err
}

func (sim *SIM900) CheckRegistered() error {
	registered, err := sim.SendCommand(CMD_REGISTERED)
	if err != nil {
		return err
	}

	values := ParseResponse(registered, CMD_REGISTERED)

	if values[1] != "1" {
		return errors.New("modem is not registered in GSM")
	}

	return err
}

func (sim *SIM900) CheckGPRSEnabled() error {
	enabled, err := sim.SendCommand(CMD_GPRS_ENABLED)
	if err != nil {
		return err
	}

	values := ParseResponse(enabled, CMD_GPRS_ENABLED)

	if values[0] != "1" {
		return errors.New("GPRS is not allowed")
	}

	return err
}

// Open GPRS bearer
func (sim *SIM900) OpenGPRSContextBearer() (string, error) {
	cmd := fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_OPEN)

	enabled, err := sim.SendCommand(cmd)
	if err != nil {
		return "", err
	}

	values := ParseResponse(enabled, CMD_GPRS_CONTEXT_SET)

	return values[len(values)-1], err
}

// Activate GPRS context
func (sim *SIM900) ActivateGPRSContext() error {
	cmd := fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_ACTIVATE)

	if _, err := sim.SendCommand(cmd); err != nil {
		// try to deactivate and activate again
		sim.DeactivateGPRSContext()

		if _, err = sim.SendCommand(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Deactivate GPRS context
func (sim *SIM900) DeactivateGPRSContext() error {
	cmd := fmt.Sprintf(CMD_GPRS_CONTEXT_SET, CONTEXT_CLOSE)

	if _, err := sim.SendCommand(cmd); err != nil {
		sim.logger.Println("DEACTIVATE CONTEXT ERROR:", err)
		return err
	}

	return nil
}

// Initialize HTTP
func (sim *SIM900) InitHTTP() error {
	if _, err := sim.SendCommand(CMD_HTTP_INIT); err != nil {
		// try to terminate and initialize again
		sim.TerminateHTTP()

		if _, err := sim.SendCommand(CMD_HTTP_INIT); err != nil {
			return err
		}
	}

	return nil
}

// Terminate HTTP
func (sim *SIM900) TerminateHTTP() error {
	if _, err := sim.SendCommand(CMD_HTTP_TERM); err != nil {
		sim.logger.Println("TERMINATE HTTP ERROR:", err)
		return err
	}

	return nil
}

// Make HTTP Request
func (sim *SIM900) HTTPRequest(method int) (int, int, error) {
	if method != 0 && method != 1 {
		return 0, 0, errors.New("only GET (0) or POST (1) method are allowed")
	}

	cmd := fmt.Sprintf(CMD_HTTP_ACTION, method)
	if _, err := sim.port.Write([]byte(CMD_AT + cmd + CMD_CR)); err != nil {
		return 0, 0, err
	}

	data, err := sim.AwaitCommand(CMD_HTTP_ACTION_RESPONSE, true, time.Second*15)
	if err != nil {
		return 0, 0, err
	}

	values := ParseResponse(data, CMD_HTTP_ACTION_RESPONSE)

	status, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, err
	}

	size, err := strconv.Atoi(values[2])
	if err != nil {
		return 0, 0, err
	}

	return status, size, nil
}

// Read HTTP Response
func (sim *SIM900) ReadHTTPResponse() ([]byte, error) {
	timeout := time.Second * 30

	result, err := sim.SendCommand(CMD_HTTP_READ, timeout)
	if err != nil {
		return nil, err
	}

	data := strings.Join(result, CMD_CR)

	return []byte(data), nil
}

package sim900

// commands
const (
	CMD_AT           string = "AT"
	CMD_CMGF         string = "AT+CMGF?"
	CMD_CMGF_SET     string = "AT+CMGF=%s"
	CMD_CMGF_RX      string = "+CMGF: "
	CMD_CTRL_Z       string = "\x1A"
	CMD_CMGS         string = "AT+CMGS=\"%s\""
	CMD_CSCS         string = "AT+CSCS=\"%s\""
	CMD_CMGD         string = "AT+CMGD=%s"
	CMD_CMGR         string = "AT+CMGR=%s"
	CMD_CMGR_RX      string = "+CMGR: "
	CMD_CMTI_RX      string = "+CMTI: \"SM\","
	CMD_CMGL         string = "AT+CMGL=\"%s\""
	CMD_CR           string = "\r\n"
	CMD_FINISH       string = "\r\nOK"
	CMD_DISABLE_ECHO string = "ATE0"
	CMD_ENABLE_ECHO  string = "ATE1"
)

const (
	SMS_ALL           string = "ALL"
	SMS_READED        string = "REC READ"
	SMS_UNREADED      string = "REC UNREAD"
	SMS_SENT          string = "STO SENT"
	SMS_UNSENT        string = "STO UNSENT"
	SMS_STATUS_UNREAD int    = 0
	SMS_STATUS_READ   int    = 1
)

// responses
const (
	CMD_OK             string = "(^OK$)"
	CMD_ERROR          string = "(^ERROR$)"
	CMD_CMGF_REGEXP    string = "(^[+]CMGF[:] [0-9]+$)"
	CMD_CMGS_RX_REGEXP string = "(^[+]CMGS[:] [0-9]+$)"
	CMD_CMGR_REGEXP    string = "(^[+]CMGR[:] .*)"
	CMD_CMTI_REGEXP    string = "(^[+]CMTI[:] \"SM\",[0-9]+$)"
)

// SMS Message Format
const (
	PDU_MODE  string = "0"
	TEXT_MODE string = "1"
)

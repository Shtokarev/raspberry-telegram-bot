package sim900

// commands
const (
	CMD_AT                   string = "AT"
	CMD_CHECK_SIM            string = "+CPIN?"
	CMD_CMGF                 string = "AT+CMGF?"
	CMD_CMGF_SET             string = "+CMGF=%s"
	CMD_CMGF_RX              string = "+CMGF: "
	CMD_CTRL_Z               string = "\x1A"
	CMD_CMGS                 string = "AT+CMGS=\"%s\""
	CMD_CSCS                 string = "AT+CSCS=\"%s\""
	CMD_CMGD                 string = "AT+CMGD=%s"
	CMD_CMGR                 string = "+CMGR=%s"
	CMD_CMGR_RX              string = "+CMGR: "
	CMD_CMTI_RX              string = "+CMTI: \"SM\","
	CMD_CMGL                 string = "AT+CMGL=\"%s\""
	CMD_CR                   string = "\r\n"
	CMD_DISABLE_ECHO         string = "E0"
	CMD_ENABLE_ECHO          string = "E1"
	CMD_VENDOR               string = "+GMI"
	CMD_MODEL                string = "+GMM"
	CMD_REVISION             string = "+GMR"
	CMD_SERIAL               string = "+GSN"
	CMD_SIGNAL_LVL           string = "+CSQ"
	CMD_REGISTERED           string = "+CREG?"
	CMD_GPRS_ENABLED         string = "+CGATT?"
	CMD_GPRS_CONTEXT_SET     string = "+SAPBR=%s"
	CMD_HTTP_INIT            string = "+HTTPINIT"
	CMD_HTTP_PARAMETERS      string = "+HTTPPARA=%s"
	CMD_HTTP_SSL             string = "+HTTPSSL=1"
	CMD_HTTP_ACTION          string = "+HTTPACTION=%d"
	CMD_HTTP_ACTION2         string = "+HTTPACTION=0"
	CMD_HTTP_ACTION_RESPONSE string = "+HTTPACTION:"
	CMD_HTTP_READ            string = "+HTTPREAD"
	CMD_HTTP_TERM            string = "+HTTPTERM"
	CMD_TURN_OFF_CONNECTION  string = "+CIPSHUT"
	CMD_DEBUG_LOG_ON         string = "+CMEE=2"
	CMD_DEBUG_LOG_OFF        string = "+CMEE=0"
	CMD_SET_BAUD_RATES       string = "+IPR=%d"
	CMD_GET_SUPPORTED_RATES  string = "+IPR=?"
	CMD_GET_RATES_RESPONSE   string = "+IPR:"
)

const (
	CONTEXT_CLOSE    string = "0,1"
	CONTEXT_TYPE     string = "3,1,\"Contype\",\"GPRS\""
	CONTEXT_APN      string = "3,1,\"APN\",\"internet.mts.ru\""
	CONTEXT_ACTIVATE string = "1,1"
	CONTEXT_OPEN     string = "2,1"
)

const (
	HTTP_CID string = "\"CID\",1"
	HTTP_URL string = "\"URL\",\"%s\""
)

const (
	CMD_FINISH_OK     string = "OK"
	CMD_FINISH_ERROR  string = "ERROR"
	CMD_VERBOSE_ERROR string = "+CME ERROR:"
	CMD_FINISH_READY  string = "READY"
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

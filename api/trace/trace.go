package trace

import (
	"fmt"
	"log"
)

// Info Traces
const (
	AppStarted                   = "APP_STARTED"
	RedirectedSuccessfully       = "REDIRECTED_SUCCESSFULLY"
	SuccessfullyConnectedToRedis = "SUCCESSFULLY_CONNECTED_TO_REDIS"
	HttpEnforced                 = "HTTP_ENFORCED"
)

// Errors
const (
	ErrorLoadingEnvFiles          = "ERROR_LOADING_ENV_FILES"
	ErrorParsingBodyToJSON        = "ERROR_PARSING_BODY_TO_JSON"
	ErrorKeyNotFoundInRedis       = "ERROR_KEY_NOT_FOUND_IN_REDIS"
	ErrorConnectingToRedis        = "ERROR_CONNECTING_TO_REDIS"
	ErrorInvalidUrl               = "ERROR_INVALID_URL"
	ErrorSelfDomainLoopPrevented  = "ERROR_SELF_DOMAIN_LOOP_PREVENTED"
	ErrorCustomShortAlreadyExists = "ERROR_CUSTOM_SHORT_ALREADY_EXISTS"
)

// Fatal Errors
const (
	ErrorInStartingApp = "ERROR_IN_STARTING_APP"
)

func LogInfo(traceCode string, data map[string]interface{}) {
	if data == nil {
		fmt.Println(traceCode)
	}
	fmt.Println(traceCode, "data - ", data)
}

func LogError(traceCode string, err error, data map[string]interface{}) {
	if data == nil {
		fmt.Println(traceCode, "error - ", err)
	}
	fmt.Println(traceCode, "error - ", err, "data - ", data)
}

func LogFatalError(traceCode string, err error, data map[string]interface{}) {
	if data == nil {
		log.Fatal(traceCode, err)
	}
	log.Fatal(traceCode, "error - ", err, "data - ", data)
}

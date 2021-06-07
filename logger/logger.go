package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "[info] ", log.LstdFlags|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[warning] ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "[error] ", log.LstdFlags|log.Lshortfile)
}

func TenantInfoLog(tenantId string) *log.Logger {
	str := fmt.Sprintf("[info] tenant=%s ", tenantId)
	return log.New(os.Stderr, str, log.LstdFlags|log.Lshortfile)
}

func TenantWarningLog(tenantId string) *log.Logger {
	str := fmt.Sprintf("[warning] tenant=%s ", tenantId)
	return log.New(os.Stderr, str, log.LstdFlags|log.Lshortfile)
}
func TenantErrorLog(tenantId string) *log.Logger {
	str := fmt.Sprintf("[error] tenant=%s ", tenantId)
	return log.New(os.Stderr, str, log.LstdFlags|log.Lshortfile)
}

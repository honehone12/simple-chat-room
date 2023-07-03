package common

import "fmt"

func SETUPMsg(msg string) string {
	return fmt.Sprintf("[SETUP] %s", msg)
}

func INFOMsg(msg string) string {
	return fmt.Sprintf("[INFO] %s", msg)
}

func PlayerMsg(name string, msg string) string {
	return fmt.Sprintf("(%s) %s", name, msg)
}

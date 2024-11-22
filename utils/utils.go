package utils

import "fmt"

func PrintLine(format string, arg ...interface{}) {
	fmt.Println(fmt.Sprintf(format, arg...))
}

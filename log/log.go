package log

import "fmt"

func Deubugf(format string, v ...interface{}){
	fmt.Printf(format, v...)
}

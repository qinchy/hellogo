package scheduler

import (
	"fmt"
	"time"
)

func PrintTimeEveryMinute() {
	for true {
		fmt.Println(time.Now())
		time.Sleep(1 * time.Second)
	}
}

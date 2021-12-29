package main

import (
	"testing"
	"fmt"
)
//1638334800
func TestDateToDiscordTime(t *testing.T) {
    time, err := DateToDiscordTime("2021-12-01")
    if err !=  nil{
       fmt.Println(time)
    }
}
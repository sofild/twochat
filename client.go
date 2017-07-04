package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	//"log"
	"bufio"
	"os"
	"time"
)

var origin = "http://127.0.0.1:8080/"
var url = "ws://127.0.0.1:8080/echo"
var username string

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		fmt.Println(17, err.Error())
	}
	fmt.Print("请输入您的英文名：")
	reader1 := bufio.NewReader(os.Stdin)
	sms1, _, err := reader1.ReadLine()
	if err != nil {
		fmt.Println(24, err.Error())
	}
	username = string(sms1)
	for {
		fmt.Print("请输入您要发送的消息：")
		reader := bufio.NewReader(os.Stdin)
		sms, _, err1 := reader.ReadLine()
		if err1 != nil {
			fmt.Println(24, err1.Error())
		}
        ch := make(chan string, 1)
        ch <- string(sms)
        go pullMsg(ch, username, ws)
    }
	ws.Close() //关闭连接
}

func pullMsg(sms chan string, username string, ws *websocket.Conn) {
	tm := time.NewTimer(time.Second * 15)
    var data string
    for {
		select {
		case data = <- sms:
            msgHandle(ws, username, data)
		    break
        case <-tm.C:
            tm.Reset(time.Second * 15)
            msgHandle(ws, username, "pull message")
		}
	}
}

func msgHandle(ws *websocket.Conn, username string, sms string) {
	sms_str := fmt.Sprintf("@%s:%s", username, sms)
	_, err2 := ws.Write([]byte(sms_str))
	if err2 != nil {
		fmt.Println(28, err2.Error())
	}
	msg := make([]byte, 10240)
	m, err3 := ws.Read(msg)
	if err3 != nil {
		fmt.Println(33, err3.Error())
	}
    receivedMsg := string(msg[0:m])
    if receivedMsg != "no message" {
	    fmt.Printf("%s\n", receivedMsg)
    }
}

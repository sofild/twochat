package main

import (
	"fmt"
	"websocket"
	//"log"
	"net/http"
	//"strconv"
    "strings"
)

type message struct{
    username string
    content string
}

var msgData []*message

func get(username string) string {
   var message string = ""
   for i, data := range msgData {
        if data.username != username {
            message += ""+data.username+":"+data.content+"\n"
            msgData[i].content = ""
        }
   }
   formatMsg()
   if message==""{
        message = "no message"
   }
   return message
}

func put(username string, content string){
    if content=="pull message"{
        return
    }
    message := new(message)
    message.username = username
    message.content = content
    msgData = append(msgData, message)
}

func formatMsg(){
    var datas []*message
    for _, data := range msgData {
        if data.content!=""{
            datas = append(datas, data)
        }
    }
    msgData = datas
}

func echoHandler(ws *websocket.Conn) {
	for {
		msg := make([]byte, 102400)
		n, err := ws.Read(msg)
		checkErr(err, 19)
        fmt.Printf("Receive: %s\n", msg[:n])
        username,content := getMsg(msg[:n])
        put(username,content)
        data := get(username)
	    _, err = ws.Write([]byte(data))
        checkErr(err, 26)
        if err != nil {
            break
        }
        fmt.Printf("Send: %s\n", data)
    }
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func checkErr(err error, line int){
    if err!=nil{
        fmt.Println(line, err.Error())
    }
}

func getMsg(msg []byte) (string, string) {
    msg_str := string(msg)
    if strings.HasPrefix(msg_str,"@") {
        msg_arr := strings.Split(msg_str, ":")
        username := strings.Replace(msg_arr[0], "@", "", 1)
        var content string = ""
        for _,key := range msg_arr[1:] {
            content += ""+key+""
        }
        return username, content
    }else{
        return "guest","^-^"
    }
}

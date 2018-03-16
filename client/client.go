package main

import(
    "fmt"
    "net"
    "os"
    "time"
    "encoding/json"
)

type Client struct{
    conn *net.UDPConn
    gkey bool   //用来判断用户退出
    userID int
    userName string
    sendMessages chan string
    receiveMessages chan string
}

type Message struct{
	Status int
	UserID int
	UserName string
	Content string
}


func (c *Client) func_sendMessage(sid int,msg string){

	m:= Message{
		Status:2,
		UserID:c.userID,
		UserName: c.userName,
		Content: msg,
	}
	str, err := json.Marshal(m)

	if err != nil {
		fmt.Println("json err:", err)
	}

	_, err = c.conn.Write([]byte(string(str)))
    checkError(err,"func_sendMessage")
}

func (c *Client) sendMessage() {
    for c.gkey {
        msg := <- c.sendMessages
		m:= Message{
			Status:1,
			UserID:c.userID,
			UserName: c.userName,
			Content: msg,
		}
		str, err := json.Marshal(m)

		if err != nil {
			fmt.Println("json err:", err)
		}
        _,err = c.conn.Write([]byte(string(str)))
        checkError(err,"sendMessage")
    }

}

//接收
func (c *Client) receiveMessage() {
    var buf [512]byte
    for c.gkey {
        n,err := c.conn.Read(buf[0:])
        checkError(err, "receiveMessage")
        c.receiveMessages <- string(buf[0:n])
    }
    
}
//获得输入并处理之，这里有Println
func (c *Client) getMessage() {
    var msg string
    for c.gkey {
        _,err := fmt.Scanln(&msg)
        checkError(err, "getMessage")
        if msg == ":quit" {
            c.gkey = false
        }else{
            c.sendMessages <- msg
        }
    }
}
//打印，这里有Println
func (c *Client) printMessage() {
    //var msg string
    for c.gkey {
        msg := <- c.receiveMessages
        var m Message
        json.Unmarshal([]byte(msg),&m)
        fmt.Println(m.UserName,":",m.Content)
    }
}

func nowTime() string {
    return time.Now().String()
}

func checkError(err error, funcName string){
    if err != nil{
        fmt.Fprintf(os.Stderr,"Fatal error:%s-----in func:%s",err.Error(), funcName)
        os.Exit(1)
    }
}

func main(){
    if len(os.Args) != 2{
        fmt.Fprintf(os.Stderr, "Usage:%s host:port", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]
    udpAddr, err := net.ResolveUDPAddr("udp4",service)
    checkError(err,"main")

    var c Client
    c.gkey = true
    c.sendMessages = make(chan string)
    c.receiveMessages = make(chan string)

    fmt.Print("用户id: ")
    _,err = fmt.Scanln(&c.userID)
    checkError(err,"main")
    fmt.Print("用户昵称: ")
    _,err = fmt.Scanln(&c.userName)
    checkError(err,"main")

    c.conn,err = net.DialUDP("udp",nil,udpAddr)
    checkError(err,"main")
    //fmt.Println(c)
    defer c.conn.Close()


    //发送进入聊天室消息,类型1，###1##uid##uName##进入聊天室###
    //messagestr := fmt.Sprintf("###1##%d##%s###", c.userID, c.userName)
    //_,err = c.conn.Write([]byte(messagestr))
    //checkError(err)
    c.func_sendMessage(1,c.userName + "进入聊天室")

    //go c.getMessage()
    go c.printMessage()
    go c.receiveMessage()

    go c.sendMessage()
    c.getMessage()

    c.func_sendMessage(3,c.userName + "离开聊天室")
    fmt.Println("退出成功!")


    os.Exit(0)
}

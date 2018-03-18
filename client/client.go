package main

import(
    "fmt"
    "net"
    "os"
    "time"
    "encoding/json"
    "../core"
	"log"
)

type Client struct{
    conn *net.UDPConn
    gkey bool   //用来判断用户退出
    userID int
    userName string
    sendMessages chan string
    receiveMessages chan string
    chiper core.Cipher
}

type Message struct{
	Status int
	UserID int
	UserName string
	Content string
}


func (c *Client) func_sendMessage(sid int,msg string){

	m:= core.Message{
		Status:2,
		UserID:c.userID,
		UserName: c.userName,
		Content: msg,
	}
	str, err := json.Marshal(m)

	if err != nil {
		fmt.Println("json err:", err)
	}

	str= []byte(string(str))
	str, err = c.chiper.EncryptMessage(str)

	_, err = c.conn.Write(str)
    checkError(err,"func_sendMessage")
}

func (c *Client) sendMessage() {
    for c.gkey {
        msg := <- c.sendMessages
		m:= core.Message{
			Status:1,
			UserID:c.userID,
			UserName: c.userName,
			Content: msg,
		}
		str, err := json.Marshal(m)

		if err != nil {
			fmt.Println("json err:", err)
		}

		str= []byte(string(str))
		str, err = c.chiper.EncryptMessage(str)

		_,err = c.conn.Write(str)
        checkError(err,"sendMessage")
    }

}

func (c *Client) receiveMessage() {
    var buf [512]byte
    for c.gkey {
        n,err := c.conn.Read(buf[0:])
        checkError(err, "receiveMessage")
        c.receiveMessages <- string(buf[0:n])
    }
    
}

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

func (c *Client) printMessage() {
    //var msg string
    for c.gkey {
        msg := <- c.receiveMessages
		dmsg,_:= c.chiper.DecryptMessage([]byte(msg))
        var m core.Message
        json.Unmarshal(dmsg,&m)
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
	log.SetFlags(log.Lshortfile)
	config := &core.Config{}
	config.ReadConfig()
	key, err := core.ParsePassword(config.Key)

	log.Println("使用配置：", fmt.Sprintf(`
远程服务地址 remote：%s
远程服务端口 port %s
密码 password：%s
	`, config.RemoteAddr,config.ListenAddr, key))

    service := config.RemoteAddr + config.ListenAddr
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err,"main")

    var c Client
    c.gkey = true
    c.chiper.Key = key[:]
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
    defer c.conn.Close()

    c.func_sendMessage(1,c.userName + "进入聊天室")

    go c.printMessage()
    go c.receiveMessage()

    go c.sendMessage()
    c.getMessage()

    c.func_sendMessage(3,c.userName + "离开聊天室")
    fmt.Println("退出成功!")


    os.Exit(0)
}

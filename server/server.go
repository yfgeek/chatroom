package main

import(
    "fmt"
    "net"
    "os"
    "encoding/json"
    "../core"
)

const(
    PORT string = ":1200"
    KEY string = "chatroom12345678"
)

type Server struct{
    conn     *net.UDPConn
    messages chan string
    clients  map [int]Client
    cipher   core.Cipher
}

type Client struct{
    userID int
    userName string
    userAddr *net.UDPAddr

}

func (s *Server) handleMessage(){
    var buf [512]byte
    n, addr, err := s.conn.ReadFromUDP(buf[0:])
    if err != nil{
        return
    }
    msg := buf[0:n]
    //分析消息
    fmt.Println("收到数据包",msg)
    m := s.analyzeMessage(msg)
    switch m.Status{
        //进入聊天室消息
        case 1:
            var c Client
            c.userAddr = addr
            c.userID = m.UserID
            c.userName = m.UserName
            s.clients[c.userID] = c //添加用户
            s.messages <- string(msg)
            fmt.Println("3")
        //用户发送消息
        case 2:
            s.messages <- string(msg)
        //client发来的退出消息
        case 3:
            delete(s.clients, m.UserID)
            s.messages <- string(msg)
        default:
            fmt.Println("未识别消息", string(msg))
    }

}

func (s *Server) analyzeMessage(msg []byte) (m core.Message) {
    msg,_ = s.cipher.DecryptMessage(msg)
    json.Unmarshal(msg, &m)
    return
}

func (s *Server) sendMessage() {
    for{
        msg := <- s.messages
        //daytime := time.Now().String()
        sendstr := msg
        for _,c := range s.clients {
            fmt.Println(c)
            fmt.Println("分发数据包",sendstr,c.userAddr)
            n,err := s.conn.WriteToUDP([]byte(sendstr),c.userAddr)
            fmt.Println(n,err)
        }
    }

}

func checkError(err error){
    if err != nil{
        fmt.Fprintf(os.Stderr,"Fatal error:%s",err.Error())
        os.Exit(1)
    }
}

func main(){
    udpAddr, err := net.ResolveUDPAddr("udp4",PORT)
    checkError(err)

    var s Server
    s.messages = make(chan string,20)
    s.clients =make(map[int]Client,0)

    s.conn,err = net.ListenUDP("udp",udpAddr)
    s.cipher.Key = []byte(KEY)
    checkError(err)

    go s.sendMessage()

    for{
        s.handleMessage()
    }
}

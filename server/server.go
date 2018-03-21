package main

import(
    "fmt"
    "net"
    "os"
    "encoding/json"
    "../core"
    "github.com/phayes/freeport"
	"log"
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

const(
	NEW_USER = iota + 1
	NEW_MESSAGE
	DELETE_USER
)

var userInitialID = 0

func (s *Server) handleMessage(){
    var buf [512]byte
    n, addr, err := s.conn.ReadFromUDP(buf[0:])
    if err != nil{
        return
    }
    msg := buf[0:n]
    fmt.Println("Got the packet",msg)
    m := s.analyzeMessage(msg)
    switch m.Status{
        case NEW_USER:
            var c Client
            c.userAddr = addr
            c.userID = userInitialID
			userInitialID++
            c.userName = m.UserName
            s.clients[c.userID] = c
            s.messages <- string(msg)
        case NEW_MESSAGE:
            s.messages <- string(msg)
        case DELETE_USER:
            delete(s.clients, m.UserID)
            s.messages <- string(msg)
        default:
            fmt.Println("Cannot read the message:", string(msg))
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
            fmt.Println("Transfer the packet",c.userAddr)
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
	log.SetFlags(log.Llongfile)
	// Generate a random server port
	port, err := freeport.GetFreePort()
	if err != nil {
		// 随机端口失败就采用 7448
		port = 1200
	}
	// Default config
	k:=core.RandPassword()
	config := &core.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Key: k.String(),
	}
	config.ReadConfig()
	config.SaveConfig()

	// Parse the config
	key, err := core.ParsePassword(config.Key)

	log.Println("Configruation", fmt.Sprintf(`
Server Port：
%s
Key：
%s
	`, config.ListenAddr, config.Key))

	if err != nil {
		log.Fatalln(err)
	}
	udpAddr, err := net.ResolveUDPAddr("udp4",config.ListenAddr)

	if err != nil {
		log.Fatalln(err)
	}

    var s Server
    s.messages = make(chan string,20)
    s.clients =make(map[int]Client,0)

    s.conn,err = net.ListenUDP("udp",udpAddr)
    s.cipher.Key = key[:]
    checkError(err)

    go s.sendMessage()

    for{
        s.handleMessage()
    }
}

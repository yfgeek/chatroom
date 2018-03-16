package core

type Message struct{
	Status int
	UserID int
	UserName string
	Content string
}



//
//func main(){
//	var m Message
//	chiper :=m.NewChiper()
//	t:= []byte("这是一个秘密")
//	text,_ := m.EncryptMessage(chiper,t)
//	fmt.Println(text)
//	d,_ := m.DecryptMessage(chiper,text)
//
//	fmt.Println(string(d))
//
//
//}
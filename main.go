package main

import (
	"bufio"
	"log"
	"net"
	"time"
)

var (
	usersMap     = make(map[string]User, 0)
	publicChan   = make(chan UserMsg, 100) // 存放公共消息
	filterString = []string{"", "cao"}     // 存放需要过滤的字符串
)

func main() {
	listenter, err := net.Listen("tcp", ":9091")
	if err != nil {
		panic(err)
	}

	// 启动广播系统
	go Broadcast()

	// 为每个连接的用户启动一个协程
	for {
		conn, err := listenter.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go ConnHandler(conn)
	}
}

func ConnHandler(conn net.Conn) {
	defer conn.Close()

	// 1. Create a user object
	userId := GetUserID(conn)
	user := User{
		UserID:      userId,
		Addr:        conn.RemoteAddr().String(),
		CreateTime:  time.Now(),
		MessageChan: make(chan string, 100),
	}

	// 2. Start user message handler
	go SendMessage(conn, user)

	// 3. Register user info to user state list just like `usersMap`
	// Send user login message to all online user at first
	saveUserOnline(userId, user)
	FirstLoginMsg(user)

	// 4. Watch user input
	reader := bufio.NewReader(conn)

	msg := UserMsg{
		UserID: user.UserID,
		Msg:    "",
	}

	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			continue
		}
		strData := string(data)
		// filter some string from user input
		if HasStrContain(strData, filterString) {
			continue
		} else if HasStrContain(strData, allowActions) {
			actionOpt(strData, conn, user)
			continue
		}

		msg.Msg = strData + "\n"

		// write message to publice chan, so as all online receiver these massage, except current user
		SendMessageToPublic(user, msg)
		// publicChan <- msg

	}

}

// 公共广播系统，1）用户首次登陆和退出通知、2）用户消息发送
func Broadcast() {
	select {
	case <-publicChan:
		for userMsg := range publicChan {
			// 获取用户chan
			for uid, user := range usersMap {
				// 跳过产生消息的用户，因为用户自己产生的消息不需要被自己看到，只需要除了自己以外的当前在线的用户看到就行。
				if uid == userMsg.UserID {
					continue
				}
				// 发送消息
				SendMessageToUser(user, userMsg.Msg)
			}
		}
	}
}

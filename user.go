package main

import (
	"fmt"
	"net"
	"time"
)

// 用户结构体
type User struct {
	UserID      string      // 用户ID（暂时用用户IP代替）
	Addr        string      // 用户访问的本地IP
	CreateTime  time.Time   // 创建时间
	MessageChan chan string // 消息通道
}

// 用户消息结构体，记录消息是哪个用户产生的
type UserMsg struct {
	UserID string // 用户ID（暂时用用户IP代替）
	Msg    string // 消息
}

// 一个用户首次登陆，	将消息发送给其它用户
func FirstLoginMsg(user User) {
	// 发送消息到公共频道
	msg := UserMsg{
		UserID: user.UserID,
		Msg:    fmt.Sprintf("用户：[%v] 登陆了\n", user.UserID),
	}
	publicChan <- msg

}

// 生成用户ID
func GetUserID(conn net.Conn) string {
	return conn.RemoteAddr().String()
}

// 将消息发送给用户
func SendMessage(conn net.Conn, user User) {
	for msg := range user.MessageChan {
		conn.Write([]byte(msg)) // 将用户chan下的消息发送给用户
	}
}

// 封装消息发送
func SendMessageToPublic(user User, msg UserMsg) {
	publicChan <- msg
}

func SendMessageToUser(user User, msg string) {
	user.MessageChan <- msg
}

// save user oneline info
func saveUserOnline(userId string, user User) {
	usersMap[userId] = user

}

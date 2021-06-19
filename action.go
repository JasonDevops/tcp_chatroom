package main

// 根据用户输入的内容执行相应的动作

import (
	"fmt"
	"net"
)

// The variable show user that user allow action and detail
var actionList = `
===========================================
|    User allow to execute action lists    |
|                                          |
===========================================
1）quit，msg: quit current user and close connection
2）ulist，msg：show all oneline user lists
`

// Allow action that user execute in the command
// if you have to want add other actions, you just need add action to `allowActions`
var allowActions = []string{
	"list", "quit",
}

// action executeing from user input
func actionOpt(str string, conn net.Conn, user User) {
	switch str {
	case "list":
		user.MessageChan <- actionList
	case "quit":
		actionClose(conn, user)
	}
}

// close user connection ....
func actionClose(conn net.Conn, user User) {
	msg := UserMsg{
		UserID: "",
		Msg:    fmt.Sprintf("用户：[%v] 下线了\n", user.UserID),
	}

	// close tcp connection
	conn.Close()

	// del user info from online map
	delete(usersMap, user.UserID)

	// notify all current online user that the user not online
	publicChan <- msg

}

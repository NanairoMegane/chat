package main

import (
	"github.com/gorilla/websocket"
)

/*
クライアント・モデル。
	socket : ユーザ毎のwebsocket情報
	send : メッセージが送られるチャネル
	room : クライアントが参加しているroomの情報
*/
type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

/*
クライアントのソケットに送られたメッセージを読む混むメソッド。
	・ソケットに届いたメッセージを読み出し、roomに送る。
	・自クライアントの書き込みを察知する？
	・エラーが起きた場合は、ソケットを閉じる。
*/
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

/*
クライアントからの送信メッセージをソケットに書き出すメソッド。
	・メッセージがsendチャンネルに到達したとき、ソケットへの書き出しを行う。
	・他クライアントのメッセージを察知する？
	・エラーが起きた場合は、ソケットを閉じる。
*/
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

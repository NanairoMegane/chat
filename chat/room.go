package main

import (
	"log"
	"net/http"

	"github.com/stretchr/objx"

	"../trace"

	"github.com/gorilla/websocket"
)

/*
チャットルーム・モデル
	forward : 他のクライアントへメッセージを送信するためのチャネル
	join : クライアントのルーム参加を管理するチャネル
	leave : クライアントのルーム退室を管理するチャネル
	clients : ルームに存在するクライアントの情報を保持するmap
*/
type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
	avatar  Avatar
}

/*
ルームの開始メソッド。ルームは無限ループで常に起動する。
	・joinチャンネルでメッセージを受診したら、クライアントmapにクライアントを追加
	・leaveチャンネルでメッセージを受診したら、クライアントmapからクライアントを削除
	・forwardチャンネルでメッセージを受診したら、
	　map内のクライアントのsendチャンネルへメッセージを一斉送信
*/
func (r *room) run() {
	for {
		select {
		//入室
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが入室しました。")
		//退室
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました。")
		//メッセージの送信
		case msg := <-r.forward:
			r.tracer.Trace("新しいメッセージを受信しました。:", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("クライアントに送信されました。")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("送信に失敗しました。")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

/*
room構造体のハンドラの実装。http.Handleのインターフェースに適合する。
	・ページを開いた時にクライアント毎に実行される
	・実行毎にクライアントが生成され、roomのjoinチャンネルに情報を送信する
	・ページ終了時には必ずroomのleaveチャンネルへ情報を送信する
	・別スレッドを構築し、クライアントのwriteチャンネルを無限ループで監視する
	・メインスレッドで、クライアントのreadチャンネルを無限ループで監視する
*/
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP : ", err)
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました：", err)
	}

	//クライアントの生成
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	//入室と退室の管理
	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

/*
新規room作成のヘルパー・メソッド
*/
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

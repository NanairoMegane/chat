package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

/*
テンプレートのハンドラ。
	once     : テンプレートを一度だけコンパイルする
	filename : テンプレートファイルの名称
	temp1    : テンプレートファイルのオブジェクト
*/
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

/*
テンプレートのハンドラ templateHandle への実装。
ここでの ServeHTTPメソッドは、http.Handlerインターフェースに適合する。
*/
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// テンプレートファイル置き場までのパスを取得
	abs, err := filepath.Abs("./chat/templates/")
	if err != nil {
		panic(err)
	}
	// 指定された名称のテンプレートファイルを一度だけコンパイルする
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(abs + t.filename))
	})
	t.temp1.Execute(w, nil)
}

func main() {
	// ルートへのアクセスに対して、chat.htmlを返却する。
	// ハンドラには独自ハンドラを使用する。templateHandleはServeHTTPを実装しているのでインターフェースに適合している。
	/*
				func Handle(pattern string, handler Handler)

					type Handler interface {
		        		ServeHTTP(ResponseWriter, *Request)
					}
	*/
	http.Handle("/", &templateHandler{filename: "/chat.html"})

	/* チャットルームを開始する */
	// 新規roomを作成
	r := newRoom()

	// /roomディレクトリにハンドラを張る。
	// ここへの直接アクセスは行われず、jsから誘導される。
	// roomディレクトリでは、room構造体のServeHTTPがリクエストを処理する。
	http.Handle("/room", r)

	// ルームを起動する
	go r.run()

	/* webサーバを開始する */
	// listenするのは8080のルートだけ。
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

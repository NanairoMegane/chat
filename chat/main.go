package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/objx"

	"../trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var avatars Avatar = TryAvatars{
	UseAuthAvatar,
	UseGravatar,
	UseFileSystemAvatar,
}

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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// テンプレートにレスポンスの入れ物とリクエストの入れ物を渡す。
	// リクエストの中身には、コマンドラインからの入力内容が含まれる。
	// Template.Execute() : テンプレートの内容をwriterに書き込む
	t.temp1.Execute(w, data)
}

func main() {

	/* コマンドラインからポート番号を読み込む */
	// flagパッケージ：コマンドラインからフラグを読み込みパースする
	// flag.String(): フラグを指定名、デフォルト値、Usageで定義する。
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// flag.Parse : コマンドラインの[1:]からフラグを読み込む。
	// これはフラグの定義後、かつプログラムの実行にされなければいけない。
	flag.Parse()

	/* Gomniauth の設定 */
	gomniauth.SetSecurityKey("chattest3594")
	gomniauth.WithProviders(
		google.New("yyy",
			"xxx",
		),
	)

	// ルートへのアクセスに対して、chat.htmlを返却する。
	// ハンドラには独自ハンドラを使用する。templateHandleはServeHTTPを実装しているのでインターフェースに適合している。
	/*
				func Handle(pattern string, handler Handler)

					type Handler interface {
		        		ServeHTTP(ResponseWriter, *Request)
					}
	*/
	http.Handle("/", MustAuth(&templateHandler{filename: "/chat.html"}))

	/* ログインページ */
	http.Handle("/login", &templateHandler{filename: "/login.html"})

	/* 認証機能 */
	http.HandleFunc("/auth/", loginHandler)

	/* ログアウト機能 */
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	/* アップロード機能 */
	http.Handle("/upload", &templateHandler{filename: "/upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	/* プロフィール画像管理機能 */
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./chat/avatars"))))

	/* チャットルームを開始する */
	// 新規roomを作成
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// /roomディレクトリにハンドラを張る。
	// ここへの直接アクセスは行われず、jsから誘導される。
	// roomディレクトリでは、room構造体のServeHTTPがリクエストを処理する。
	http.Handle("/room", r)

	// ルームを起動する
	go r.run()

	/* webサーバを開始する */
	// listenするのは8080のルートだけ。
	log.Println("Webサーバを開始します。ポート：", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

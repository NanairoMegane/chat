package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

/*
次に行うべきハンドラをラップする、認証用の構造体
*/
type authHandler struct {
	next http.Handler
}

/*
認証用の構造体に実装されたメソッド。
	・http.Handle に適合するように作られている。
	・ログイン認証機能を司る
		・過去にログインしたことがあればチャットルーム用のハンドラを起動する。
		過去にログインしていなければ、ログインページに遷移させる。
*/
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 過去に一度も認証を行っていなかった場合
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		// なんらかのエラーが生じた場合
	} else if err != nil {
		panic(err.Error)
		// 成功した場合
	} else {
		// ラップしたハンドラを呼び出す
		h.next.ServeHTTP(w, r)
	}
}

/*
認証用の構造体を生成する。
*/
func MustAuth(handler http.Handler) http.Handler {
	// 引数に受け取ったハンドラを次に行うべきハンドラとして登録しつつ認証用の構造体を生成する。
	return &authHandler{next: handler}
}

/*
ログイン用のハンドラ
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗しました。: ", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました。")
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//log.Println("TODO: ログイン処理", provider)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗しました。: ", provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした。: ", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました。: ", provider, "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(w, "アクション%sには非対応です。", action)
	}
}

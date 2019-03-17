package main

import (
	"errors"
)

// インスタンスがアバター画像を取得できない時のエラーメッセージ
var ErrNoAvatarURL = errors.New("chat:アバターのURLを取得できません。")

/*
ユーザのプロフォール画像を表す型
*/
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

/*
アバター画像の構造体(空)
*/
type AuthAvatar struct{}

// アバター画像型の実装
var UseAuthAvatar AuthAvatar

// AuthAvatar を Avatar 型に適合させる
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	// クライアントからアバターのURLを取り出せるか
	if url, ok := c.userData["avatar_url"]; ok {
		// URLを正常に文字列に変換できるか
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

/*
Gravatar 画像取得用の構造体(空)
*/
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}

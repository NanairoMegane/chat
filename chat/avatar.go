package main

import (
	"errors"
	"os"
)

// インスタンスがアバター画像を取得できない時のエラーメッセージ
var ErrNoAvatarURL = errors.New("chat:アバターのURLを取得できません。")

/*
ユーザのプロフォール画像を表す型
*/
type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			println("GetAvatarURL success.", url)
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

/*
アバター画像の構造体(空)
*/
type AuthAvatar struct{}

// アバター画像型の実装
var UseAuthAvatar AuthAvatar

// AuthAvatar を Avatar 型に適合させる
func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

/*
Gravatar 画像取得用の構造体(空)
*/
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {

	userid := u.UniqueID()
	filename := "/avatars/" + userid + ".jpg"
	if _, err := os.Stat(filename); err == nil {
		println("avatar file is exist. :", filename)
		return filename, nil
	}

	println("avatar file is not exist.")
	return "", ErrNoAvatarURL

}

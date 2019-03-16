package main

import (
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)

	/* 1. クライアントから画像を取り出せないパターン */
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、ErrNoAvatarURLを返すべきです。")
	}
	// テスト用の値のセット
	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{
		"avatar_url": testUrl,
	}

	/* 2. クライアントから画像を取り出せたパターン */
	url, err = authAvatar.GetAvatarURL(client)
	/* 2-1. 画像が取り出せているのに、エラーを返しているパターン*/
	if err != nil {
		t.Error("値が存在する場合、エラーを返すべきではありません。")
	} else {
		/* 2-2. 画像が存在するのに、URLを正しく取り出せていないパターン */
		if url != testUrl {
			t.Error("正常に動作した場合、正しいURLを返すべきです。")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData =
		map[string]interface{}{
			"userid": "8dbd58ddaaa47b6bdd6cbd4f5d7d90e4",
		}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("エラーを返すべきではありません")
	}
	if url != "//www.gravatar.com/avatar/8dbd58ddaaa47b6bdd6cbd4f5d7d90e4" {
		t.Errorf("正しくないURLを取得しました: %s", url)
	}
}

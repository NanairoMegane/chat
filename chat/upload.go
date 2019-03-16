package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userid := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename, _ := filepath.Abs("./templates/avatars/")
	filename += "/" + userid + "/" + header.Filename
	println(filename)
	rep := regexp.MustCompile(".jpg")
	filename = rep.ReplaceAllString(filename, "")
	println(filename)
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "成功")
}

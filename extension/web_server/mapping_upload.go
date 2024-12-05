package web_server

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
)

func init() {
	http.HandleFunc("/upload", upload)
}

func upload(resp http.ResponseWriter, req *http.Request) {
	// 解析multipart/form-data格式的请求体，设置最大内存使用量为1000MB，可按需调整
	err := req.ParseMultipartForm(10 << 22)
	if err != nil {
		log.Errorf("解析multipart/form-data失败: %v", err)
		http.Error(resp, "解析上传请求体失败", http.StatusInternalServerError)
		return
	}

	// 获取上传的文件，这里假设文件表单字段名为"file"，可根据实际情况修改
	file, _, err := req.FormFile("file")
	if err != nil {
		log.Printf("获取文件失败: %v", err)
		http.Error(resp, "获取文件失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//获取文件路径并创建
	filePath := req.FormValue("filePath")
	_, err2 := os.Stat(path.Dir(filePath))
	if err2 != nil {
		os.MkdirAll(path.Dir(filePath), os.ModePerm)
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
		http.Error(resp, "创建文件失败", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// 将上传文件的内容复制到新创建的文件中
	_, err = io.Copy(newFile, file)
	if err != nil {
		log.Printf("复制文件内容失败: %v", err)
		http.Error(resp, "复制文件内容失败", http.StatusInternalServerError)
		return
	}

	resp.Write([]byte("文件上传成功"))
}

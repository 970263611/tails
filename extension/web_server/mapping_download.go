package web_server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/download", download)
}

func download(resp http.ResponseWriter, req *http.Request) {
	// 获取请求的文件名参数，这里假设通过URL路径中的最后一段作为文件名，例如 /download/file.txt 中的 file.txt
	filePath := req.URL.Query().Get("filepath")
	if filePath == "" {
		http.Error(resp, "请指定要下载的文件名", http.StatusBadRequest)
		return
	}

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(resp, "文件不存在", http.StatusNotFound)
		return
	}

	// 打开文件
	fileObj, err := os.Open(filePath)
	if err != nil {
		log.Errorf("打开文件失败: %v", err)
		http.Error(resp, "打开文件失败", http.StatusInternalServerError)
		return
	}
	defer fileObj.Close()

	// 获取文件信息，用于设置响应头相关内容
	fileInfo, err := fileObj.Stat()
	if err != nil {
		log.Printf("获取文件信息失败: %v", err)
		http.Error(resp, "获取文件信息失败", http.StatusInternalServerError)
		return
	}

	// 设置响应头，告知客户端这是一个文件下载，以及文件名、文件大小等信息
	resp.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 将文件内容复制到响应体，实现文件下载
	_, err = io.Copy(resp, fileObj)
	if err != nil {
		log.Errorf("复制文件内容到响应体失败: %v", err)
		http.Error(resp, "文件下载失败", http.StatusInternalServerError)
		return
	}
}

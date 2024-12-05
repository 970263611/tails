package file_download

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
)

type FileDownload struct {
	iface.Context
}

func GetInstance(globalContext iface.Context) iface.Component {
	return &FileDownload{globalContext}
}

func (d *FileDownload) GetName() string {
	return "file_download"
}

func (d *FileDownload) GetDescribe() string {
	return "文件下载 \n例1：浏览器下载 file_download -addr 127.0.0.1:17001 -i /home/test/abc.zip " +
		"\n例2：直接写入本地 file_download -addr 127.0.0.1:17001 -i /home/test/abc.zip -o /home/file/abc.zip"
}

func (d *FileDownload) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "", "addr", true, nil, "下载地址，ip:port")
	cm.AddParameters(cons.STRING, cons.LOWER_I, "", "inputpath", true, nil, "下载文件的全路径(含文件名)")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "", "outputpath", false, nil, "写入本地的全路径(含文件名)，如果该字段为空则直接输出文件流")
}

func (d *FileDownload) Do(params map[string]any) []byte {
	addr := params["addr"].(string)
	inputPath := params["inputpath"].(string)
	opath, ok := params["outputpath"] //如果没有送-o，则直接将文件返回
	var writer http.ResponseWriter
	var outputPath string
	if !ok {
		writer = d.GetResponseWriter()
		if writer == nil {
			msg := fmt.Sprintf("本地请求时 写入本地路径 不能为空")
			log.Error(msg)
			return []byte(msg)
		}
	} else {
		outputPath = opath.(string)
	}
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/download?filepath="+inputPath, nil)
	if err != nil {
		msg := fmt.Sprintf("创建请求失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("下载远程文件失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	body := resp.Body
	defer body.Close()
	result, _ := io.ReadAll(body)
	if resp.StatusCode != http.StatusOK {
		msg := "下载远程文件失败：" + string(result)
		log.Error(msg)
		return []byte(msg)
	}
	if outputPath == "" {
		writer.Header().Set("Content-Disposition", "attachment; filename="+path.Base(inputPath))
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(result)))
		return result
	}
	//获取文件路径并创建
	_, err2 := os.Stat(path.Dir(outputPath))
	if err2 != nil {
		os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	}
	newFile, err := os.Create(outputPath)
	if err != nil {
		msg := fmt.Sprintf("创建文件失败: %v", err)
		log.Error("创建文件失败: %v", err)
		return []byte(msg)
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, body)
	newFile.Write(result)
	if err != nil {
		msg := fmt.Sprintf("写出文件失败: %v", err)
		log.Error(msg)
		return []byte(msg)
	}
	return []byte("文件下载成功")
}

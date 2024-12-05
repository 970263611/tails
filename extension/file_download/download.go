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

type FileDownload struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &FileDownload{}
}

func (c *FileDownload) GetName() string {
	return "file_download"
}

func (c *FileDownload) GetDescribe() string {
	return "文件下载"
}

func (d *FileDownload) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "", "addr", true, nil, "下载地址，ip:port")
	cm.AddParameters(cons.STRING, cons.LOWER_I, "", "inputpath", true, nil, "下载文件的全路径(含文件名)")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "", "outputpath", true, nil, "写入本地的全路径(含文件名)")
}

func (d *FileDownload) Do(params map[string]any) []byte {
	addr := params["addr"].(string)
	inputPath := params["inputpath"].(string)
	outputPath := params["outputpath"].(string)
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/download?filepath="+inputPath, nil)
	if err != nil {
		msg := fmt.Sprintf("创建请求失败：%v", err.Error())
		log.Error(msg)
		return []byte(msg)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("下载远程文件失败：%v", err.Error())
		log.Error(msg)
		return []byte(msg)
	}
	body := resp.Body
	defer body.Close()
	if resp.StatusCode != http.StatusOK {
		result, _ := io.ReadAll(body)
		msg := "下载远程文件失败：" + string(result)
		log.Error(msg)
		return []byte(msg)
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
	if err != nil {
		msg := fmt.Sprintf("写出文件失败: %v", err)
		log.Error(msg)
		return []byte(msg)
	}
	return []byte("文件下载成功")
}

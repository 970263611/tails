package file_upload

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileUpload struct{}

func GetInstance(globalContext iface.Context) iface.Component {
	return &FileUpload{}
}

func (c *FileUpload) GetName() string {
	return "file_upload"
}

func (c *FileUpload) GetDescribe() string {
	return "文件上传 \n例：file_upload -a 127.0.0.1:17001  -i /home/test/abc.zip -o /home/file/abc.zip"
}

func (d *FileUpload) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "", "addr", true, nil, "上传地址，ip:port")
	cm.AddParameters(cons.STRING, cons.LOWER_I, "", "inputpath", true, nil, "本地文件全路径(含文件名)")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "", "outputpath", true, nil, "输出文件全路径(含文件名)")
}

func (d *FileUpload) Do(params map[string]any) []byte {
	addr := params["addr"].(string)
	inputPath := params["inputpath"].(string)
	outputPath := params["outputpath"].(string)
	//读取文件
	file, err := os.Open(inputPath)
	if err != nil {
		msg := fmt.Sprintf("读取文件失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	defer file.Close()
	//创建表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", inputPath)
	if err != nil {
		msg := fmt.Sprintf("创建文件表单字段失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		msg := fmt.Sprintf("复制文件内容到表单字段失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	field, err := writer.CreateFormField("filepath")
	if err != nil {
		msg := fmt.Sprintf("创建文件表单字段失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	_, err = field.Write([]byte(outputPath))
	if err != nil {
		msg := fmt.Sprintf("复制文件内容到表单字段失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	err = writer.Close()
	if err != nil {
		msg := fmt.Sprintf("关闭multipart写入器失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	//发送请求
	req, err := http.NewRequest(http.MethodPost, "http://"+addr+"/upload", body)
	if err != nil {
		msg := fmt.Sprintf("创建请求失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("上传文件失败：%v", err)
		log.Error(msg)
		return []byte(msg)
	}
	//响应返回
	respBody := resp.Body
	defer respBody.Close()
	result, _ := io.ReadAll(respBody)
	return result
}

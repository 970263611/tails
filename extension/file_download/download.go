package file_download

import (
	cons "basic/constants"
	iface "basic/interfaces"
	"basic/tool/utils"
	"errors"
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
	return "文件下载，支持相对和绝对路径，相对路径的根路径tails文件服务所在目录，绝对路径为文件全路径 \n例1：浏览器下载 file_download -addr 127.0.0.1:17001 -i /home/test/abc.zip " +
		"\n例2：直接写入本地 file_download -addr 127.0.0.1:17001 -i /home/test/abc.zip -o /home/file/abc.zip" +
		"\n例3：本地文件操作 file_download -i /home/test/abc.zip -o /home/file/abc.zip，相当于cp命令"
}

func (d *FileDownload) Register(cm iface.ComponentMeta) {
	cm.AddParameters(cons.STRING, cons.LOWER_A, "", "addr", false, nil, "下载地址：ip:port，该字段为空时读取本地文件")
	cm.AddParameters(cons.STRING, cons.LOWER_I, "", "inputpath", true, nil, "远程文件路径(含文件名)")
	cm.AddParameters(cons.STRING, cons.LOWER_O, "", "outputpath", false, nil, "本地输出路径(含文件名)，该字段为空则直接输出文件流")
}

func (d *FileDownload) Do(params map[string]any) []byte {
	iaddr, ok1 := params["addr"]
	inputpath := params["inputpath"].(string)
	ipath, ok2 := params["outputpath"] //如果没有送-o，则直接将文件返回
	var writer http.ResponseWriter
	//输出到本地还是http下载
	if !ok2 {
		writer = d.GetResponseWriter()
		if writer == nil {
			msg := fmt.Sprintf("本地请求时 写入本地路径 不能为空")
			log.Error(msg)
			return []byte(msg)
		}
	}
	var result []byte
	var err error
	//判断是读取本地还是远程下载
	if ok1 {
		result, err = downloadFile(iaddr.(string), inputpath)
	} else {
		result, err = readLocal(inputpath)
	}
	if err != nil {
		return []byte(err.Error())
	}
	//判断是写出到本地还是返回二进制流
	if !ok2 {
		writer.Header().Set("Content-Disposition", "attachment; filename="+path.Base(inputpath))
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(result)))
		return result
	} else {
		err = writeFile(ipath.(string), result)
	}
	if err != nil {
		return []byte(err.Error())
	}
	return []byte("文件下载成功")
}

/*
*
写出文件到本地
*/
func writeFile(outputPath string, file []byte) error {
	//获取文件路径并创建
	outputPath = utils.GetAbsolutePath(outputPath)
	_, err2 := os.Stat(utils.PathDir(outputPath))
	if err2 != nil {
		os.MkdirAll(utils.PathDir(outputPath), os.ModePerm)
	}
	newFile, err := os.Create(outputPath)
	if err != nil {
		msg := fmt.Sprintf("创建文件失败: %v", err)
		log.Error("创建文件失败: %v", err)
		return errors.New(msg)
	}
	defer newFile.Close()
	newFile.Write(file)
	if err != nil {
		msg := fmt.Sprintf("写出文件失败: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

/*
*
远程下载文件
*/
func downloadFile(addr string, inputpath string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+"/download?filepath="+inputpath, nil)
	if err != nil {
		msg := fmt.Sprintf("创建请求失败：%v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("下载远程文件失败：%v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	body := resp.Body
	defer body.Close()
	result, _ := io.ReadAll(body)
	if resp.StatusCode != http.StatusOK {
		msg := "下载远程文件失败：" + string(result)
		log.Error(msg)
		return nil, errors.New(msg)
	}
	return result, nil
}

/*
*
读取本地文件
*/
func readLocal(inputpath string) ([]byte, error) {
	// 检查文件是否存在
	filebyte, err := os.ReadFile(utils.GetAbsolutePath(inputpath))
	if err != nil {
		msg := "读取本地文件失败：" + err.Error()
		log.Error(msg)
		return nil, err
	}
	return filebyte, nil
}

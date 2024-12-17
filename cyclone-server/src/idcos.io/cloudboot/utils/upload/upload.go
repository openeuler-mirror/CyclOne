package upload

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/voidint/binding"
)

var (
	//UploadDir 批量导入放置的位置
	UploadDir = filepath.Join(os.TempDir(), "cloudboot-server") + string(os.PathSeparator)
)

//FileExist 检查文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//GenerateTempFile 加载导入文件
func GenerateTempFile(r *http.Request, prefix string) (string, error) {
	r.ParseForm()
	file, handle, err := r.FormFile("files[]")
	if err != nil {
		return "", err
	}

	if !FileExist(UploadDir) {
		err := os.MkdirAll(UploadDir, 0777)
		if err != nil {
			return "", err
		}
	}

	list := strings.Split(handle.Filename, ".")
	ext := list[len(list)-1]

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano()) + handle.Filename))
	cipherStr := h.Sum(nil)
	md5 := fmt.Sprintf("%s", hex.EncodeToString(cipherStr))
	filename := prefix + md5 + "." + ext

	if FileExist(UploadDir + filename) {
		os.Remove(UploadDir + filename)
	}

	f, err := os.OpenFile(UploadDir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer file.Close()
	return filename, nil
}

//ImportReq 导入请求参数
type ImportReq struct {
	// 文件名称
	FileName string `json:"file_name"`
	// 页大小
	Limit uint `json:"limit"`
	// 页号
	Offset uint `json:"offset"`
	//UserName 登入用户
	UserName string `json:"-"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *ImportReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Limit:    "limit",
		&reqData.Offset:   "offset",
		&reqData.FileName: "file_name",
	}
}

//Validate 参数校验
func (reqData *ImportReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	//这里的校验要与model状态保持一致
	if reqData.FileName == "" {
		errs.Add([]string{"filename"}, binding.RequiredError, "文件名不能为空")
		return errs
	}
	if !FileExist(UploadDir + reqData.FileName) {
		errs.Add([]string{"filename"}, binding.RequiredError, fmt.Sprintf("文件名(%s)不存在", UploadDir+reqData.FileName))
		return errs
	}
	return errs
}

const (
	//Continue 可以继续检查
	Continue = 1
	//Return 不能继续检查
	Return = 2
	//DO 可以继续向下执行
	DO = 3
)

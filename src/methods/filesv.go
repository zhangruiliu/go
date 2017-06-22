package methods
import (
	"sync"
	"mime/multipart"
	"errors"
	"os"
	"io"
)

var bufsize = 1024 * 10

//池
var pool = sync.Pool{
	New : func() interface{} {
		buf := make([]byte, bufsize)
		return buf
	},
}

//SaveMultiPartFileAs 将http 上传的二进制文件存放
func SaveMultiPartFileAs(f multipart.File, fileName string) error {
	if f == nil {
		return errors.New("multipart.File 内容为空")
	}
	defer f.Close()
	sf, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer sf.Close()

	buf := pool.Get().([]byte)
	for {

		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if n <= 0 {
			break
		}
		_, err = sf.Write(buf[0:n])
		if err != nil {
			return err
		}
	}
	pool.Put(buf)
	return nil

}
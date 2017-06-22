package jsonconf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// ConfFromFile 从文件中生成配置
func ConfFromFile(filename string, data interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return ConfFromBytes(buf, data)
}

// ConfFromString 从字符串中生成配置
func ConfFromString(str string, data interface{}) error {
	return ConfFromBytes([]byte(str), data)
}

// ConfFromBytes 从字节中生成配置
func ConfFromBytes(buf []byte, data interface{}) error {
	nbuf, err := removeComment(buf)
	if err != nil {
		return err
	}
	if nbuf[0] == byte(0xEF) && nbuf[1] == byte(0xBB) && nbuf[2] == byte(0xBF) { // bom
		nbuf = nbuf[3:]
	}
	return json.Unmarshal(nbuf, data)
}

func removeComment(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	bufreader := bufio.NewReader(bytes.NewBuffer(data))
	for {
		line, _, err := bufreader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if bytes.HasPrefix((bytes.TrimSpace(line)), []byte("//")) {
			continue
		}
		buf.Write(line)
	}
	return buf.Bytes(), nil

}

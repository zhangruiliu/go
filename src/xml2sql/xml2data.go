package xml2sql

import (
	"io"
	"encoding/xml"
)

var (
	sqlStartElemet = "sql"
)



//解析元素类型为sql的
func NewDataFromReaders(rs []io.Reader, config *DataConfig, data *Data) (*Data, error) {
	//data := New()
	if data == nil {
		data = New()
	}
	if config == nil {

	}else {
		data.Config = config
		sqlStartElemet = data.Config.SqlStartElemet
	}
	for _, r := range rs {
		decode := xml.NewDecoder(r)
		var (
			token xml.Token
			err error
			hasStartElement bool
			name string
		)
		for {
			token, err = decode.Token()
			if err != nil {
				if err == io.EOF {
					break
				}else {
					return nil, err
				}
			}
			switch token.(type) {
			case xml.StartElement:
				ele := token.(xml.StartElement)
				if ele.Name.Local != sqlStartElemet {
					continue
				}
				for _, v := range ele.Attr {
					if v.Name.Local == "name" {
						name = v.Value
					}
					if v.Name.Local == "type" && v.Value == "sql" {
						hasStartElement = true
					}
				}
				if name != "" && hasStartElement {
					hasStartElement = true
				}
			case xml.CharData:
				ele := token.(xml.CharData)
				if hasStartElement {
					data.XML[name] += string(ele)
				}
			case xml.EndElement:
				hasStartElement = false
				name = ""
			}
		}

	}
	return data, nil
}
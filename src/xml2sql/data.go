package xml2sql
import (
	"encoding/json"
)


type Data struct {
	XML    map[string]string
	Config *DataConfig
}



func New() *Data {
	data := &Data{}
	data.XML = make(map[string]string)
	return data
}

type DataConfig struct {
	SqlStartElemet string
}

func NewConfig(js string) *DataConfig {
	config := &DataConfig{}
	err := json.Unmarshal([]byte(js), config)
	if err != nil{
		panic(err)
	}
	return config

}
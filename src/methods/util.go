package methods

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"strings"
	//"sql2ety"
	//"dbls/lsentity"
	//"logic/limsError"
)

func NotEscapeSqlStringT(str string) string {
	t1 := strings.Replace(str, `'`, `''`, -1)
	return t1
}

//EscapeSqlString 转义sql 中的字符串
func EscapeSqlString(str string) string {
	t1 := strings.Replace(str, `'`, `''`, -1)
	t1 = strings.TrimSpace(t1)
	t1 = strings.Replace(t1, "\u0000", "", -1)
	// t1 = strings.TrimFunc(t1, func(r rune) bool {
	// 	if r == '\u0000' {
	// 		return true
	// 	}
	// 	return false
	// })
	t1 = strings.TrimSpace(t1)

	return `'` + t1 + `'`
}

//从http.Request中读取 http.Request 并且 关闭
func ReadFromRequest(r *http.Request) (ret []byte) {
	defer func() {
		if r := recover(); r != nil {
			ret = nil
		}
	}()

	//defer r.Body.Close()

	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return rb
}

////从http.Request中读取 http.Request 并且 关闭
//func ReadFromRequest(r *http.Request) []byte {
//	defer r.Body.Close()
//	rB, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		panic(err)
//	}
//	return rB
//}

func AlleleCount(gene string) int64 {
	gene = strings.TrimSpace(gene)
	if gene == "" {
		return 0
	}

	ss := strings.Split(gene, ";")
	var count = int64(0)
	for _, v := range ss {
		if v != "" {
			count++
		}
	}
	return count
}

//sample_dna_source 表中位点个数个数
func SampleDnaSourceAlleleCount(gene string) int64 {
	gene = strings.TrimSpace(gene)

	nlgene := strings.Split(gene, "#")
	ngene := nlgene[0]
	lgene := ""
	if len(nlgene) > 1 {
		lgene = nlgene[1]
	}

	var nCount = int64(0)
	sngene := strings.Split(ngene, ";")
	for _, v := range sngene {
		if v != "" {
			nCount++
		}
	}
	lCount := AlleleCount(lgene)
	return int64(nCount) + lCount

}

func GetString(runs ...rune) string {
	return string(runs)
}

//
//func Insert(data interface{}, db sql2ety.LimsDBInterface, tbname ... string) error {
//	stmt, err := lsentity.MakeInsertSql(data, db.GetDbType(), tbname...)
//	if err != nil {
//		return limsError.New().Add(err)
//	}
//
//	_, err = db.Exec(stmt)
//	if err != nil {
//		return limsError.New().Add(err)
//	}
//	return nil
//}

package card

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"g"
	"io"
	"io/ioutil"
	"lsent"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"timex"
	"unsafe"
	"web/httpcontext"
	"xlsxtool"
	"xml2sql"

	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	"github.com/tealeg/xlsx"
)

type Card struct {
	httpcontext.HttpContext
}
type Personsex struct {
	lsent.Personinfo
	Idcardnocheck string `json:"id_card_no_check"` //0重复不插 1必须插入
	Photo         string `json:"photo"`
}
type IdCardNoCheck struct {
	Id_card_no_check string `json:"id_card_no_check"`
}
type Cperson struct {
	Person_Info   Personsex     `json:"person_info"`
	IdCardNoCheck IdCardNoCheck `json:"id_card_no_check"`
}

func StoB(s string) []byte {
	return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s))))
}

//新增 修改
func (card *Card) Createperandsamp() {
	var (
		someper  = Cperson{}
		sys_user = lsent.Sys_User{}
	)
	card.UnJson(&someper)
	card.DB().SingularTable(true)
	r, err := card.R.Cookie("uid")
	card.CheckWrite(err)
	if r != nil {
		card.CheckWrite(card.DB().Where("id=?", r.Value).Find(&sys_user).Error)
	}
	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()

	if someper.Person_Info.Id != "" {
		if someper.Person_Info.SampleLabNo == "" {
			err = errors.New("样本实验室编号不能为空")
			card.CheckWrite(err)
		}
		if someper.Person_Info.IdCardNo == "" {
			err = errors.New("身份证号不能为空")
			card.CheckWrite(err)
		}
		perlist := SelectpersoBySamplelabno(someper.Person_Info.SampleLabNo)
		checkper := SelectcheckpersonBySamplelabno(someper.Person_Info.SampleLabNo)
		if len(perlist) > 0 || len(checkper) > 0 {
			err = errors.New("样本实验室编号已存在")
			card.CheckWrite(err)
		}
		re, _ := SelectCheckPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
		offre, _ := SelectPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
		if len(re) == 0 && len(offre) == 0 {
			someper.Person_Info.Race = SelectDictByCategory(someper.Person_Info.Race)
			someper.Person_Info.Gender = SelectDictByCategory(someper.Person_Info.Gender)
			someper.Person_Info.UpdateUser = sys_user.UserName
			someper.Person_Info.UpdateDatetime = time.Now().Format("2006-01-02 15:04:05")
			card.CheckWrite(db.Save(&someper.Person_Info.Personinfo).Error)
			card.CheckWrite(db.Commit().Error)
			card.JsonWrite(someper.Person_Info)
		} else if len(offre) > 0 {
			card.JsonWritefailure(offre)
		} else if len(re) > 0 {
			card.JsonWritefailure(re)
		}
	} else {
		// if someper.Person_Info.SampleLabNo == "" {
		// 	err = errors.New("样本实验室编号为空")
		// 	card.CheckWrite(err)
		// }
		if someper.Person_Info.IdCardNo == "" {
			err = errors.New("身份证号不能为空")
			card.CheckWrite(err)
		}
		if someper.IdCardNoCheck.Id_card_no_check == "1" {
			// re, _ := SelectCheckPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
			offre, err4 := SelectPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
			card.CheckWrite(err4)
			if len(offre) == 0 {
				someper.Person_Info.Id = card.Uuid()
				someper.Person_Info.Race = SelectDictByCategory(someper.Person_Info.Race)
				someper.Person_Info.Gender = SelectDictByCategory(someper.Person_Info.Gender)
				someper.Person_Info.CreateUser = sys_user.UserName
				someper.Person_Info.CollectUser = sys_user.UserName
				someper.Person_Info.CreateDatetime = time.Now().Format("2006-01-02 15:04:05")
				someper.Person_Info.CollectDatetime = time.Now().Format("2006-01-02 15:04:05")
				someper.Person_Info.SampleType = "01"
				if someper.Person_Info.Photo != "" {
					files := strings.Split(someper.Person_Info.Photo, ",")[1]
					res, err8 := base64.StdEncoding.DecodeString(files)
					card.CheckWrite(err8)
					someper.Person_Info.Foot = res
				}
				// code, err11 := GetConfigvaluelocal("collect-agency-code")
				// card.CheckWrite(err11)
				code, err11 := Getconfigno("7821")
				card.CheckWrite(err11)
				someper.Person_Info.CollectAgencyCode = code
				// code2, err112 := GetConfigvaluelocal("collect-agency-name")
				// card.CheckWrite(err112)
				code2, err112 := Getconfigno("7822")
				card.CheckWrite(err112)
				someper.Person_Info.CollectAgencyName = code2
				// server, err112 := GetConfigvaluelocal("server")
				// card.CheckWrite(err112)
				server, err113 := Getconfigno("7824")
				card.CheckWrite(err113)
				someper.Person_Info.Server = server
				someper.Person_Info.Flag = 0
				if someper.Person_Info.SampleLabNo == "" {
					labno, err0 := Getsamplelabno()
					if err0 != nil {
						card.CheckWrite(err0)
					}
					someper.Person_Info.SampleLabNo = labno
					perlist := SelectpersoBySamplelabno(someper.Person_Info.SampleLabNo)
					checkper := SelectcheckpersonBySamplelabno(someper.Person_Info.SampleLabNo)
					if len(perlist) > 0 || len(checkper) > 0 {
						err = errors.New("样本实验室编号已存在")
						card.CheckWrite(err)
					}
				} else {
					perlist := SelectpersoBySamplelabno(someper.Person_Info.SampleLabNo)
					checkper := SelectcheckpersonBySamplelabno(someper.Person_Info.SampleLabNo)
					if len(perlist) > 0 || len(checkper) > 0 {
						err = errors.New("样本实验室编号已存在")
						card.CheckWrite(err)
					}
				}
				fmt.Println(someper.Person_Info.Personinfo.CreateDatetime)
				card.CheckWrite(db.Create(&someper.Person_Info.Personinfo).Error)
				card.CheckWrite(db.Commit().Error)
				card.JsonWrite(someper.Person_Info)
			} else if len(offre) > 0 {
				card.JsonWritefailure(offre)
			}
		} else if someper.IdCardNoCheck.Id_card_no_check == "" || someper.IdCardNoCheck.Id_card_no_check == "0" {
			re, _ := SelectCheckPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
			offre, err4 := SelectPersonByIdcardNoLine(someper.Person_Info.IdCardNo)
			card.CheckWrite(err4)
			if len(re) == 0 && len(offre) == 0 {
				someper.Person_Info.Id = card.Uuid()
				someper.Person_Info.Race = SelectDictByCategory(someper.Person_Info.Race)
				someper.Person_Info.Gender = SelectDictByCategory(someper.Person_Info.Gender)
				someper.Person_Info.CreateUser = sys_user.UserName
				someper.Person_Info.CollectUser = sys_user.UserName
				someper.Person_Info.CreateDatetime = time.Now().Format("2006-01-02 15:04:05")
				someper.Person_Info.CollectDatetime = time.Now().Format("2006-01-02 15:04:05")
				someper.Person_Info.SampleType = "01"
				if someper.Person_Info.Photo != "" {
					files := strings.Split(someper.Person_Info.Photo, ",")[1]
					res, err8 := base64.StdEncoding.DecodeString(files)
					card.CheckWrite(err8)
					someper.Person_Info.Foot = res
				}
				// code, err11 := GetConfigvaluelocal("collect-agency-code")
				// card.CheckWrite(err11)
				code, err11 := Getconfigno("7821")
				card.CheckWrite(err11)
				someper.Person_Info.CollectAgencyCode = code
				// code2, err112 := GetConfigvaluelocal("collect-agency-name")
				// card.CheckWrite(err112)
				code2, err112 := Getconfigno("7822")
				card.CheckWrite(err112)
				someper.Person_Info.CollectAgencyName = code2
				// server, err112 := GetConfigvaluelocal("server")
				// card.CheckWrite(err112)
				server, err113 := Getconfigno("7824")
				card.CheckWrite(err113)
				someper.Person_Info.Server = server
				someper.Person_Info.Flag = 0
				if someper.Person_Info.SampleLabNo == "" {
					labno, err0 := Getsamplelabno()
					if err0 != nil {
						card.CheckWrite(err0)
					}
					someper.Person_Info.SampleLabNo = labno
					perlist := SelectpersoBySamplelabno(someper.Person_Info.SampleLabNo)
					checkper := SelectcheckpersonBySamplelabno(someper.Person_Info.SampleLabNo)
					if len(perlist) > 0 || len(checkper) > 0 {
						err = errors.New("样本实验室编号已存在")
						card.CheckWrite(err)
					}
				} else {
					perlist := SelectpersoBySamplelabno(someper.Person_Info.SampleLabNo)
					checkper := SelectcheckpersonBySamplelabno(someper.Person_Info.SampleLabNo)
					if len(perlist) > 0 || len(checkper) > 0 {
						err = errors.New("样本实验室编号已存在")
						card.CheckWrite(err)
					}
				}
				fmt.Println(someper.Person_Info.Personinfo.CreateDatetime)
				card.CheckWrite(db.Create(&someper.Person_Info.Personinfo).Error)
				card.CheckWrite(db.Commit().Error)
				card.JsonWrite(someper.Person_Info)
			} else if len(offre) > 0 {
				card.JsonWritefailure(offre)
			} else if len(re) > 0 {
				card.JsonWritefailure(re)
			}
		}
	}
}

type DemoT struct {
	A string `json:"a"`
	B string `json:"b"`
}

// func (card *Card) Selectperandsamp() {
// 	var person []lsent.Per
// 	card.UnJson(&person)

// 	stmt := `select count(*) from per`
// 	card.DB().Raw(stmt).Scan(&person)
// 	card.JsonWrite(person)
// 	var (
// 		sql string
// 		ok  bool
// 	)
// 	if sql, ok = g.Sqls["demo"]; !ok {

// 	}
// 	Demo := DemoT{
// 		"A",
// 		"",
// 	}
// 	sql, _ = xml2sql.GenerateSql(sql, &Demo)
// 	fmt.Println(sql)
// 	card.DB()
// }

type someperson struct {
	Person_Info lsent.Personinfo `json:"person_info"`
}

// func (card *Card) Createperandsamp() {
// 	var (
// 		person  *lsent.Personinfo
// 		someper *someperson
// 	)
// 	card.DB().SingularTable(true)
// 	card.UnJson(&someper)
// 	db := card.DB().Begin()
// 	defer func() {
// 		if db.Error != nil {
// 			db.Rollback()
// 		}
// 	}()
// 	if someper.Person_Info.Id != "" {
// 		card.CheckWrite(card.DB().Where("id = ?", someper.Person_Info.Id).Find(&person).Error)
// 		personi := &lsent.Personinfo{}
// 		personi.PersonName = someper.Person_Info.PersonName
// 		personi.Gender = someper.Person_Info.Gender
// 		personi.Race = someper.Person_Info.Race
// 		personi.BirthDatetime = someper.Person_Info.BirthDatetime
// 		personi.CollectAgencyCode = someper.Person_Info.CollectAgencyCode
// 		personi.CollectAgencyName = someper.Person_Info.CollectAgencyName
// 		personi.CollectDatetime = someper.Person_Info.CollectDatetime
// 		personi.IdCardNo = someper.Person_Info.IdCardNo
// 		personi.NativePlaceAddr = someper.Person_Info.NativePlaceAddr
// 		personi.NativePlaceRegion = someper.Person_Info.NativePlaceRegion
// 		personi.UpdateDatetime = timex.Time(time.Now())
// 		personi.CreateDatetime = someper.Person_Info.CreateDatetime
// 		card.CheckWrite(db.Save(&personi).Error)
// 		card.CheckWrite(db.Commit().Error)
// 		card.JsonWrite(personi)
// 	} else {
// 		someper.Person_Info.Id = card.Uuid()
// 		// person.UpdateDatetime = timex.Time(time.Now())
// 		someper.Person_Info.CollectDatetime = timex.Time(time.Now())
// 		someper.Person_Info.CreateDatetime = timex.Time(time.Now())
// 		card.CheckWrite(db.Create(&someper.Person_Info).Error)
// 		card.CheckWrite(db.Commit().Error)
// 		card.JsonWrite(someper.Person_Info)
// 	}
// }

type line struct {
	Id string `json:"id"`
}

func (card *Card) GetConfigvalue() {
	var linevalue line
	card.UnJson(&linevalue)
	buf, err := ioutil.ReadFile("conf.json")
	card.CheckWrite(err)
	data, err := removeComment(buf)
	if data[0] == byte(0xEF) && data[1] == byte(0xBB) && data[2] == byte(0xBF) { // bom
		data = data[3:]
	}
	card.CheckWrite(err)
	con, err := simplejson.NewJson(data)
	card.CheckWrite(err)
	if linevalue.Id == "" {
		card.JsonWrite(con)
	} else {
		if con != nil {
			if _, b := con.CheckGet(linevalue.Id); b {
				port, _ := con.Get(linevalue.Id).String()
				fmt.Println(port)
				card.JsonWrite(port)
			} else {
				err = errors.New("无此字段")
				card.CheckWrite(err)
			}
		}

	}

	// card.CheckWrite(err)
	// ra := reflect.ValueOf(&con)
	// fmt.Println(ra.Elem().FieldByName("off-line").Interface())
}
func GetConfigvaluelocal(id string) (result string, err7 error) {
	var (
		err error
	)
	buf, err2 := ioutil.ReadFile("conf.json")
	err = err2
	if err != nil {
		return "", err
	}
	data, err3 := removeComment(buf)
	err = err3
	if err != nil {
		return "", err
	}
	if data[0] == byte(0xEF) && data[1] == byte(0xBB) && data[2] == byte(0xBF) { // bom
		data = data[3:]
	}
	con, err4 := simplejson.NewJson(data)
	if err4 != nil {
		return "", err4
	}
	err = err4
	if con != nil {
		if _, b := con.CheckGet(id); b {
			port, _ := con.Get(id).String()
			fmt.Println(port)
			result = port
			return result, err
		} else {
			return result, err
		}
	}
	return result, err
	// card.CheckWrite(err)
	// ra := reflect.ValueOf(&con)
	// fmt.Println(ra.Elem().FieldByName("off-line").Interface())
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

type syslist struct {
	Result []lsent.Sys_Dict `json:"result"`
	Status string           `json:"status"`
}

//下载字典
func (card *Card) Getallsysdict() {
	sqls := "select * from sys_dict where dict_key in('7821','7822','7823','7824','7825','7826','7827','7828')"
	var (
		sys []lsent.Sys_Dict
	)
	card.CheckWrite(card.DB().Raw(sqls).Scan(&sys).Error)
	sys_list := syslist{}
	buf, err := json.Marshal("")
	online, err11 := GetConfigvaluelocal("public")
	card.CheckWrite(err11)
	resp, err := http.Post(online+"/api/"+"SelectSysDict", "application/json", bytes.NewReader(buf))
	card.CheckWrite(err)
	body, err := ioutil.ReadAll(resp.Body)
	card.CheckWrite(err)
	defer resp.Body.Close()
	err = json.Unmarshal(body, &sys_list)
	card.DB().SingularTable(true)

	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	dict := []lsent.Sys_Dict{}
	card.CheckWrite(db.Raw("delete  from sys_dict").Scan(&dict).Error)
	for _, v := range sys_list.Result {
		v.Id = card.Uuid()
		fmt.Println(v)
		card.CheckWrite(db.Create(&v).Error)
	}
	for _, k := range sys {
		if k.DictKey == "7821" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7822" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7823" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7824" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7825" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7826" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7827" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		} else if k.DictKey == "7828" {
			err3 := db.Exec("update sys_dict set dict_value1=" + "'" + k.DictValue1 + "'" + "where dict_key='" + k.DictKey + "'").Error
			card.CheckWrite(err3)
		}
	}
	var (
		dict2 lsent.Sys_Dict
	)
	dict2.Id = card.Uuid()
	dict2.DictKey = "01"
	dict2.DictNationalKey = "01"
	dict2.DictValue1 = "血斑/血液"
	dict2.DictCategory = "sample_type"
	card.CheckWrite(db.Create(&dict2).Error)
	card.CheckWrite(db.Commit().Error)
	card.JsonWrite(sys_list)
}

type Server struct {
	Result []lsent.Server `json:"result"`
	Status string         `json:"status"`
}

//下载服务器列表
func (card *Card) Getallserver() {
	var server Server
	updatedb, err := GetConfigvaluelocal("getcheckserver")
	card.CheckWrite(err)
	resp, err := http.Get(updatedb)
	card.CheckWrite(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &server)
	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	serverLocal := []lsent.Server{}
	card.CheckWrite(db.Raw("delete from server").Scan(&serverLocal).Error)
	for _, k := range server.Result {
		card.CheckWrite(db.Create(&k).Error)
	}
	// card.CheckWrite(db.Commit().Error)
	err = db.Commit().Error
	if err != nil {
		db.Error = err
		panic(db.Error)
	}
	card.JsonWrite(serverLocal)
}

type sysuserlist struct {
	Result []lsent.Sys_User `json:"result"`
	Status string           `json:"status"`
}

//下载用户SelectSysUser
func (card *Card) Getallsysuser() {
	sys_user_list := sysuserlist{}
	buf, err := json.Marshal("")
	online, err11 := GetConfigvaluelocal("public")
	card.CheckWrite(err11)
	resp, err := http.Post(online+"/api/"+"SelectSysUser", "application/json", bytes.NewReader(buf))
	card.CheckWrite(err)
	body, err := ioutil.ReadAll(resp.Body)
	card.CheckWrite(err)
	defer resp.Body.Close()
	err = json.Unmarshal(body, &sys_user_list)
	card.DB().SingularTable(true)
	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	dict := []lsent.Sys_User{}
	card.CheckWrite(db.Raw("delete  from sys_user").Scan(&dict).Error)
	for _, k := range sys_user_list.Result {
		k.Id = card.Uuid()
		fmt.Println(k)
		card.CheckWrite(db.Create(&k).Error)
	}
	card.CheckWrite(db.Commit().Error)
	card.JsonWrite(sys_user_list)
}

//登录
func (card *Card) Login() {
	card.DB().SingularTable(true)
	loguss := &lsent.Sys_User{}
	// loguser := &lsent.Usr{}
	var (
		password string
	)
	if card.R.Method == "GET" {
		loguss.UserName = card.R.URL.Query().Get("user_name")
		loguss.Password = card.R.URL.Query().Get("password")
		password = loguss.Password
		md5Pwd := md5.Sum([]byte(loguss.Password))
		loguss.Password = strings.ToUpper(hex.EncodeToString(md5Pwd[:]))
	} else {
		card.UnJson(loguss)
		password = loguss.Password
		md5Pwd := md5.Sum([]byte(loguss.Password))
		loguss.Password = strings.ToUpper(hex.EncodeToString(md5Pwd[:]))
	}
	fmt.Print(loguss.UserName)
	fmt.Println(loguss.Password)
	nomiUser := lsent.Sys_User{}
	err2 := card.DB().Where("user_name = ?", loguss.UserName).Find(&nomiUser).Error
	card.CheckWrite(err2)
	if card.R.Method == "GET" {
		if strings.ToUpper(nomiUser.Password) == strings.ToUpper(password) || strings.ToUpper(nomiUser.Password) == strings.ToUpper(loguss.Password) {
			cookie2 := http.Cookie{
				Name:  "init_server_no",
				Value: string(nomiUser.InitServerNo),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie2)
			cookie3 := http.Cookie{
				Name:  "user_name",
				Value: url.QueryEscape(nomiUser.UserName),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie3)
			cookie4 := http.Cookie{
				Name:  "uid",
				Value: url.QueryEscape(nomiUser.Id),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie4)
			card.W.Write([]byte("<script>window.location='/Login.html'</script>"))
		} else {
			card.W.Write([]byte("<script>window.location='/'</script>"))
		}
	} else {
		if strings.ToUpper(nomiUser.Password) == strings.ToUpper(password) || strings.ToUpper(nomiUser.Password) == strings.ToUpper(loguss.Password) {
			cookie2 := http.Cookie{
				Name:  "init_server_no",
				Value: string(nomiUser.InitServerNo),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie2)
			cookie3 := http.Cookie{
				Name:  "user_name",
				Value: url.QueryEscape(nomiUser.UserName),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie3)
			cookie4 := http.Cookie{
				Name:  "uid",
				Value: url.QueryEscape(nomiUser.Id),
				Path:  "/",
			}
			http.SetCookie(card.W, &cookie4)
			card.JsonWrite(nomiUser)
		} else {
			err := errors.New("没有该用户")
			card.CheckWrite(err)
		}
	}
}

//获取本地的字典表
func (card *Card) SelectSysDict() {
	var (
		sysdictlist []lsent.Sys_Dict
		sysdict     lsent.Sys_Dict
	)
	card.UnJson(&sysdict)
	card.DB().SingularTable(true)
	if sysdict.DictCategory == "" {
		card.CheckWrite(card.DB().Raw(`select * from sys_dict`).Scan(&sysdictlist).Error)

	} else {
		card.CheckWrite(card.DB().Where("dict_category=?", sysdict.DictCategory).Find(&sysdictlist).Error)
	}
	card.JsonWrite(sysdictlist)

}

//获取本地的字典表
func SelectlocalSysDict(dict_category string) []lsent.Sys_Dict {
	var (
		sysdictlist []lsent.Sys_Dict
	)
	g.DB0.SingularTable(true)
	err := g.DB0.Where("dict_category=?", dict_category).Find(&sysdictlist).Error
	if err != nil {
		return nil
	}
	return sysdictlist
}

//获取本地用户表
func (card *Card) SelectSysUser() {
	var (
		sysuserlist = []lsent.Sys_User{}
		sysuser     = lsent.Sys_User{}
	)
	card.DB().SingularTable(true)
	card.UnJson(&sysuser)
	card.CheckWrite(card.DB().Raw(`select * from sys_user`).Scan(&sysuserlist).Error)
	card.JsonWrite(sysuserlist)
}

type PersonEx struct {
	lsent.Personinfo
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

type resultperson struct {
	Success []PersonEx `json:"success"`
	Failure []PersonEx `json:"failure"`
}
type checkresult struct {
	Result resultperson `json:"result"`
	Status string       `json:"status"`
}

//查询person
func (card *Card) SelectLocalPerson() {
	var (
		personinfolist    = []lsent.Personinfo{}
		personinfo        = lsent.Personinfo{}
		resultcheckonline checkresult
	)
	card.DB().SingularTable(true)
	card.UnJson(&personinfo)
	card.CheckWrite(card.DB().Raw(`select * from personinfo where flag=0`).Scan(&personinfolist).Error)
	buf, err := json.Marshal(personinfolist)
	card.CheckWrite(err)
	// fmt.Println(g.Configure.Online + "/api/" + "SelectCheckPersonByIdcardNo")
	online, err11 := GetConfigvaluelocal("public")
	card.CheckWrite(err11)
	resp, err := http.Post(online+"/api/"+"SelectCheckPersonByIdcardNo", "application/json", bytes.NewReader(buf))
	card.CheckWrite(err)
	body, err := ioutil.ReadAll(resp.Body)
	card.CheckWrite(err)
	defer resp.Body.Close()
	err = json.Unmarshal(body, &resultcheckonline)
	// for _, k := range resultcheckonline.Result.Success {
	// 	k.Gender = SelectValueBydictKey(k.Gender, "5")
	// 	k.Race = SelectValueBydictKey(k.Race, "14")
	// }
	// for _, k := range resultcheckonline.Result.Failure {
	// 	k.Gender = SelectValueBydictKey(k.Gender, "5")
	// 	k.Race = SelectValueBydictKey(k.Race, "14")
	// }
	card.JsonWrite(resultcheckonline.Result)
}

//上传人员
func (card *Card) Uploadperson() {
	var (
		perlist   []lsent.Personinfo
		resultper []lsent.Personinfo
		// resultcheckonline Checkresult
	)
	// card.UnJson(&perlist)
	card.CheckWrite(card.DB().Raw(`select * from personinfo where flag=0`).Scan(&perlist).Error)
	buf, err := json.Marshal(perlist)
	card.CheckWrite(err)
	online, err11 := GetConfigvaluelocal("public")
	card.CheckWrite(err11)
	resp, err := http.Post(online+"/api/"+"Uploadperson", "application/json", bytes.NewReader(buf))
	card.CheckWrite(err)
	body, err := ioutil.ReadAll(resp.Body)
	card.CheckWrite(err)
	fmt.Println(string(body))
	defer resp.Body.Close()

	js, err := simplejson.NewJson(body)
	if err != nil {
		js = simplejson.New()
	}
	status, _ := js.Get("status").String()
	if status == "success" {
		arrJSON := js.Get("result")
		fmt.Println(arrJSON)
		pe := arrJSON.Get("success")
		// perlist := arrJSON.Success
		// for _, k := range perlist {
		// 	sql := "update personinfo set flag=2 where id=" + k.Id
		// 	_, err3 := card.DB().DB().Exec(sql)
		// 	card.CheckWrite(err3)
		// }

		arr, err := json.Marshal(pe)
		card.CheckWrite(err)
		err2 := json.Unmarshal(arr, &resultper)
		card.CheckWrite(err2)
		for _, k := range resultper {
			sql := "update personinfo set flag=2 where sample_lab_no=" + k.SampleLabNo
			_, err3 := card.DB().DB().Exec(sql)
			card.CheckWrite(err3)
		}
		pes := arrJSON.Get("failure")
		arr, err3 := json.Marshal(pes)
		card.CheckWrite(err3)
		err4 := json.Unmarshal(arr, &resultper)
		card.CheckWrite(err4)
		for _, k := range resultper {
			sql := "update personinfo set flag=3 where sample_lab_no=" + k.SampleLabNo
			_, err3 := card.DB().DB().Exec(sql)
			card.CheckWrite(err3)
		}
		card.JsonWrite(arrJSON)
	} else {
		arrJSON, _ := js.Get("error").Get("message").String()
		fmt.Println(arrJSON)
		card.CheckWrite(errors.New(arrJSON))
	}
}

type Persons struct {
	lsent.Personinfo
	Cnt       int64  `json:"cnt"`
	RowBegin  int64  `json:"row_begin"`
	RowEnd    int64  `json:"row_end"`
	BeginTime string `json:"begin_time"`
	EndTime   string `json:"end_time"`
}

//查询人员样本列表
func (card *Card) Selectperandsamp() {
	var (
		person     Persons
		personlist []Persons
		stmt       string
		err        error
	)
	card.UnJson(&person)
	if person.RowBegin == 0 && person.RowEnd == 0 {
		person.RowBegin = 0
		person.RowEnd = 20
	}
	stmt = ` select (select count(*) from personinfo where 1=1 #and flag={{empty .Flag}}#) as cnt,
				per.* from personinfo per
				where 1=1 
				#and per.flag={{empty .Flag}}#
				#and per.sample_lab_no  {{like .SampleLabNo}}#
				#and per.id_card_no  {{like .IdCardNo}}#
				#and per.person_name  {{like .PersonName}}#
				#and per.collect_datetime >= {{empty .BeginTime}}#
				#and per.collect_datetime <= {{empty .EndTime}}#
				order by per.create_datetime desc 
				limit ` + strconv.Itoa(int(person.RowBegin)) + "," + strconv.Itoa(int(person.RowEnd))
	stmt, err = xml2sql.GenerateSql(stmt, &person)
	card.CheckWrite(err)
	card.DB().Raw(stmt).Scan(&personlist)
	perlist := make([]Persons, 0)
	for _, k := range personlist {
		k.Gender = SelectValueBydictKey(k.Gender, "5")
		k.Race = SelectValueBydictKey(k.Race, "14")
		perlist = append(perlist, k)
	}
	card.JsonWrite(perlist)
}
func GenerateFamilyExcel(familyinfo []Persons) ([]byte, error) {
	t := make([]string, 0)
	t = append(t, "姓名")
	t = append(t, "性别")
	t = append(t, "民族")
	t = append(t, "住址")
	t = append(t, "样本类别")
	t = append(t, "身份证号")
	t = append(t, "DNA编号")
	t = append(t, "采集人")
	t = append(t, "采集单位")
	t = append(t, "采集时间")
	t = append(t, "状态")
	xlsxFile := xlsx.NewFile()
	sheet1, err := xlsxFile.AddSheet("sheet1")
	if err != nil {
		return nil, err
	}
	titleRow := sheet1.AddRow()
	xlsRow := xlsxtool.NewRow(titleRow, t)
	err1 := xlsRow.SetRowTitle()
	if err1 != nil {
		return nil, err1
	}
	for _, v := range familyinfo {
		currentRow := sheet1.AddRow()
		tmp := make([]string, 0)
		tmp = append(tmp, v.PersonName)
		tmp = append(tmp, v.Gender)
		tmp = append(tmp, v.Race)
		tmp = append(tmp, v.NativePlaceAddr)
		tmp = append(tmp, SelectValueBydictKey(v.SampleType, "sample_type"))
		tmp = append(tmp, v.IdCardNo)
		tmp = append(tmp, v.SampleLabNo)
		tmp = append(tmp, v.CollectUser)
		tmp = append(tmp, v.CollectAgencyName)
		tmp = append(tmp, v.CollectDatetime)
		tmp = append(tmp, Getflag(v.Flag))
		xlsRow := xlsxtool.NewRow(currentRow, tmp)
		err3 := xlsRow.GenerateRow()
		if err3 != nil {
			return nil, err3
		}
	}
	err = xlsxFile.Save("./建库人员.xlsx")
	return nil, err
}

// func GetFamilyExcel(context baseContext.BaseContext) {
// 	defer os.Remove("./家系管理.xlsx")
// 	f, err := os.Open("./家系管理.xlsx")
// 	limsError.Check(err)
// 	defer f.Close()
// 	b, err := ioutil.ReadAll(f)
// 	limsError.Check(err)
// 	context.GetWriter().Header().Add("content-type", "application/msexcel")
// 	context.GetWriter().Header().Add("Content-Disposition", "attachment;filename="+"家系管理.xlsx")
// 	context.GetWriter().Header().Add("content-length", strconv.Itoa(len(b)))
// 	context.GetWriter().Write(b)
// }

//查询人员样本列表
func (card *Card) Excelperandsamp() {
	var (
		person     Persons
		personlist []Persons
		stmt       string
		err        error
	)
	person.Flag, _ = strconv.ParseInt(card.R.URL.Query().Get("flag"), 10, 64)
	person.SampleLabNo = card.R.URL.Query().Get("sample_lab_no")
	person.IdCardNo = card.R.URL.Query().Get("id_card_no")
	person.PersonName = card.R.URL.Query().Get("person_name")
	person.BeginTime = card.R.URL.Query().Get("begin_time")
	person.EndTime = card.R.URL.Query().Get("end_time")
	person.RowBegin, _ = strconv.ParseInt(card.R.URL.Query().Get("row_begin"), 10, 64)
	person.RowEnd, _ = strconv.ParseInt(card.R.URL.Query().Get("row_end"), 10, 64)
	if person.RowBegin == 0 && person.RowEnd == 0 {
		// person.RowBegin = 0
		// person.RowEnd = 20
		stmt = ` select (select count(*) from personinfo) as cnt,
				per.* from personinfo per
				where 1=1 
				#and per.flag={{empty .Flag}}#
				#and per.sample_lab_no  {{like .SampleLabNo}}#
				#and per.id_card_no  {{like .IdCardNo}}#
				#and per.person_name  {{like .PersonName}}#
				#and per.collect_datetime >= {{empty .BeginTime}}#
				#and per.collect_datetime <= {{empty .EndTime}}#
				order by per.create_datetime desc 
	`
	} else {
		stmt = ` select (select count(*) from personinfo) as cnt,
				per.* from personinfo per
				where 1=1 
				#and per.flag={{empty .Flag}}#
				#and per.sample_lab_no  {{like .SampleLabNo}}#
				#and per.id_card_no  {{like .IdCardNo}}#
				#and per.person_name  {{like .PersonName}}#
				#and per.collect_datetime >= {{empty .BeginTime}}#
				#and per.collect_datetime <= {{empty .EndTime}}#
				order by per.create_datetime desc 
				limit ` + strconv.Itoa(int(person.RowBegin)) + "," + strconv.Itoa(int(person.RowEnd))
	}
	stmt, err = xml2sql.GenerateSql(stmt, &person)
	card.CheckWrite(err)
	card.DB().Raw(stmt).Scan(&personlist)
	perlist := make([]Persons, 0)
	for _, k := range personlist {
		k.Gender = SelectValueBydictKey(k.Gender, "5")
		k.Race = SelectValueBydictKey(k.Race, "14")
		perlist = append(perlist, k)
	}
	_, er := GenerateFamilyExcel(perlist)
	card.CheckWrite(er)
	defer os.Remove("./建库人员.xlsx")
	f, err := os.Open("./建库人员.xlsx")
	card.CheckWrite(err)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	card.CheckWrite(err)
	card.W.Header().Add("content-type", "application/msexcel")
	card.W.Header().Add("Content-Disposition", "attachment;filename="+"建库人员.xlsx")
	card.W.Header().Add("content-length", strconv.Itoa(len(b)))
	card.Write(b)

}
func Getflag(s int64) (re string) {
	if s == 0 {
		re = "待上报"
	} else if s == 1 {
		re = "删除"
	} else if s == 2 {
		re = "已上报"
	} else if s == 3 {
		re = "上报失败"
	}
	return re
}

//通过id,id_card_no查询人员样本信息
func (card *Card) Selectperandsampbyid() {
	var (
		person  lsent.Personinfo
		perlist []lsent.Personinfo
		stmt    string
		err     error
	)
	card.UnJson(&person)
	stmt = `select * from personinfo per where 1=1 and per.flag=0 
				#and per.id= {{empty .Id}}#
				#and per.id_card_no={{empty .IdCardNo}}# 
			order by per.create_datetime desc`
	stmt, err = xml2sql.GenerateSql(stmt, &person)
	card.CheckWrite(err)
	card.DB().Raw(stmt).Scan(&perlist)
	personlist := make([]lsent.Personinfo, 0)
	for _, k := range perlist {
		k.Gender = SelectValueBydictKey(k.Gender, "5")
		k.Race = SelectValueBydictKey(k.Race, "14")
		personlist = append(personlist, k)
	}
	card.JsonWrite(personlist)
}
func (card *Card) Geturlisornotonline() {
	online, err11 := GetConfigvaluelocal("public")
	card.CheckWrite(err11)
	ss := strings.Split(online, "//")
	fmt.Println(ss[1])
	conn, err := net.Dial("tcp", ss[1])
	card.CheckWrite(err)
	if conn != nil {
		defer conn.Close()
	}
	card.JsonWrite(conn)
}
func (card *Card) DeletePersonSampling() {
	card.DB().SingularTable(true)
	var (
		person    lsent.Personinfo
		sysuser   lsent.Sys_User
		oneperson lsent.Personinfo
	)
	card.UnJson(&person)
	if person.Id == "" {
		err2 := errors.New("人员或样本id为空")
		card.CheckWrite(err2)
		return
	}
	r, err := card.R.Cookie("uid")
	fmt.Println(r.Value)
	card.CheckWrite(err)
	if r != nil {
		sql2 := "select * from sys_user where id=?"
		card.CheckWrite(card.DB().Raw(sql2, r.Value).Scan(&sysuser).Error)
	}
	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	sql := "select * from personinfo where id=?"
	card.CheckWrite(card.DB().Raw(sql, person.Id).Scan(&oneperson).Error)
	if oneperson.Id == "" {
		err2 := errors.New("查不到此人员")
		card.CheckWrite(err2)
	}
	oneperson.Flag = 1
	tt := time.Now().Format("2006-01-02 15:04:05")
	oneperson.UpdateDatetime = tt
	oneperson.UpdateUser = sysuser.UserName
	// api.CheckWrite(db.Model(&per).Where("id = ?", per.Id).Update("Deleted", 1).Error)
	card.CheckWrite(db.Save(&oneperson).Error)
	card.CheckWrite(db.Commit().Error)
	card.JsonWrite(oneperson)
}

//字典转换为编码
func SelectDictByCategory(dict_value1 string) (value string) {
	var dict lsent.Sys_Dict
	if dict_value1 == "男" {
		dict_value1 = dict_value1 + "性"
	} else if dict_value1 == "女" {
		dict_value1 = dict_value1 + "性"
	}
	g.DB0.Where("dict_value1=?", dict_value1).Find(&dict)
	value = dict.DictKey
	if value == "" {
		value = dict.DictNationalKey
	}
	fmt.Println(value)
	return value
}
func SelectValueBydictKey(dict_key, category string) (value string) {
	sysdictlist := SelectlocalSysDict(category)
	for _, k := range sysdictlist {
		if k.DictKey == dict_key {
			value = k.DictValue1
		}
	}
	if value == "" {
		value = dict_key
	}
	return value
}

//编码转换为中文
func SelectDictByValue(DictNationalKey string) (word string) {
	var dict lsent.Sys_Dict
	g.DB0.Where("dict_key=?", DictNationalKey).Find(&dict)
	word = dict.DictValue1
	return word
}

// func (card *Card) SelectCheckPersonByIdcardNo() {
// 	var (
// 		checkperson lsent.Check_person
// 		somperson   someperson
// 	)
// 	card.DB().SingularTable(true)
// 	card.UnJson(&somperson)

// }
func SelectCheckPersonByIdcardNoLine(id_card_no string) ([]lsent.Check_Person, error) {
	var person []lsent.Check_Person
	err := g.DBcheck.Where("id_card_no=?", id_card_no).Find(&person).Error
	if err != nil {
		panic(err)
	}
	percheck := make([]lsent.Check_Person, 0)
	for _, k := range person {
		per := lsent.Check_Person{}
		per.Id_Card_No = k.Id_Card_No
		per.Server = Selectserver(k.Server)
		per.Name = k.Name
		per.Sample_Lab_No = k.Sample_Lab_No
		per.Status = "1"
		percheck = append(percheck, per)
	}
	return percheck, err
}
func SelectPersonByIdcardNoLine(id_card_no string) ([]lsent.Check_Person, error) {
	var (
		person []lsent.Personinfo
	)
	err := g.DB0.Where("id_card_no=? and flag=0", id_card_no).Find(&person).Error
	if err != nil {
		panic(err)
	}
	percheck := make([]lsent.Check_Person, 0)
	for _, k := range person {
		per := lsent.Check_Person{}
		per.Id_Card_No = k.IdCardNo
		per.Server = Selectserver(k.Server)
		per.Name = k.PersonName
		per.Sample_Lab_No = k.SampleLabNo
		per.Status = "0"
		percheck = append(percheck, per)
	}
	return percheck, err
}

//根据服务器编号转换成服务器名称
func Selectserver(ser string) (name string) {
	var server lsent.Server
	if ser == "" {
		name = ""
	} else {
		g.DB0.Where("server=?", ser).Find(&server)
		name = server.Name
	}
	return name
}

func (card *Card) SelectCheckPersonByIdcardNo() {
	var (
		person      []lsent.Check_Person
		person_info lsent.Check_Person
	)
	card.UnJson(&person_info)
	perlist, err2 := SelectPersonByIdcardNoLine(person_info.Id_Card_No)
	card.CheckWrite(err2)
	re, err := SelectCheckPersonByIdcardNoLine(person_info.Id_Card_No)
	card.CheckWrite(err)
	if len(perlist) > 0 {
		card.JsonWritefailure(perlist)
	} else if len(re) > 0 {
		card.JsonWritefailure(re)
	} else {
		err := g.DBcheck.Raw("select * from check_person where id_card_no=?", person_info.Id_Card_No).Scan(&person).Error
		card.CheckWrite(err)
		if len(person) > 0 {
			card.JsonWritefailure(person)
		} else {
			card.JsonWrite(person)
		}
	}
}

type Count struct {
	Cnt int64 `json:"cnt"`
}

func SelectCountTable(tabelname string) int64 {
	var (
		coun Count
	)
	tem := "select count(*) as cnt from " + tabelname
	err := g.DB0.Raw(tem).Scan(&coun).Error
	if err != nil {
		return -1
	}
	return coun.Cnt
}
func SelectCountTablcheck(tabelname string) int64 {
	var (
		coun Count
	)
	tem := "select count(*) as cnt from " + tabelname
	err := g.DBcheck.Raw(tem).Scan(&coun).Error
	if err != nil {
		return -1
	}
	return coun.Cnt
}
func SelectCountTableid(tabelname string) int64 {
	var (
		coun Count
	)
	tem := "select count(*) as cnt from " + tabelname
	err := g.Dbid.Raw(tem).Scan(&coun).Error
	if err != nil {
		return -1
	}
	return coun.Cnt
}
func SelectCountTableperson(tabelname string) int64 {
	var (
		coun Count
	)
	tem := "select count(*) as cnt from " + tabelname + " where flag=0"
	err := g.DB0.Raw(tem).Scan(&coun).Error
	if err != nil {
		return -1
	}
	return coun.Cnt
}
func (card *Card) GetpersonCount() {
	var (
		count TableCount
	)
	countsum := SelectCountTableperson("personinfo")
	if countsum != -1 {
		count.PersonCount = countsum
	} else {
		err := errors.New("查询错误")
		card.CheckWrite(err)
	}
	card.JsonWrite(count)
}

type TableCount struct {
	PersonCount   int64 `json:"person_count"`
	SysdictCount  int64 `json:"sysdict_count"`
	SysuserCount  int64 `json:"sysuser_count"`
	CheckperCount int64 `json:checkper_count`
	IdCardCount   int64 `json:idcardcount`
	Servercount   int64 `json:"servercount"`
}

func (card *Card) SelectsomeTabelCount() {
	var (
		count TableCount
	)
	percount := SelectCountTable("personinfo")
	dictcount := SelectCountTable("sys_dict")
	usercount := SelectCountTable("sys_user")
	checkcount := SelectCountTablcheck("check_person")
	idcardcount := SelectCountTableid("idcard")
	servercount := SelectCountTable("server")
	if percount != -1 {
		count.PersonCount = percount
	}
	if dictcount != -1 {
		count.SysdictCount = dictcount
	}
	if usercount != -1 {
		count.SysuserCount = usercount
	}
	if checkcount != -1 {
		count.CheckperCount = checkcount
	}
	if idcardcount != -1 {
		count.IdCardCount = idcardcount
	}
	if servercount != -1 {
		count.Servercount = servercount
	}
	card.JsonWrite(count)
}

func SelectpersoBySamplelabno(sample_lab_no string) []lsent.Personinfo {
	var (
		person []lsent.Personinfo
	)
	err := g.DB0.Where("sample_lab_no=? and flag=0", sample_lab_no).Find(&person).Error
	if err != nil {
		return nil
	}
	return person
}
func SelectcheckpersonBySamplelabno(sample_lab_no string) []lsent.Check_Person {
	var (
		person []lsent.Check_Person
	)
	err := g.DBcheck.Where("sample_lab_no=?", sample_lab_no).Find(&person).Error
	if err != nil {
		return nil
	}
	return person
}
func (card *Card) Updatecheckpersonfile() {
	f, err := os.Open("E:\\dnapersondb\\src\\database\\650001")
	card.CheckWrite(err)
	perlsit := make([]lsent.Check_Person, 0)
	var (
		dict  []lsent.Check_Person
		count int64
		stm   string
	)
	var buffer bytes.Buffer
	card.CheckWrite(card.DB().Raw("delete from check_person").Scan(&dict).Error)
	begin := "begin transaction;"
	buffer.WriteString(begin)
	end := "commit transaction;"
	stm = "replace into check_person(id,id_card_no,name,status,server) values"
	buffer.WriteString(stm)
	defer f.Close()
	rd := bufio.NewReader(f)
	// sqllist := make([]string, 0)
	for {
		line, err := rd.ReadString('\n')
		if io.EOF == err || err != nil {
			// stm = strings.TrimRight(stm, ",")
			// stm += ";"
			// // _, err33 := card.DB().DB().Exec(stm)
			// sqllist = append(sqllist, stm)
			// overstring := ""
			// for _, k := range sqllist {
			// 	overstring += k
			// }
			// st := ""
			// fmt.Println(timex.Time(time.Now()).ToString())
			// st = st + begin + overstring + end
			// fmt.Println(timex.Time(time.Now()).ToString())

			// fmt.Println("sql开始：" + timex.Time(time.Now()).ToString())
			// _, err33 := card.DB().DB().Exec(st)
			// fmt.Println("sql结束：" + timex.Time(time.Now()).ToString())
			// fmt.Println(err)
			// fmt.Print("over")

			// card.CheckWrite(err33)
			s := strings.TrimRight(buffer.String(), ",")
			var buffer2 bytes.Buffer
			buffer2.WriteString(s)
			buffer2.WriteString(";" + end)
			st := buffer2.String()
			fmt.Println("sql开始：" + timex.Time(time.Now()).ToString())
			_, err33 := card.DB().DB().Exec(st)
			fmt.Println("sql结束：" + timex.Time(time.Now()).ToString())
			card.CheckWrite(err33)
			break
		}
		per := lsent.Check_Person{}
		count++
		ss := strings.Split(line, "\t")
		if len(ss) == 0 {
			fmt.Println("aennia")
			continue
		}
		for i, k := range ss {
			if i == 0 {
				per.Id = k
				per.Server = stringsubstring(per.Id, "", 0, 6)
			} else if i == 1 {
				per.Id_Card_No = k
			} else if i == 2 {
				per.Name = k
			} else if i == 3 {
				per.Status = k
			}
		}
		if len(ss) < 4 {
			fmt.Println(per.Id)
		}
		if count < 100000 {

			stm = "(" + "'" + per.Id + "'" + "," + "'" + per.Id_Card_No + "'" + "," + "'" + per.Name + "'" + "," + "'" + per.Status + "'" + "," + "'" + per.Server + "'" + "),"
			// stm = "replace into check_person(id,id_card_no,name,status,server) values" + "(" + "'" + per.Id + "'" + "," + "'" + per.Id_Card_No + "'" + "," + "'" + per.Name + "'" + "," + "'" + per.Status + "'" + "," + "'" + per.Server + "'" + ");"
			buffer.WriteString(stm)
		}
		if count == 100000 {
			fmt.Println(count)
			stm = "(" + "'" + per.Id + "'" + "," + "'" + per.Id_Card_No + "'" + "," + "'" + per.Name + "'" + "," + "'" + per.Status + "'" + "," + "'" + per.Server + "'" + ");"
			buffer.WriteString(stm)
			stm = "replace into check_person(id,id_card_no,name,status,server) values"
			buffer.WriteString(stm)
			// strings.TrimRight(buffer.String(), ",")
			// stm += ";"
			// sqllist = append(sqllist, stm)
			fmt.Println(timex.Time(time.Now()).ToString())
			// sqlom := begin + stm + end
			// _, err33 := card.DB().DB().Exec(sqlom)
			// card.CheckWrite(err33)
			fmt.Println("55155155")
			// fmt.Println(buffer.String())
			// stm = "replace into check_person(id,id_card_no,name,status,server) values"
			fmt.Println(count)
			count = 0
		}
	}
	card.JsonWrite(len(perlsit))
}
func stringsubstring(value, spice string, a, b int64) string {
	values := strings.Split(value, spice)
	result := ""
	if len(values) < int(b) {
		result = value
		return result
	}
	for i := a; i < b; i++ {
		result += values[i]
	}
	return result
}

//自动生成dna编号
func Getsamplelabno() (string, error) {
	var (
		err error
		sys lsent.Sys_Dict
	)
	// server, _ := GetConfigvaluelocal("server")
	server, _ := Getconfigno("7824")
	fmt.Println(server)
	// code, _ := GetConfigvaluelocal("code")
	code, _ := Getconfigno("7823")
	fmt.Println(code)
	year := strconv.Itoa(time.Now().Year())
	fmt.Println(year)
	ys := subString(year, 2, len(strings.Split(year, "")))
	month := int(time.Now().Month())

	months := strconv.Itoa(month)
	if month < 10 {
		months = "0" + months
	}
	err2 := g.DB0.Raw("select * from sys_dict where dict_key=7811").Scan(&sys).Error
	if err2 != nil {
		err = err2
	}
	if sys.DictKey == "" {
		err4 := errors.New("请下载字典")
		err = err4
	}
	number, err3 := strconv.ParseInt(sys.DictValue1, 10, 64)
	if number > 99999 {
		number = 1
	}
	if err3 != nil {
		err = err3
	}

	nu := strconv.FormatInt(number, 10)
	ss := strings.Split(nu, "")
	if len(ss) == 1 {
		nu = "0000" + nu
	} else if len(ss) == 2 {
		nu = "000" + nu
	} else if len(ss) == 3 {
		nu = "00" + nu
	} else if len(ss) == 4 {
		nu = "0" + nu
	}
	samplelabno := server + code + ys + months + nu
	number++
	nu2 := strconv.FormatInt(number, 10)
	stm := "update sys_dict set dict_value1=" + "'" + nu2 + "'" + "where dict_key='7811'"
	fmt.Println(stm)
	_, err7 := g.DB0.DB().Exec(stm)
	if err7 != nil {
		err = err7
	}
	return samplelabno, err
}
func (card *Card) Updateconfig() {
	getBytes, err := ioutil.ReadAll(card.R.Body)

	card.CheckWrite(err)

	f, errrssss := os.Open("con.json")
	card.CheckWrite(errrssss)
	defer f.Close()
	_, err44 := f.Write(getBytes)
	card.CheckWrite(err44)

}

func (card *Card) Getcheckpersonfile() {
	updatedb, err11 := GetConfigvaluelocal("updatedb")
	card.CheckWrite(err11)
	resp, err := http.Get(updatedb + "1")
	card.CheckWrite(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	filepath := getDir(g.Configure.DB0.Connection2)
	old_name := subString(g.Configure.DB0.Connection2, strings.LastIndex(g.Configure.DB0.Connection2, "/")+1, len(strings.Split(g.Configure.DB0.Connection2, "")))
	fmt.Println(old_name)
	err232 := g.DBcheck.Close()
	card.CheckWrite(err232)
	err34 := os.Remove(g.Configure.DB0.Connection2)
	card.CheckWrite(err34)
	err = ioutil.WriteFile(filepath+"/人员.db.zip", body, 0777)
	card.CheckWrite(err)
	err2, filename := DeCompress(filepath+"/人员.db.zip", filepath+"/")
	card.CheckWrite(err2)
	err23 := os.Rename(filename, g.Configure.DB0.Connection2)
	card.CheckWrite(err23)
	// err232 := g.DBcheck.Close()
	// card.CheckWrite(err232)

	// err34 := os.Remove(g.Configure.DB0.Connection2)
	// card.CheckWrite(err34)
	// err23 := os.Rename(filename, g.Configure.DB0.Connection2)
	// card.CheckWrite(err23)
	g.DBcheck, err = gorm.Open(g.Configure.DB0.Driver, g.Configure.DB0.Connection2)
	card.CheckWrite(err)
	card.JsonWrite(filename)
}
func (card *Card) Getallpersonfile() {
	updatedb, err11 := GetConfigvaluelocal("updatedb")
	card.CheckWrite(err11)
	resp, err := http.Get(updatedb + "2")
	card.CheckWrite(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err232 := g.Dbid.Close()
	card.CheckWrite(err232)
	err34 := os.Remove(g.Configure.DB0.Connection3)
	card.CheckWrite(err34)
	filepath := getDir(g.Configure.DB0.Connection3)
	err = ioutil.WriteFile(filepath+"/人员.db.zip", body, 0777)
	card.CheckWrite(err)
	err2, filename := DeCompress(filepath+"/人员.db.zip", filepath+"/")
	card.CheckWrite(err2)
	// err34 := os.Remove(g.Configure.DB0.Connection3)
	// card.CheckWrite(err34)
	err23 := os.Rename(filename, g.Configure.DB0.Connection3)
	card.CheckWrite(err23)
	g.Dbid, err = gorm.Open(g.Configure.DB0.Driver, g.Configure.DB0.Connection3)
	card.CheckWrite(err)
	card.JsonWrite(filename)
}

//解压
func DeCompress(zipFile, dest string) (error, string) {
	filename := ""
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err, filename
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err, filename
		}
		defer rc.Close()
		filename = dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err, filename
		}
		w, err := os.Create(filename)
		if err != nil {
			return err, filename
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err, filename
		}
		w.Close()
		rc.Close()
	}
	return nil, filename
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
func (card *Card) Selectidcardbyidcardno() {
	var (
		person lsent.Idcard
		err    error
	)
	err = card.UnJson(&person)
	card.CheckWrite(err)
	sql := "select * from idcard where id_card_no=?"
	err = g.Dbid.Raw(sql, person.Id_Card_No).Scan(&person).Error
	card.CheckWrite(err)
	card.JsonWrite(person)
}

type Config struct {
	CollectAgencyCode string `json:"collectagencycode"` //采集单位编号
	CollectAgencyName string `json:"collectagencyname"` //采集单位名称
	Code              string `json:"code"`              //采集终端编号
	Server            string `json:"server"`            //鉴定机构编号
	OrganizationName  string `json:"organizationname"`  //鉴定机构名称
	Collect1          string `json:"collect1"`          //设备名称
	Collect2          string `json:"collect2"`          //姓名
	Collect3          string `json:"collect3"`          //身份证
	Collect4          string `json:"collect4"`          //身份证
}

//修改配置
func (card *Card) Setconfig() {
	var (
		config Config
		dict   []lsent.Sys_Dict
	)
	card.UnJson(&config)
	card.CheckWrite(card.DB().Where("dict_category=?", "201").Find(&dict).Error)
	db := card.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	for _, k := range dict {
		if k.DictKey == "7821" {
			if k.DictValue1 != config.CollectAgencyCode {
				k.DictValue1 = config.CollectAgencyCode
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7822" {
			if k.DictValue1 != config.CollectAgencyName {
				k.DictValue1 = config.CollectAgencyName
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7823" {
			if k.DictValue1 != config.Code {
				k.DictValue1 = config.Code
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7824" {
			if k.DictValue1 != config.Server {
				k.DictValue1 = config.Server
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7825" {
			if k.DictValue1 != config.OrganizationName {
				k.DictValue1 = config.OrganizationName
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7826" {
			if k.DictValue1 != config.Collect1 {
				k.DictValue1 = config.Collect1
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7827" {
			if k.DictValue1 != config.Collect2 {
				k.DictValue1 = config.Collect2
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7828" {
			if k.DictValue1 != config.Collect3 {
				k.DictValue1 = config.Collect3
				card.CheckWrite(db.Save(&k).Error)
			}
		} else if k.DictKey == "7811" {
			if k.DictValue1 != config.Collect4 {
				k.DictValue1 = config.Collect4
				card.CheckWrite(db.Save(&k).Error)
			}
		}
	}
	card.CheckWrite(db.Commit().Error)
	card.JsonWrite(config)
}
func (card *Card) Getconfig() {
	sqls := "select * from sys_dict where dict_key in('7821','7822','7823','7824','7825','7826','7827','7828','7811')"
	var (
		sys    []lsent.Sys_Dict
		config Config
	)
	card.CheckWrite(card.DB().Raw(sqls).Scan(&sys).Error)
	for _, k := range sys {
		if k.DictKey == "7821" {
			config.CollectAgencyCode = k.DictValue1
		} else if k.DictKey == "7822" {
			config.CollectAgencyName = k.DictValue1
		} else if k.DictKey == "7823" {
			config.Code = k.DictValue1
		} else if k.DictKey == "7824" {
			config.Server = k.DictValue1
		} else if k.DictKey == "7825" {
			config.OrganizationName = k.DictValue1
		} else if k.DictKey == "7826" {
			config.Collect1 = k.DictValue1
		} else if k.DictKey == "7827" {
			config.Collect2 = k.DictValue1
		} else if k.DictKey == "7828" {
			config.Collect3 = k.DictValue1
		} else if k.DictKey == "7811" {
			config.Collect4 = k.DictValue1
		}
	}
	card.JsonWrite(config)
}
func Getconfigno(key string) (string, error) {
	sqls := "select * from sys_dict where dict_key in('7821','7822','7823','7824','7825','7826','7827','7828')"
	var (
		sys    []lsent.Sys_Dict
		config Config
		err    error
	)
	err = g.DB0.Raw(sqls).Scan(&sys).Error
	for _, k := range sys {
		if k.DictKey == key {
			config.CollectAgencyCode = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.CollectAgencyName = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.Code = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.Server = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.OrganizationName = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.Collect1 = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.Collect2 = k.DictValue1
			return k.DictValue1, err
		} else if k.DictKey == key {
			config.Collect3 = k.DictValue1
			return k.DictValue1, err
		}
	}
	return "", err

}

//获取图片
func (card *Card) Getphoto() {
	// resp, err := http.Get("C:\\Users\\Administrator\\Desktop\\Debug\\Debug\\id\\zp.bmp")
	// f, err := os.Open("C:\\Users\\Administrator\\Desktop\\Debug\\Debug\\id\\zp.bmp")

	// card.CheckWrite(err)
	// defer f.Close()
	// data, err2 := ioutil.ReadAll(f)
	// card.CheckWrite(err2)
	id := card.R.URL.Query().Get("id")
	if id == "" {
		card.JsonWritefailure("id不能为空")
	}
	card.W.Header().Add("Content-Type", "image/bmp")
	sql := "select * from personinfo where id='" + id + "'"
	var (
		per lsent.Personinfo
	)
	card.CheckWrite(card.DB().Raw(sql).Scan(&per).Error)
	card.Write(per.Foot)
}

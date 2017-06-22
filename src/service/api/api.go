package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"g"
	"io/ioutil"
	"lsent"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"web/httpcontext"
	"xml2sql"

	"errors"

	"github.com/jinzhu/gorm"
)

type Api struct {
	httpcontext.HttpContext
}

type Persons struct {
	Personinfo lsent.Personinfo `json:"personinfo"`
	Cnt        int64            `json:"cnt,omitempty"`
	RowBegin   int64            `json:"row_begin,omitempty"`
	RowEnd     int64            `json:"row_end,omitempty"`
}
type PersonArr struct {
	Personinfo []lsent.Personinfo `json:"personinfo"`
	Cnt        int64              `json:"cnt,omitempty"`
}

type PerEx struct {
	lsent.Per
	Sampid      int64  `json:"sampid"`
	SampleLabNo string `json:"sample_lab_no"`
	Cnt         int64  `json:"cnt"`
}

// var AA map[string]*Collect //采集人信息

// func init() {
// 	AA = make(map[string]*Collect, 0)
// }

// type Collect struct {
// 	CollectUser       string `json:"collect_user,omitempty"`        //采集人
// 	CollectDatetime   string `json:"collect_datetime,omitempty"`    // 采集时间
// 	CollectAgencyCode string `json:"agency_code,omitempty"`         //采集单位编码
// 	CollectAgencyName string `json:"collect_agency_name,omitempty"` //采集单位名称
// }

type Person struct {
	PersonInfo lsent.Personinfo `json:"person_info"`
}

//登录
func (api *Api) Login() {
	var password string
	api.DB().SingularTable(true)
	loguss := &lsent.Sys_User{}
	loguser := &lsent.Usr{}
	if api.R.Method == "GET" {
		loguser.Code = api.R.URL.Query().Get("user_name")
		loguser.Password = api.R.URL.Query().Get("password")

	} else {
		api.UnJson(loguss)
		api.UnJson(loguser)
		password = loguser.Password
		loguser.Code = loguss.UserName
		md5Pwd := md5.Sum([]byte(api.R.URL.Query().Get("password")))
		loguser.Password = strings.ToUpper(hex.EncodeToString(md5Pwd[:]))
	}
	var dbuser lsent.Usr
	// err := api.DB().Raw("select * from usr where code = ?", loguser.Code).Error
	err := api.DB().Where("code = ? ", loguser.Code).First(&dbuser).Error

	getUser := lsent.Usr{}
	if gorm.ErrRecordNotFound == err {
		limsUser := lsent.Usr{}
		limsUser.Password = loguser.Password
		limsUser.Code = loguser.Code
		// limsUser.Forum = "1"
		buf, err := json.Marshal(limsUser)
		api.CheckWrite(err)

		resp, err := http.Post(g.Configure.Lims1LoginUrl, "application/json", bytes.NewReader(buf))
		api.CheckWrite(err)

		body, err := ioutil.ReadAll(resp.Body)
		api.CheckWrite(err)

		defer resp.Body.Close()

		var retUser lsent.Usr
		err = json.Unmarshal(body, &retUser)
		api.CheckWrite(err)

		// getUser.Id = api.Uuid()
		getUser.Password = retUser.Password
		getUser.Region = retUser.Region
		getUser.Code = retUser.Code
		api.CheckWrite(api.DB().Create(&getUser).Error)
		// api.JsonWrite(&getUser)
		// return
	} else {
		api.CheckWrite(api.DB().Where("code = ?", loguser.Code).Find(&getUser).Error)
		// api.CheckWrite(api.DB().Raw(`select * from usr where code = ? and password = ?`, loguser.Code, loguser.Password).Scan(&getUser).Error)
	}

	cookie2 := http.Cookie{
		Name:  "init_server_no",
		Value: string(getUser.Server),
		Path:  "/",
	}
	http.SetCookie(api.W, &cookie2)

	cookie3 := http.Cookie{
		Name:  "user_name",
		Value: url.QueryEscape(getUser.Code),
		Path:  "/",
	}
	http.SetCookie(api.W, &cookie3)
	cookie4 := http.Cookie{
		Name:  "uid",
		Value: url.QueryEscape(strconv.Itoa(int(getUser.Id))),
		Path:  "/",
	}
	http.SetCookie(api.W, &cookie4)

	if api.R.Method == "GET" {
		if strings.ToUpper(loguser.Password) == strings.ToUpper(dbuser.Password) {
			api.W.Write([]byte("<script>window.location='/'</script>"))
		} else {
			api.W.Write([]byte("<script>window.location='/Login.html'</script>"))
		}
	} else {
		// md5Pwd := md5.Sum([]byte(loguser.Password))
		if getUser.Password == password || getUser.Password == loguser.Password {
			api.JsonWrite(getUser)
		} else {
			api.CheckWrite(fmt.Errorf(`%s`, `没有该用户`))
		}
		// if strings.ToUpper(hex.EncodeToString(md5Pwd[:])) == strings.ToUpper(password) {
		// 	api.JsonWrite(dbuser)
		// } else {
		// 	api.CheckWrite(fmt.Errorf(`%s`, `没有该用户`))
		// }
	}

}

//查询人员样本列表
func (api *Api) Selectperandsamp() {
	var (
		person    Persons
		personArr PersonArr
		perEx     []PerEx
		stmt      string
		err       error
	)
	api.UnJson(&person)
	api.UnJson(&perEx)
	if person.RowBegin == 0 && person.RowEnd == 0 {
		person.RowBegin = 0
		person.RowEnd = 20
	}
	stmt = ` select (select count(*) from per inner join sampling on sampling.perid=per.id and per.id is not null) as cnt,
				per.*,sampling.Code as sample_lab_no,sampling.id as sampid from per inner join sampling on sampling.perid=per.id and  per.id is not null
				where 1=1 and per.deleted = 0 and sampling.deleted = 0
				#and sampling.Code = {{empty .SampleLabNo}}#
				#and per.Passno = {{empty .IdCardNo}}#
				#and per.Name = {{empty .PersonName}}#
				#and per.Fillingdate = {{empty .CollectDatetime}}#
				order by per.Insertdatetime desc 
				limit ` + strconv.Itoa(int(person.RowBegin)) + "," + strconv.Itoa(int(person.RowEnd))
	stmt, err = xml2sql.GenerateSql(stmt, &person.Personinfo)
	api.CheckWrite(err)
	api.DB().Raw(stmt).Scan(&perEx)
	PersonJson := make([]lsent.Personinfo, len(perEx))
	for i, k := range perEx {
		PersonJson[i].Id = qianID(k.Sampid) //返回给前端 服务器编码(6b)+10b
		PersonJson[i].PersonName = k.Name
		PersonJson[i].Gender = SelectDictByValue(k.Gender, "5")
		PersonJson[i].Race = SelectDictByValue(k.Folk, "14")
		PersonJson[i].BirthDatetime = k.Birthday
		PersonJson[i].NativePlaceRegion = k.Permanentaddr
		PersonJson[i].NativePlaceAddr = k.Houseaddr
		PersonJson[i].IdCardNo = k.Passno
		PersonJson[i].SampleLabNo = k.SampleLabNo
		PersonJson[i].CollectUser = k.Fillingperson
		PersonJson[i].CollectDatetime = k.Fillingdate
		PersonJson[i].CollectAgencyCode = strconv.Itoa(int(k.Region))
		PersonJson[i].CollectAgencyName = k.Fillingunit
		PersonJson[i].Server = strconv.Itoa(int(k.Server))
		PersonJson[i].Remark = k.Memo
		PersonJson[i].PersonId = qianID(k.Id)
		PersonJson[i].CreateUser = strconv.Itoa(int(k.Insertoperatorid))
		PersonJson[i].CreateDatetime = k.Insertdatetime
		PersonJson[i].UpdateUser = strconv.Itoa(int(k.Updateoperatorid))
		PersonJson[i].UpdateDatetime = k.Lastupdatetime
		// PersonJson[i].Cnt = k.Cnt
	}
	personArr.Personinfo = PersonJson
	personArr.Cnt = perEx[0].Cnt
	api.JsonWrite(personArr)
}

//通过id,id_card_no查询人员样本信息
func (api *Api) Selectperandsampbyid() {
	var (
		person lsent.Personinfo
		perEx  PerEx
		stmt   string
		err    error
	)
	api.UnJson(&person)
	person.Id = strconv.Itoa(int(houID(person.Id)))
	stmt = `select per.*,
				   sampling.Code as sample_lab_no,
				   sampling.id as sampid
				from per inner join sampling on (sampling.perid=per.id and per.id is not null)
				where 1=1 and per.deleted = 0 and sampling.deleted = 0
				#and sampling.id = {{empty .Id}}#
				#and per.Passno = {{empty .IdCardNo}}#`
	stmt, err = xml2sql.GenerateSql(stmt, &person)
	api.CheckWrite(err)
	api.DB().Raw(stmt).Scan(&perEx)
	person.Id = qianID(perEx.Sampid)
	person.PersonName = perEx.Name
	person.Gender = perEx.Gender
	person.Race = perEx.Folk
	person.BirthDatetime = perEx.Birthday
	person.NativePlaceRegion = perEx.Permanentaddr
	person.NativePlaceAddr = perEx.Houseaddr
	person.IdCardNo = perEx.Passno
	person.SampleLabNo = perEx.SampleLabNo
	person.CollectUser = perEx.Fillingperson
	person.CollectDatetime = perEx.Fillingdate
	person.CollectAgencyCode = strconv.Itoa(int(perEx.Region))
	person.CollectAgencyName = perEx.Fillingunit
	person.Server = strconv.Itoa(int(perEx.Server))
	person.Remark = perEx.Memo
	person.PersonId = qianID(perEx.Id)
	person.CreateUser = strconv.Itoa(int(perEx.Insertoperatorid))
	person.CreateDatetime = perEx.Insertdatetime
	person.UpdateUser = strconv.Itoa(int(perEx.Updateoperatorid))
	person.UpdateDatetime = perEx.Lastupdatetime
	api.JsonWrite(person)
}

//新增，修改
func (api *Api) Createperandsamp() {
	var (
		person   Person
		per      lsent.Per
		pers     []lsent.Per
		sampling lsent.Sampling
		usr      lsent.Usr
		err      error
	)
	api.UnJson(&person)
	api.UnJson(&per)
	api.UnJson(&sampling)
	r, err := api.R.Cookie("uid")
	api.CheckWrite(err)
	if r != nil {
		uidstr, _ := strconv.Atoi(r.Value)
		api.CheckWrite(api.DB().First(&usr, int64(uidstr)).Error)
	}

	if person.PersonInfo.SampleLabNo == "" {
		err = errors.New("样本实验室编号为空")
		api.CheckWrite(err)
	}
	db := api.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	if person.PersonInfo.Id != "" {
		if person.PersonInfo.IdCardNo != "" {
			api.CheckWrite(db.Where("Passno = ? ", person.PersonInfo.IdCardNo).Find(&pers).Error)
			if len(pers) > 0 {
				if len(pers) > 1 || (len(pers) == 1 && strconv.Itoa(int(pers[0].Id)) != person.PersonInfo.Id) {
					err = errors.New("身份证号重复")
					api.CheckWrite(err)
				}
			}
		}
		personID := houID(person.PersonInfo.PersonId)
		samplingID := houID(person.PersonInfo.Id)
		api.CheckWrite(api.DB().Where("id = ?", personID).Find(&per).Error)
		api.CheckWrite(api.DB().Where("id = ?", samplingID).Find(&sampling).Error)
		per, sampling = PersonToPerAndSam(person.PersonInfo)
		per.Updateoperatorid = usr.Id
		// per.Lastupdatetime = timex.Time(time.Now())
		//时间转化为string，layout必须为 "2006-01-02 15:04:05"
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		per.Lastupdatetime = timeNow
		if per.Jaildate == "" {
			per.Jaildate = "0000-00-00 00:00:00"
		}

		sampling.Updateoperatorid = usr.Id
		// sampling.Lastupdatetime = timex.Time(time.Now())
		sampling.Lastupdatetime = timeNow
		api.CheckWrite(db.Save(&per).Error)
		api.CheckWrite(db.Save(&sampling).Error)
		// person.PersonInfo.UpdateDatetime = timex.Time(time.Now())
		person.PersonInfo.UpdateDatetime = timeNow
		person.PersonInfo.UpdateUser = usr.Name
	} else {
		if person.PersonInfo.IdCardNo != "" {
			api.CheckWrite(db.Where("Passno = ? ", person.PersonInfo.IdCardNo).Find(&pers).Error)
			if len(pers) > 0 {
				err = errors.New("身份证号重复")
				api.CheckWrite(err)
			}
		}
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		per, sampling = PersonToPerAndSam(person.PersonInfo)
		if person.PersonInfo.SampleLabNo == "" { //样本实验室编号为空则自动生成
			sampling.Code = "t" + g.Configure.Server
		}
		per.Insertoperatorid = usr.Id
		// per.Insertdatetime = timex.Time(time.Now())
		per.Insertdatetime = timeNow
		sampling.Insertoperatorid = usr.Id
		if per.Jaildate == "" {
			per.Jaildate = "0000-00-00 00:00:00"
		}
		if sampling.Samplingdate == "" {
			sampling.Samplingdate = "0000-00-00 00:00:00"
		}
		if sampling.Accepttime == "" {
			sampling.Accepttime = "0000-00-00 00:00:00"
		}
		if sampling.Checkingdate == "" {
			sampling.Checkingdate = "0000-00-00 00:00:00"
		}
		if per.Abscondencedate == "" {
			per.Abscondencedate = "0000-00-00 00:00:00"
		}
		if sampling.Lastupdatetime == "" {
			sampling.Lastupdatetime = "0000-00-00 00:00:00"
		}
		if per.Lastupdatetime == "" {
			per.Lastupdatetime = "0000-00-00 00:00:00"
		}
		sampling.Insertdatetime = timeNow

		per.Fillingperson = usr.Name
		per.Fillingdate = timeNow
		per.Fillingunit = g.Configure.CollectAgencyName
		region, _ := strconv.Atoi(Substr(g.Configure.CollectAgencyCode, 4, 8))
		per.Region = int64(region)
		api.CheckWrite(db.Create(&per).Error)
		sampling.Perid = per.Id
		api.CheckWrite(db.Create(&sampling).Error)
		// person.PersonInfo.CreateDatetime = timex.Time(time.Now())
		person.PersonInfo.CreateDatetime = timeNow
		person.PersonInfo.CreateUser = usr.Name
		person.PersonInfo.Id = qianID(sampling.Id)
		person.PersonInfo.PersonId = qianID(per.Id)
	}
	api.CheckWrite(db.Commit().Error)
	api.JsonWrite(person.PersonInfo)

}

//字典查找
func (api *Api) SelectSysDict() {
	var (
		sysDict lsent.Sys_Dict
		dict    []lsent.Dict
		dd      lsent.Dict
	)
	api.UnJson(&sysDict)
	api.UnJson(&dict)
	k := sysDict.DictCategory
	if k == "" {
		err := api.DB().Where("value=?", "7811").Find(&dd).Error
		if err != nil {
			db := api.DB().Begin()
			defer func() {
				if db.Error != nil {
					db.Rollback()
				}
			}()
			api.CheckWrite(db.Exec(`INSERT INTO dict (category,sub,value,word,ord,note,deleted) VALUES(?,?,?,?,?,?,?)`, 201, 0, "7811", "1", 0, "DNA编号规则", 0).Error)
			api.CheckWrite(db.Commit().Error)
		}
		api.CheckWrite(api.DB().Raw(`select * from dict where (category in (5,11,14,22,32,33,37,45,83,90))or value =7811`).Scan(&dict).Error)
	} else {
		api.CheckWrite(api.DB().Where("category=?", sysDict.DictCategory).Find(&dict).Error)
	}
	var sysDicts = make([]lsent.Sys_Dict, 0)
	if k == "" {
		sysDictsTmp := make([]lsent.Dict, 0)
		confjson := `[{"category": 201,"value": "7821","word": "6500001"},{"category": 201,"value": "7822","word": "新疆省厅"},
				{"category": 201,"value": "7823","word": "001"},{"category": 201,"value": "7824","word": "650000"},
				{"category": 201,"value": "7825","word": "新疆省厅"},{"category": 201,"value": "7826","word": "1"},
				{"category": 201,"value": "7827","word": "1"},{"category": 201,"value": "7828","word": "1"}]`
		api.CheckWrite(json.Unmarshal([]byte(confjson), &sysDictsTmp))
		for _, t := range sysDictsTmp {
			dict = append(dict, t)
		}
	}
	for _, k := range dict {
		sysDict.DictCategory = strconv.Itoa(int(k.Category))
		sysDict.DictNationalKey = k.Value
		sysDict.DictKey = k.Value
		sysDict.DictValue1 = k.Word
		sysDict.Ord = k.Ord
		sysDicts = append(sysDicts, sysDict)
	}
	api.JsonWrite(sysDicts)
}

//用户查找
func (api *Api) SelectSysUser() {
	var (
		sysUser lsent.Sys_User
		usr     []lsent.Usr
	)
	api.UnJson(sysUser)
	api.UnJson(usr)
	api.DB().Raw(`select * from usr`).Scan(&usr)
	var sysUsers = make([]lsent.Sys_User, 0)
	for _, k := range usr {
		sysUser = UsrToSysUser(k)
		sysUsers = append(sysUsers, sysUser)
	}
	api.JsonWrite(sysUsers)
}

//删除人员样本信息
func (api *Api) DeletePersonSampling() {
	var (
		person   Person
		per      lsent.Per
		sampling lsent.Sampling
	)
	api.UnJson(&person)
	api.UnJson(&per)
	api.UnJson(&sampling)
	per, sampling = PersonToPerAndSam(person.PersonInfo)
	if per.Id == 0 || sampling.Id == 0 {
		err := errors.New("人员id或样本id为空")
		api.CheckWrite(err)
	}
	db := api.DB().Begin()
	defer func() {
		if db.Error != nil {
			db.Rollback()
		}
	}()
	per.Deleted = 1
	sampling.Deleted = 1
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	per.Lastupdatetime = timeNow
	sampling.Lastupdatetime = timeNow
	api.CheckWrite(db.Save(&per).Error)
	api.CheckWrite(db.Save(&sampling).Error)
	api.CheckWrite(db.Commit().Error)
	api.JsonWrite(person.PersonInfo)
}

//获取采集人信息
// func (api *Api) SelectInformation() {
// 	var usr lsent.Usr
// 	r, err := api.R.Cookie("uid")
// 	api.CheckWrite(err)
// 	if r != nil {
// 		uidstr, _ := strconv.Atoi(r.Value)
// 		api.CheckWrite(api.DB().First(&usr, int64(uidstr)).Error)
// 	}
// 	pp := AA[strconv.Itoa(int(usr.Id))]
// 	if pp == nil {
// 		f, err := os.Open(filepath.Join(g.Configure.PickInformation, "information.txt"))
// 		api.CheckWrite(err)
// 		defer f.Close()
// 		buf, err := ioutil.ReadAll(f)
// 		api.CheckWrite(err)
// 		js, err := simplejson.NewJson(buf)
// 		if err != nil {
// 			js = simplejson.New()
// 		}
// 		arrJSON := js.Get(strconv.Itoa(int(usr.Id)))
// 		fmt.Println(arrJSON)
// 		api.JsonWrite(arrJSON)
// 	} else {
// 		api.JsonWrite(pp)
// 	}
// }

//获取当前登录用户
func (api *Api) QuerySysUser() {
	var (
		usr     lsent.Usr
		sysuser lsent.Sys_User
	)
	api.UnJson(&usr)
	r, err := api.R.Cookie("uid")
	api.CheckWrite(err)
	if r != nil {
		uidstr, _ := strconv.Atoi(r.Value)
		api.CheckWrite(api.DB().First(&usr, int64(uidstr)).Error)
	}
	sysuser = UsrToSysUser(usr)
	api.JsonWrite(sysuser)

}

type PersonInfo struct {
	Success []PersonEx `json:"success"`
	Failure []PersonEx `json:"failure"`
}
type PersonEx struct {
	lsent.Personinfo
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

//查重身份证
func (api *Api) SelectCheckPersonByIdcardNo() {
	var (
		person     []PersonEx
		pers       []lsent.Per
		samplings  []lsent.Sampling
		personinfo PersonInfo
	)
	api.UnJson(&person)
	perSucc := make([]PersonEx, 0)
	perFail := make([]PersonEx, 0)
	for _, k := range person {
		if k.IdCardNo != "" || k.SampleLabNo != "" {
			k.Gender = SelectDictByValue(k.Gender, "5")
			k.Race = SelectDictByValue(k.Race, "14")
			if k.Gender == "" || len(k.Gender) == 0 {
				k.Status = "3" //3表示性别为空
				k.Msg = "性别为空"
				perFail = append(perFail, k)

			} else if k.Race == "" || len(k.Race) == 0 {
				k.Status = "4" //3表示民族为空
				k.Msg = "民族为空"
				perFail = append(perFail, k)
			} else if k.IdCardNo != "" {
				api.CheckWrite(api.DB().Where("Passno = ? ", k.IdCardNo).Find(&pers).Error)
				if len(pers) > 0 {
					k.Status = "1" //1表示身份证重复
					k.Msg = "身份证重复"
					perFail = append(perFail, k)
				}
			} else if k.SampleLabNo != "" {
				api.CheckWrite(api.DB().Where("code = ? ", k.SampleLabNo).Find(&samplings).Error)
				if len(pers) > 0 {
					k.Status = "2" //1表示样本dna编号重复
					k.Msg = "样本实验室编号重复"
					perFail = append(perFail, k)
				}
			}
			if len(k.Status) == 0 {
				perSucc = append(perSucc, k)
			}
		}
	}
	personinfo.Failure = perFail
	personinfo.Success = perSucc
	api.JsonWrite(personinfo)
}

//批量插入建库人员
func (api *Api) Uploadperson() {
	var (
		person     []PersonEx
		personinfo PersonInfo
		usr        lsent.Usr
		err        error
	)
	api.UnJson(&person)
	perSucc := make([]PersonEx, 0)
	perFail := make([]PersonEx, 0)
	for _, k := range person {
		db := api.DB().Begin()
		defer func() {
			if db.Error != nil {
				db.Rollback()
			}
		}()
		k.Id = ""
		k.PersonId = ""
		per, sampling := PersonToPerAndSam(k.Personinfo)
		if k.SampleLabNo == "" {
			sampling.Code = "t" + g.Configure.Server + ""
		}
		err = db.Where("code=?", k.CreateUser).Find(&usr).Error
		if err != nil {
			k.Msg = "未找到该用户"
			k.Flag = 1
			perFail = append(perFail, k)
			continue
		}
		per.Insertoperatorid = usr.Id
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		if per.Birthday == "" {
			per.Birthday = "0000-00-00 00:00:00"
		}
		if per.Jaildate == "" {
			per.Jaildate = "0000-00-00 00:00:00"
		}
		if per.Abscondencedate == "" {
			per.Abscondencedate = "0000-00-00 00:00:00"
		}
		if sampling.Lastupdatetime == "" {
			sampling.Lastupdatetime = "0000-00-00 00:00:00"
		}
		if per.Lastupdatetime == "" {
			per.Lastupdatetime = "0000-00-00 00:00:00"
		}
		if sampling.Samplingdate == "" {
			sampling.Samplingdate = "0000-00-00 00:00:00"
		}
		if sampling.Accepttime == "" {
			sampling.Accepttime = "0000-00-00 00:00:00"
		}
		if sampling.Checkingdate == "" {
			sampling.Checkingdate = "0000-00-00 00:00:00"
		}
		per.Fillingperson = usr.Name
		per.Fillingdate = timeNow
		per.Fillingunit = g.Configure.CollectAgencyName
		if !Checkpercode(per.Code) {
			k.Msg = "人员编号已存在"
			k.Flag = 1
			perFail = append(perFail, k)
			continue
		}
		if !Checkidcardno(per.Passno) {
			k.Msg = "身份证号已存在"
			k.Flag = 1
			perFail = append(perFail, k)
			continue
		}
		if per.Folk == "" {
			k.Msg = "民族为空"
			k.Flag = 1
			perFail = append(perFail, k)
			continue
		}
		if per.Gender == "" {
			k.Msg = "性别为空"
			k.Flag = 1
			perFail = append(perFail, k)
			continue
		}
		region, _ := strconv.Atoi(Substr(g.Configure.CollectAgencyCode, 4, 8))
		per.Region = int64(region)
		per.Insertdatetime = timeNow
		sampling.Insertoperatorid = usr.Id
		sampling.Insertdatetime = timeNow
		api.CheckWrite(db.Create(&per).Error)
		sampling.Perid = per.Id
		api.CheckWrite(db.Create(&sampling).Error)
		if k.Personinfo.Foot != nil {
			fmt.Println(len(k.Personinfo.Foot))
			photo := PerandSamToPhoto(per, sampling, k.Personinfo)
			api.CheckWrite(db.Create(&photo).Error)
		}
		k.CreateDatetime = timeNow
		k.CreateUser = usr.Name
		k.Flag = 2
		k.Id = qianID(sampling.Id)
		k.PersonId = qianID(per.Id)
		perSucc = append(perSucc, k)
		api.CheckWrite(db.Commit().Error)
	}
	perSucc2 := make([]PersonEx, 0)
	perFail2 := make([]PersonEx, 0)
	for _, k2 := range perSucc {
		k2.Gender = SelectDictByValue(k2.Gender, "5")
		k2.Race = SelectDictByValue(k2.Race, "14")
		perSucc2 = append(perSucc2, k2)

	}
	for _, k3 := range perFail {
		k3.Gender = SelectDictByValue(k3.Gender, "5")
		k3.Race = SelectDictByValue(k3.Race, "14")
		perFail2 = append(perFail2, k3)
	}
	personinfo.Success = perSucc2
	personinfo.Failure = perFail2
	api.JsonWrite(personinfo)
}

//检查per_code唯一性
func Checkpercode(code string) bool {
	var per lsent.Per
	err := g.DB0.Where("code=?", code).Find(&per).Error
	if err != nil {
		return true
	}
	return false
}

//检查身份证唯一性
func Checkidcardno(idcardno string) bool {
	var per lsent.Per
	err := g.DB0.Where("Passno=?", idcardno).Find(&per).Error
	if err != nil {
		return true
	}
	return false
}
func PerandSamToPhoto(per lsent.Per, sam lsent.Sampling, person lsent.Personinfo) (ph lsent.Photo) {
	ph.Id = per.Id
	ph.Deleted = 0
	ph.Photo = person.Foot
	ph.Photolen = strconv.Itoa(len(person.Foot))
	ph.Server = per.Server
	ph.Insertdatetime = time.Now().Format("2006-01-02 15:04:05")
	ph.Updatedatetime = time.Now().Format("2006-01-02 15:04:05")
	return ph

}

//数据转换
func PersonToPerAndSam(person lsent.Personinfo) (per lsent.Per, sampling lsent.Sampling) {
	per.Name = person.PersonName
	// per.Gender = SelectDictByCategory(person.Gender)
	per.Gender = person.Gender
	per.Nationality = SelectDictByCategory("中国")
	region, _ := strconv.Atoi(Substr(person.CollectAgencyCode, 4, 8))
	per.Region = int64(region)
	per.Folk = person.Race
	per.Birthday = person.BirthDatetime
	per.Permanentaddr = person.NativePlaceRegion
	per.Houseaddr = person.NativePlaceAddr
	per.Fillingdate = person.CollectDatetime
	per.Fillingperson = person.CollectUser
	per.Passno = person.IdCardNo
	per.Fillingunit = person.CollectAgencyName
	// sever, _ := strconv.Atoi(person.Server)
	per.Server = 0
	per.Photolen = int64(len(person.Foot))
	per.Memo = person.Remark
	per.Code = "p" + person.SampleLabNo
	per.Category = 4
	per.Passtype = "1"
	per.Id = houID(person.PersonId)
	sampling.Code = person.SampleLabNo
	sampling.Category = 4
	sampling.Region = int64(region)
	sampling.Id = houID(person.Id)
	return per, sampling
}

//用户表信息转换
func UsrToSysUser(usr lsent.Usr) (sysUser lsent.Sys_User) {
	sysUser.Id = strconv.Itoa(int(usr.Id))
	sysUser.Password = usr.Password
	sysUser.UserName = usr.Code
	sysUser.InitServerNo = strconv.Itoa(int(usr.Server))
	sysUser.FullName = usr.Name
	// Region
	// sysUser.Id = strconv.Itoa(int(usr.Id))
	// sysUser.Id = strconv.Itoa(int(usr.Id))
	// sysUser.Id = strconv.Itoa(int(usr.Id))
	// sysUser.Id = strconv.Itoa(int(usr.Id))
	// sysUser.Id = strconv.Itoa(int(usr.Id))
	return sysUser
}

//字典中文转换为编码
func SelectDictByCategory(dict_value1 string) (value string) {
	var dict lsent.Dict
	if dict_value1 == "" {
		value = ""
	} else {
		g.DB0.Where("word=?", dict_value1).Find(&dict)
		value = dict.Value
	}
	return value
}

//字典编码转换为中文
func SelectDictByValue(DictNationalKey, category string) (word string) {
	var dict []lsent.Dict
	if DictNationalKey == "" || category == "" {
		word = ""
	} else {
		g.DB0.Where("category=?", category).Find(&dict)
		for _, k := range dict {
			if DictNationalKey == k.Value {
				word = k.Word
			}
		}
		if len(word) == 0 {
			word = ""
		}
	}
	return word
}

//截取字符串
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//截取前端id转为后台id
func houID(id string) (houId int64) {
	tt := Substr(id, 6, 10)
	ff, _ := strconv.Atoi(tt)
	houId = int64(ff)
	return houId
}

//后台id转为前端id
func qianID(id int64) (qianId string) {
	if len(g.Configure.Server) == 4 {
		g.Configure.Server = g.Configure.Server + "00"
	}
	var tt string = "" //id后10位补0
	for j := 0; j < 10-len(strconv.Itoa(int(id))); j++ {
		tt = tt + "0"
	}
	qianId = g.Configure.Server + tt + strconv.Itoa(int(id))
	return qianId
}

//样本实验编号为空则自动生成
// card.CheckWrite(card.DB().Where("value=?", "7811").Find(&sysdict).Error)
// sysdict.DictValue1 = sysdict.DictValue1 + "1"
// someper.Person_Info.SampleLabNo = g.Configure.Server + g.Configure.Code + sysdict.DictValue1
// card.CheckWrite(db.Save(&sysdict).Error)
// card.CheckWrite(err)
// body, err := ioutil.ReadAll(resp.Body)
// card.CheckWrite(err)
// defer resp.Body.Close()
// err = json.Unmarshal(body, &resultcheckonline)
// card.JsonWrite(resultcheckonline.Result)

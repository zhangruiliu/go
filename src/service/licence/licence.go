package service

import (
	"crypto/md5"
	"strings"
	"encoding/hex"
	"time"
	"strconv"
	"dbdefault"
	"fmt"
	"net"
	"os/exec"
	"io/ioutil"
	"web/httpcontext"
	"g"
	"config"
)

//var prefix = "haihuaxinan:"

type Licence struct {
	httpcontext.HttpContext
}
type SecretInfo struct{
	Mac_Addr string
	Cpu_Id string
	Ip_Addr string
}

func (l *Licence)CheckLicence() {
	confLicence := config.IniConfig{}
	licenceConfig,err := confLicence.ParseFile("licence.txt")
	l.Check(err)

	licenceValue := licenceConfig.String("LICENCE")
	dateLimitValue := licenceConfig.String("DATELIMIT")
	if dbdefault.IsStringNullOrEmpty(licenceValue){
		l.CheckWrite(fmt.Errorf("请验证你的许可文件是否合法"))
	}

	Secret := make([]SecretInfo,0)
	mac,_ := net.Interfaces()//获取mac addr
	for _,v := range mac {
		if v.HardwareAddr != nil {
			tmp := SecretInfo{}
			tmp.Mac_Addr = fmt.Sprintf("%s", v.HardwareAddr)
			Secret = append(Secret, tmp)
		}
	}
	ip,_ := net.InterfaceAddrs()
	var i1 = 0
	for _, addr := range ip {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Secret[i1].Ip_Addr = ipnet.IP.To4().String()
				i1 += 1
			}
		}
	}
	cmd := exec.Command("cmd" ,"/C","wmic cpu get processorid")//获取cpu id
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	opBytes, _ := ioutil.ReadAll(stdout)
	a := strings.Replace( strings.TrimSpace(string(opBytes)),"\n","",-1)
	b := strings.Replace(a,"\r", "",-1)
	b = strings.Replace(b,"\t","", -1)
	b = strings.Replace(b," ","", -1)
	opTmpString := strings.Split(string(b),"ProcessorId")
	if len(opTmpString) > 1 {
		Secret[0].Cpu_Id = opTmpString[1]
	}
	/*
	sumByte := md5.Sum([]byte(ip+mac+cpu_id))
	LICENCE := strings.ToUpper(hex.EncodeToString(sumByte[:]))
	*/

	md5Arr := make([]string,len(Secret))
	for i,_ := range Secret {
		//fmt.Println(Secret[i].Ip_Addr+Secret[i].Mac_Addr+Secret[0].Cpu_Id)
		res := md5.Sum([]byte(Secret[i].Ip_Addr+Secret[i].Mac_Addr+Secret[0].Cpu_Id))
		md5Arr[i] = strings.ToUpper(hex.EncodeToString(res[:]))
	}
	//设置软件的激活日期
	timeString := time.Now()
	src,err := hex.DecodeString(dateLimitValue)
	l.Check(err)
	dateValue := strings.Split(string(src),":")
	date := strings.Split(dateValue[1],"-")
	yyyy,_ := strconv.Atoi(date[0])
	mm,_ := strconv.Atoi(date[1])
	dd,_ := strconv.Atoi(date[2])
	setTimeDate := time.Date(yyyy,time.Month(mm),dd,0,0,0,0,time.Local)
	flag := false
	for  _,v := range md5Arr {
		//fmt.Println(v)
		if licenceValue == v && timeString.Before(setTimeDate){
			flag = true
			l.JsonWrite(g.SUCCESS_INFO)
			break
		}
	}
	if !flag{
		l.JsonWrite(g.FAILURE_INFO)
	}
}


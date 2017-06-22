package lsent

type Per struct {
	Id            int64  `json:"id,omitempty"`
	Server        int64  `json:"server,omitempty"`   //服务器编号0
	Uploaded      int64  `json:"uploaded,omitempty"` //上报 0
	Region        int64  `json:"region,omitempty"`   //0
	Category      int64  `json:"category,omitempty"` //默认为4
	Commitid      int64  `json:"commitid,omitempty"` //默认为0
	Code          string `json:"code,omitempty"`     //本地人员编号
	Name          string `json:"name,omitempty"`     //人员姓名
	Kind          string `json:"kind,omitempty"`     //词典表category=90的第一个
	Class         string `json:"class,omitempty"`
	Codesampling  string `json:"codesampling,omitempty"`  //外部系统人员编号
	Alias         string `json:"alias,omitempty"`         //人员别名
	Gender        string `json:"gender,omitempty"`        //性别
	Birthday      string `json:"birthday,omitempty"`      //出生日期
	Nationality   string `json:"nationality,omitempty"`   //国籍
	Folk          string `json:"folk,omitempty"`          //民族
	Passtype      string `json:"passtype,omitempty"`      //证件类型（1为身份证）
	Passno        string `json:"passno,omitempty"`        //证件号码
	Permanentaddr string `json:"permanentaddr,omitempty"` //户籍地
	Houseaddr     string `json:"houseaddr,omitempty"`     //现住址
	// JailLocation        string    `json:"jail_location,omitempty"`
	Jaildate            string `json:"jaildate,omitempty"`
	Casetype            string `json:"casetype,omitempty"`
	Caselevel           string `json:"caselevel,omitempty"`
	Fillingunit         string `json:"fillingunit,omitempty"`   //采样单位
	Fillingperson       string `json:"fillingperson,omitempty"` //采样人
	Fillingdate         string `json:"fillingdate,omitempty"`   //采样时间
	Memo                string `json:"memo,omitempty"`          //备注
	Labassemblyno       string `json:"labassemblyno,omitempty"`
	Labmemo             string `json:"labmemo,omitempty"`
	Deductedreason      string `json:"deductedreason,omitempty"`
	Featureinfo         string `json:"featureinfo,omitempty"`
	Undertakeunit       string `json:"undertakeunit,omitempty"`
	Undertakeperson     string `json:"undertakeperson,omitempty"`
	Briefinfo           string `json:"briefinfo,omitempty"`
	Caseno              string `json:"caseno,omitempty"`
	Caseinfo            string `json:"caseinfo,omitempty"`
	Abscondencedate     string `json:"abscondencedate,omitempty"`
	Abscondencelocation string `json:"abscondencelocation,omitempty"`
	Occupation          string `json:"occupation,omitempty"`
	Studytarget         string `json:"studytarget,omitempty"`
	Samplingrange       string `json:"samplingrange,omitempty"`
	Samplecount         int64  `json:"samplecount,omitempty"`
	Status              int64  `json:"status,omitempty"`
	Insertoperatorid    int64  `json:"insertoperatorid,omitempty"` //插入人
	Insertdatetime      string `json:"insertdatetime,omitempty"`   //插入时间
	Updateoperatorid    int64  `json:"updateoperatorid,omitempty"` //更新人
	Lastupdatetime      string `json:"lastupdatetime,omitempty"`   //更新时间
	Photolen            int64  `json:"photolen,omitempty"`
	Photo               string `json:"photo,omitempty"`
	Keepingperson       int64  `json:"keepingperson,omitempty"`
	Deleted             int64  `json:"deleted,omitempty"` //0
}

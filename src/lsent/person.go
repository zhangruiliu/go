package lsent

type Personinfo struct {
	Id                string `json:"id,omitempty"`        //代表样本id
	PersonName        string `json:"person_name"`         //人员名称
	Gender            string `json:"gender"`              //性别
	Race              string `json:"race,omitempty"`      //名族
	BirthDatetime     string `json:"birth_datetime"`      //出生日期
	NativePlaceRegion string `json:"native_place_region"` //现住址
	NativePlaceAddr   string `json:"native_place_addr"`   //现住址
	IdCardNo          string `json:"id_card_no"`          //身份证号
	SampleLabNo       string `json:"sample_lab_no"`       //样本实验室编号
	CollectUser       string `json:"collect_user"`        //采集人
	CollectDatetime   string `json:"collect_datetime"`    // 采集单位时间
	CollectAgencyCode string `json:"agency_code"`         //采集单位
	CollectAgencyName string `json:"collect_agency_name"` //采集单位名称
	Server            string `json:"server,omitempty"`    //服务器编号
	Remark            string `json:"remark,omitempty"`    //备注
	PersonId          string `json:"person_id,omitempty"` //人员id
	CreateUser        string `json:"create_user"`
	CreateDatetime    string `json:"create_datetime"`
	UpdateUser        string `json:"update_user"`
	UpdateDatetime    string `json:"update_datetime"`
	SampleType        string `json:sample_type`
	Flag              int64  `json:"flag"`
	Foot              []byte `json:foot`
}

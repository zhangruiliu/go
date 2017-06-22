package lsent

import "timex"

type Sys_Dict struct {
	Id              string     `json:"id,omitempty" sql:"field" xml:"id"`                             //ID
	DictCategory    string     `json:"dict_category,omitempty" sql:"field" xml:"dictCategory"`        //字典类型
	DictSubCategory string     `json:"dict_sub_category,omitempty" sql:"field" xml:"dictSubCategory"` //字典子类型
	RootFlag        int64      `json:"root_flag,omitempty" sql:"field" xml:"rootFlag"`                //是否根节点
	DictKey         string     `json:"dict_key,omitempty" sql:"field" xml:"dictKey"`                  //字典主键
	DictNationalKey string     `json:"dict_national_key,omitempty" sql:"field" xml:"dictNationalKey"` //字典国家库主键
	DictValue1      string     `json:"dict_value1,omitempty" sql:"field" xml:"dictValue1"`            //字典值1
	DictValue2      string     `json:"dict_value2,omitempty" sql:"field" xml:"dictValue2"`            //字典值2
	DictValue3      string     `json:"dict_value3,omitempty" sql:"field" xml:"dictValue3"`            //字典值3
	DictAlias       string     `json:"dict_alias,omitempty" sql:"field" xml:"dictAlias"`              //字典别名
	ParentKey       string     `json:"parent_key,omitempty" sql:"field" xml:"parentKey"`              //字典父节点KEY
	DownloadFlag    int64      `json:"download_flag,omitempty" sql:"field" xml:"downloadFlag"`        //下载标识
	ReadonlyFlag    int64      `json:"readonly_flag,omitempty" sql:"field" xml:"readonlyFlag"`        //只读标识
	Ord             int64      `json:"ord,omitempty" sql:"field" xml:"ord"`                           //显示顺序
	DictPy          string     `json:"dict_py,omitempty" sql:"field" xml:"dictPy"`                    //字典拼音
	ActiveFlag      int64      `json:"active_flag,omitempty" sql:"field" xml:"activeFlag"`            //启用标识
	Remark          string     `json:"remark,omitempty" sql:"field" xml:"remark"`                     //备注
	CreateUser      string     `json:"create_user,omitempty" sql:"field" xml:"createUser"`            //创建人
	CreateDatetime  timex.Time `json:"create_datetime,omitempty" sql:"field" xml:"createDatetime"`    //创建时间
	UpdateUser      string     `json:"update_user,omitempty" sql:"field" xml:"updateUser"`            //更新人
	UpdateDatetime  timex.Time `json:"update_datetime,omitempty" sql:"field" xml:"updateDatetime"`    //更新时间

}

package lsent

import "timex"

type Sys_User struct {
	Id                      string     `json:"id,omitempty" sql:"field" xml:"id"`                                            //ID
	InitServerNo            string     `json:"init_server_no,omitempty" sql:"field" xml:"initServerNo"`                      //原始服务器编号
	LabId                   string     `json:"lab_id,omitempty" sql:"field" xml:"labId"`                                     //实验室ID
	UserType                string     `json:"user_type,omitempty" sql:"field" xml:"userType"`                               //用户类型
	UserName                string     `json:"user_name,omitempty" sql:"field" xml:"userName"`                               //登录用户名
	Password                string     `json:"password,omitempty" sql:"field" xml:"password"`                                //登录密码
	FullName                string     `json:"full_name,omitempty" sql:"field" xml:"fullName"`                               //姓名
	Phone                   string     `json:"phone,omitempty" sql:"field" xml:"phone"`                                      //联系方式
	CertificateType         string     `json:"certificate_type,omitempty" sql:"field" xml:"certificateType"`                 //证件类型
	CertificateNo           string     `json:"certificate_no,omitempty" sql:"field" xml:"certificateNo"`                     //证件号码
	Email                   string     `json:"email,omitempty" sql:"field" xml:"email"`                                      //电子邮件
	OrganizationRegionalism string     `json:"organization_regionalism,omitempty" sql:"field" xml:"organizationRegionalism"` //用户单位行政区划
	OrganizationName        string     `json:"organization_name,omitempty" sql:"field" xml:"organizationName"`               //用户单位名称
	LoginIpRange            string     `json:"login_ip_range,omitempty" sql:"field" xml:"loginIpRange"`                      //登陆IP段限制
	ActiveFlag              int64      `json:"active_flag,omitempty" sql:"field" xml:"activeFlag"`                           //是否启用
	IsLabEmployee           int64      `json:"is_lab_employee,omitempty" sql:"field" xml:"isLabEmployee"`                    //是否为实验室用户标识
	LabEmployeeId           string     `json:"lab_employee_id,omitempty" sql:"field" xml:"labEmployeeId"`                    //实验室用户ID
	IsTransferUser          int64      `json:"is_transfer_user,omitempty" sql:"field" xml:"isTransferUser"`                  //是否为上报用户标识
	Remark                  string     `json:"remark,omitempty" sql:"field" xml:"remark"`                                    //备注
	Va                      string     `json:"va,omitempty" sql:"field" xml:"va"`
	InspectionPersonnelId   string     `json:"inspection_personnel_id,omitempty" sql:"field" xml:"inspectionPersonnelId"` //检验人员id
	DataSource              string     `json:"data_source,omitempty" sql:"field" xml:"dataSource"`                        //数据来源
	DeleteFlag              int64      `json:"delete_flag,omitempty" sql:"field" xml:"deleteFlag"`                        //删除标识
	UserLevel               int64      `json:"user_level,omitempty" sql:"field" xml:"userLevel"`                          //数据级别
	CreateUser              string     `json:"create_user,omitempty" sql:"field" xml:"createUser"`                        //创建人
	CreateDatetime          timex.Time `json:"create_datetime,omitempty" sql:"field" xml:"createDatetime"`                //创建时间
	UpdateUser              string     `json:"update_user,omitempty" sql:"field" xml:"updateUser"`                        //更新人
	UpdateDatetime          timex.Time `json:"update_datetime,omitempty" sql:"field" xml:"updateDatetime"`                //更新时间

}

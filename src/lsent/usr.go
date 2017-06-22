package lsent

import "timex"

type Usr struct {
	Id              int64      `json:"id,omitempty" sql:"field"`
	Password        string     `json:"password,omitempty"` //密码
	Server          int64      `json:"server,omitempty"`   //服务器编号
	Uploaded        int64      `json:"uploaded,omitempty"`
	Region          int64      `json:"region,omitempty"`
	Category        int64      `json:"category,omitempty"`
	Code            string     `json:"code,omitempty"`
	Name            string     `json:"name,omitempty"`
	Role            int64      `json:"role,omitempty"`
	Access          int64      `json:"access,omitempty"`
	Sex             int64      `json:"sex,omitempty"`
	Gender          string     `json:"gender,omitempty"`
	Birthdate       timex.Time `json:"birthdate,omitempty"`
	Post            string     `json:"post,omitempty"`
	Jtitle          string     `json:"jtitle,omitempty"`
	Phone           string     `json:"phone,omitempty"`
	Perdata         string     `json:"perdata,omitempty"`
	Perdatafilename string     `json:"perdatafilename,omitempty"`
	Deleted         int64      `json:"deleted,omitempty"`
}

package lsent

type Check_Person struct {
	Id            string `json:"id,omitempty"`
	Id_Card_No    string `json:"id_card_no,omitempty"` //身份证号
	Server        string `json:"server,omitempty"`     //检验单位
	Status        string `json:"status,omitempty"`     //
	Name          string `json:"name,omitempty"`
	Sample_Lab_No string `json:"sample_lab_no,omitempty"`
}

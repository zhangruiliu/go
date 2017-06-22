package lsent

type Idcard struct {
	Id_Card_No  string `json:"id_card_no,omitempty"`
	Person_Name string `json:"person_name,omitempty"` //身份证号
	Gender      string `json:"gender,omitempty"`      //检验单位
	Race        string `json:"race,omitempty"`        //
	Address     string `json:"address,omitempty"`
}

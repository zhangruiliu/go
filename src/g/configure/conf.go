package configure

// Configure 配置信息
type Configure struct {
	// Parent  *Configure `json:"parent"`
	DB0               DriverInfo `json:"db0"`
	Static            string     `json:"static"`
	LogPath           string     `json:"log-path"`
	Port              string     `json:"port"`
	Lims1LoginUrl     string     `json:"lims1-login-url"`
	SqlTmplPath       string     `json:"sql-tmpl-path"`
	CollectAgencyName string     `json:"collect-agency-name"`
	CollectAgencyCode string     `json:"collect-agency-code"`
	Code              string     `json:"code"`
	Local             string     `json:"local"`
	Public            string     `json:"public"`
	Mode              string     `json:"mode"`
	Server            string     `json:"server"`
	Updatedb          string     `json:"updatedb"`
	Getcheckserver    string     `json:"getcheckserver"`
}

// DriverInfo 驱动信息
type DriverInfo struct {
	Driver      string `json:"driver"`
	Connection  string `json:"connection"`
	Connection2 string `json:"connection2"`
	Connection3 string `json:"connection3"`
}

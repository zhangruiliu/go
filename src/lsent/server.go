package lsent

type Server struct {
	Server       string `json:"server,omitempty"`
	Name         string `json:"name,omitempty"`
	Urlmain      string `json:"urlmain,omitempty"`
	Dsnmainmysql string `json:"dsnmainmysql,omitempty"`
}

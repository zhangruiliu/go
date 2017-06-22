package lsent

type Photo struct {
	Id             int64  `json:"id,omitempty"`
	Server         int64  `json:"server,omitempty"`
	Photolen       string `json:"photolen,omitempty"`
	Photo          []byte `json:"photo,omitempty"`
	Insertdatetime string `json:"insertdatetime,omitempty"`
	Updatedatetime string `json:"updatedatetime,omitempty"`
	Deleted        int64  `json:"deleted,omitempty"`
}

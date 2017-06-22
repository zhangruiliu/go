package lsent

type Dict struct {
	Category int64  `json:"category,omitempty"`
	Sub      int64  `json:"sub,omitempty"`
	Value    string `json:"value,omitempty"`
	Word     string `json:"word,omitempty"`
	Ord      int64  `json:"ord,omitempty"`
	Note     string `json:"note,omitempty"`
	Deleted  int64  `json:"deleted,omitempty"`
}

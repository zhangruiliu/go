package methods

func MakeInClause(s []string) string {
	if s == nil || len(s) == 0 {
		return ""
	}
	inClause := " in ("
	l := len(s)
	exist := false
	for i, v := range s {
		if v == "" {
			continue
		}

		if i == l-1 {
			exist = true
			inClause += EscapeSqlString(v)
		} else {
			exist = true
			inClause += EscapeSqlString(v) + ","
		}
	}
	inClause += ") "
	if exist {
		return inClause
	}
	return ""
}

package xml2sql
import (
	"strings"
	"time"
	"fmt"
)


//转义 ' "
func EscapeQuotes(str string) string {

	str = strings.Replace(str, `'`, `''`, -1)
	str = `'` + str + `'`
	return str
}


func TimeToString(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
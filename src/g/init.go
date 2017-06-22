package g

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"xml2sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func Before() (err error) {
	// DB0, err = db.NewDB(Configure.DB0.Driver, Configure.DB0.Connection)
	// if err != nil {
	// 	return
	// }
	if Configure.Mode == "public" {
		DB0, err = gorm.Open(Configure.DB0.Driver, Configure.DB0.Connection)
		if err != nil {
			return
		}
	} else if Configure.Mode == "local" {
		DB0, err = gorm.Open(Configure.DB0.Driver, Configure.DB0.Connection)
		if err != nil {
			return
		}
		DBcheck, err = gorm.Open(Configure.DB0.Driver, Configure.DB0.Connection2)
		if err != nil {
			return
		}
		Dbid, err = gorm.Open(Configure.DB0.Driver, Configure.DB0.Connection3)
		if err != nil {
			return
		}
		DBcheck.SingularTable(true)
		Dbid.SingularTable(true)
	}

	// DB.DB().SetMaxOpenConns(200)
	DB0.SingularTable(true)
	return
}

func Clean() {
	if DB0 != nil {
		DB0.Close()
	}

	if Writer != nil {
		Writer.Close()
	}
}

func loadxmlsqlFile(fname string, data *xml2sql.Data) (*xml2sql.Data, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		return data, nil
	}

	if !strings.HasSuffix(finfo.Name(), ".xml") {
		return data, nil
	}
	newData, err := xml2sql.NewDataFromReaders([]io.Reader{f}, nil, data)
	if err != nil {
		return nil, err
	}
	return newData, nil
}

//加载sql
func LoadSqlTmpl() error {
	fmt.Println("Configure.SqlTmplPath", Configure.SqlTmplPath)
	f, err := os.Open(Configure.SqlTmplPath)
	if err != nil {
		return err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return nil
	}
	var data *xml2sql.Data = nil
	for _, v := range names {
		fname := path.Join(Configure.SqlTmplPath, v)
		data, err = loadxmlsqlFile(fname, data)
		if err != nil {
			return err
		}
	}
	Sqls = data.XML
	return nil
}

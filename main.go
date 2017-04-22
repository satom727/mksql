package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const SELECT byte = 1
const INSERT byte = 2
const UPDATE byte = 3

const MYSQL byte = 1
const SQLSERVER byte = 2

func main() {
	// csvで各テーブルの定義を一覧出力する
	// 出力したファイルをもとに各クエリを作成する
	uOpt := flag.String("u", "", "login user")
	pOpt := flag.String("p", "", "login password")
	dOpt := flag.String("d", "MySQL", "MySQL or SQLServer")
	sOpt := flag.String("s", "", "target scheme")
	tOpt := flag.String("t", "", "target table")
	qOpt := flag.String("q", "", "querry type")
	fOpt := flag.String("f", "", "file path")

	flag.Parse()

	opt := cmdOption{columns: []string{}, query: []string{}}
	optionList := map[string]string{}

	optionList["u"] = *uOpt
	optionList["p"] = *pOpt
	optionList["d"] = *dOpt
	optionList["s"] = *sOpt
	optionList["t"] = *tOpt
	optionList["q"] = *qOpt
	optionList["f"] = *fOpt

	fmt.Println(optionList["f"])
	//opt.init(flag.Args())
	opt.init(optionList)
	opt.makeQuery()
	//fmt.Println(opt.query[0])

}

//オプション構造体
type cmdOption struct {
	queryType byte
	DBType    byte
	password  string
	user      string
	file      string
	table     string
	scheme    string
	columns   []string
	query     []string
	updOpt    updateOption
	insOpt    insertOption
	selOpt    selectOption
}

//optionの値を格納
func (opt *cmdOption) init(optList map[string]string) {
	opt.file = getValue("f", optList)
	opt.user = getValue("u", optList)
	opt.password = getValue("p", optList)
	if strings.ToLower(getValue("d", optList)) == "mysql" {
		opt.DBType = MYSQL
	}
	if strings.ToLower(getValue("q", optList)) == "update" {
		opt.queryType = UPDATE
	}
	opt.table = getValue("t", optList)
	opt.scheme = getValue("s", optList)

}

// queryを作成
func (opt cmdOption) makeQuery() {
	file, err := os.Open(opt.file)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	lineCnt := 0
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true // ダブルクオートを厳密にチェックしない
	for row, err := reader.Read(); err != io.EOF; row, err = reader.Read() {
		lineCnt++
		if err != nil {
			panic(err)
		}
		if lineCnt == 1 {
			opt.columns = row
		} else {
			opt.makeUpdateQuery(row)
		}
	}
}

//update文を作成
func (opt *cmdOption) makeUpdateQuery(values []string) {
	sql := "update " + opt.table + " SET"
	for i := 0; i < len(values); i++ {
		sql = sql + " " + opt.columns[i] + "=" + values[i] + ","
	}
	sql = strings.Trim(sql, ",") + ";"
	opt.query = append(opt.query, sql)
	fmt.Println(sql)

}

func getValue(key string, list map[string]string) string {
	val, ok := list[key]
	if ok {
		return val
	}
	return ""
}

//func (opt *cmdOption) init(optArgs []string) {
//	for i := 0; i < len(optArgs); i++ {
//		if "d" == optArgs[i] {
//			i++
//			if "mysql" == ToLower(optArgs[i]） {
//				opt.DBType = MYSQL
//			}
//		} else if "p" == optArgs[i] {
//			i++
//			opt.password = optArgs[i]
//		} else if "u" == optArgs[i] {
//			i++
//			opt.user = optArgs[i]
//		} else if "f" == optArgs[i] {
//			i++
//			opt.file = optArgs[i]
//		} else if "s" == optArgs[i] {
//			i++
//			if "update" == optArgs[i] {
//				opt.queryType = UPDATE
//			}
//		} else {
//		}
//	}
//	if opt.queryType == UPDATE {
//		opt.setUpdOpt(updateOption{})
//	}
//}

func (opt *cmdOption) setUpdOpt(updOpt updateOption) {
	opt.updOpt = updOpt
}
func (opt *cmdOption) setInsOpt(insOpt insertOption) {
	opt.insOpt = insOpt
}
func (opt *cmdOption) setSelOpt(selOpt selectOption) {
	opt.selOpt = selOpt
}

type updateOption struct {
	where []string
}
type insertOption struct {
}
type selectOption struct {
}

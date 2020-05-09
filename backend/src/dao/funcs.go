package dao
import (
	"database/sql"
	"fmt"
	"strings"
)

//根据配置初始化DSN
func initDSN(config map[string]string) string {
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// 初始化服务器的连接
	// 完成参数的获取，考虑用户未传递的情况（默认值）
	DSN := ""
	username, ok := config["username"]
	if ! ok {
		username = ""
	}
	password, ok := config["password"]
	if ! ok {
		password = ""
	}
	if "" == username {
		DSN += ""
	} else if "" == password {
		// username 存在，但是 password不存在
		DSN += username + "@"
	} else {
		// username 和 password 都存在
		DSN += username + ":" + password + "@"
	}

	protocol, ok := config["protocol"]
	if ! ok {
		protocol = "tcp"
	}
	host, ok := config["host"]
	if ! ok {
		host = "localhost"
	}
	port, ok := config["port"]
	if ! ok {
		port = "3306"
	}
	DSN += protocol + "(" + host + ":" + port + ")"

	dbname, ok := config["dbname"]
	if ! ok {
		dbname = ""
	}
	DSN += "/" + dbname

	// 参数处理
	params := []string{}
	collation, ok := config["collation"]
	if ! ok {
		collation = ""
	} else {
		params = append(params, "collation=" + collation)
	}

	// 参数处理完毕后，拼凑成一个字符串，加载DSN后边即可
	paramStr := strings.Join(params, "&")
	DSN += "?" + paramStr

	return DSN
}

//构建select语句, 私有
func (this *Dao) buildSelect() string {
	query := "SELECT"
	// distinct
	if this.queryDistinct {
		query += " DISTINCT"
	}
	// 字段部分
	query += " " + this.queryFields
	// 表名
	query += " FROM " + this.queryTable
	// Join部分
	for _, join := range this.queryJoins {
		query += " " + join
	}
	//Where
	if "" != this.queryCondition {
		query += " WHERE " + this.queryCondition
	}
	// GroupBy部分
	if "" != this.queryGroupBy {
		query += " GROUP BY " + this.queryGroupBy
	}
	// Having 部分
	if "" != this.queryHavingCondition {
		query += " HAVING " + this.queryHavingCondition
	}
	// order by
	if "" != this.queryOrderBy {
		query += " ORDER BY " + this.queryOrderBy
	}
	// limit
	if "" != this.queryOffset && "" != this.queryLimit {
		query += " LIMIT " + this.queryLimit + " OFFSET " + this.queryOffset
	} else if "" != this.queryLimit {
		query += " LIMIT " + this.queryLimit
	}
	return query
}

//执行非查询类
func (this *Dao) exec(query string, params ...interface{}) (sql.Result, error) {
	// 记录
	this.sql = query
	this.sqlParams = params

	return this.db.Exec(query, params...)
}
//执行查询类
func (this *Dao) query(query string, params ...interface{}) (*sql.Rows, error){
	// 记录
	this.sql = query
	this.sqlParams = params

	return this.db.Query(query, params...)
}

// 单行和列
func (this *Dao) rowAndCols() (map[string]string, []string, error) {
	// 利用 Rows 获取多条，只要第一条
	this.Limit(1)

	rows, cols, err := this.rowsAndCols()
	if err != nil {
		return map[string]string{}, []string{}, err
	}

	// 判断是否查询到了数据
	if len(rows) > 0 {
		return rows[0], cols, nil
	} else {
		// 没有查询到数据的情况（一条都没有）
		return map[string]string{}, []string{}, nil
	}
}

//多行和列
func (this *Dao) rowsAndCols() ([]map[string]string, []string, error) {
	//一 构建SQL
	query := this.buildSelect()

	//二 执行
	//rows, err := this.db.Query(query, append(this.queryParams, this.queryHavingParams...)...)
	rows, err := this.query(query, append(this.queryParams, this.queryHavingParams...)...)
	// 执行过后，清理掉记录的query的每个部分
	this.clearQuery()
	if err != nil {
		return []map[string]string{}, []string{}, err
	}
	defer rows.Close()

	// 2.5 确定列的数量
	cols, colErr := rows.Columns()
	if colErr != nil {
		return []map[string]string{}, []string{}, colErr
	}
	colNum := len(cols)
	// 可以确定需要多少元素
	// 为了rows.Scan 传参，不定数量的参数，展开操作
	fields := make([]interface{}, colNum)
	// values 中存储的字符串值的引用！
	values := make([]sql.NullString, colNum)
	for i, _ := range fields {
		// fields 的每个元素都是指针类型，可以被赋值 *string
		fields[i] = &values[i]
	}

	result := []map[string]string{}
	// 三 处理结果
	for rows.Next() { // 确保存在记录
		// 获取记录数据
		// 变量为指针传递，保证可以修改 变量
		scanErr := rows.Scan(fields...)
		if scanErr != nil { // 当前记录scan错误，获取下一条记录
			fmt.Println(scanErr)
			continue
		}
		// 获取到数据，整理成目标格式
		row := map[string]string{}
		for i, _ := range fields {
			// 通过接口的断言在解析地址的方式得到字符串
			//row[cols[i]] = *(fields[i].(*string))
			ns := *(fields[i].(*sql.NullString))
			if ns.Valid {
				row[cols[i]] = ns.String
			} else {
				row[cols[i]] = "" // "NULL"
			}
		}
		result = append(result, row)
	}
	return result, cols, nil
}

// 清理query部分
func (this *Dao) clearQuery() {
	this.queryTable = ""
	this.queryCondition = ""
	this.queryParams = []interface{}{}
	this.queryDistinct = false
	this.queryFields = "*" // 非 string 的零值
	this.queryJoins = []string{}
	this.queryGroupBy = ""
	this.queryHavingCondition = ""
	this.queryHavingParams = []interface{}{}
	this.queryOrderBy = ""
	this.queryLimit = ""
	this.queryOffset = ""
}

// 私有的函数，表名的包裹
func _table(table string) string {
	tableSlice := strings.Split(table, " ")
	if tslen := len(tableSlice); tslen > 1 {
		return fmt.Sprintf("`%s` AS `%s`", tableSlice[0], tableSlice[tslen-1])
	} else {
		return fmt.Sprintf("`%s`", tableSlice[0])
	}
}

// 字段的反引号包裹
func _fieldWrap(field string) string {
	// 判断字段名是否有.部分，有的话分别包裹
	fieldSlice := strings.Split(field, ".")
	if fslen := len(fieldSlice); fslen > 1 {
		// t.field
		if fieldSlice[1] == "*" {
			return fmt.Sprintf("`%s`.%s", fieldSlice[0], fieldSlice[1])
		} else {
			return fmt.Sprintf("`%s`.`%s`", fieldSlice[0], fieldSlice[1])
		}
	} else {
		// field
		if fieldSlice[0] == "*" {
			return fmt.Sprintf("%s", fieldSlice[0])
		} else {
			return fmt.Sprintf("`%s`", fieldSlice[0])
		}
	}
}
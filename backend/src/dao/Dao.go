package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

// 原生字符串类型
type Raw struct {
	String string
}

type Dao struct {
	db *sql.DB
	sql string
	sqlParams []interface{}

	// query 的每个部分
	queryTable string
	queryCondition string
	queryParams []interface{}
	queryDistinct bool
	queryFields string
	queryJoins []string // join A.... join B...
	queryGroupBy string
	queryHavingCondition string
	queryHavingParams []interface{}
	queryOrderBy string
	queryLimit string
	queryOffset string
}

// select distinct
func (this *Dao) Distinct() *Dao {
	this.queryDistinct = true
	return this
}

// 表名
func (this *Dao) Table(table string) *Dao {
	this.queryTable = _table(table)
	// 返回当前对象
	return this
}
func (this *Dao) From(table string) *Dao {
	return this.Table(table)
}

// join, join tableA as A ON T.field = A.field
func (this *Dao) Join (joinType string, joinTable string, joinOn string) *Dao {
	join := fmt.Sprintf("%s JOIN %s ON %s", strings.ToUpper(joinType), _table(joinTable),  joinOn)

	this.queryJoins = append(this.queryJoins, join)
	return this
}
func (this *Dao) LeftJoin (joinTable string, joinOn string) *Dao {
	return this.Join("LEFT", joinTable, joinOn)
}
func (this *Dao) RightJoin (joinTable string, joinOn string) *Dao {
	return this.Join("RIGHT", joinTable, joinOn)
}
func (this *Dao) InnerJoin (joinTable string, joinOn string) *Dao {
	return this.Join("INNER", joinTable, joinOn)
}

// Group by 子句
func (this *Dao) GroupBy(fields ...string) *Dao {
	fieldSlice := []string{}
	for _, field := range fields {
		fieldSlice = append(fieldSlice, _fieldWrap(field))
	}
	this.queryGroupBy = strings.Join(fieldSlice, ", ")
	return this
}


// having条件
func (this *Dao) Having(condition string, params []interface{}) *Dao {
	this.queryHavingCondition, this.queryHavingParams = condition, params
	return this
}

//排序
func (this *Dao) OrderBy(fields ...string) *Dao {
	//直接将 排序的字段连接起来即可
	//this.queryOrderBy = strings.Join(fields, ", ")
	fieldSlice := []string{}
	//若需要为字段包裹反引号
	for _, field := range fields {
		// 判断是否指定了排序方式
		if fieldSplit := strings.Split(field, " "); len(fieldSplit) > 1 {
			fieldSlice = append(fieldSlice, _fieldWrap(fieldSplit[0]) + " " + strings.ToUpper(fieldSplit[1]))
		} else {
			fieldSlice = append(fieldSlice, _fieldWrap(fieldSplit[0]))
		}
	}
	this.queryOrderBy = strings.Join(fieldSlice, ", ")
	return this
}

// limit
func (this *Dao) Limit(size int) *Dao {
	//strconv.Itoa(), 整型转换为字符串
	this.queryLimit = strconv.Itoa(size)
	return this
}
// Offset
func (this *Dao) Offset(rows int) *Dao {
	this.queryOffset = strconv.Itoa(rows)
	return this
}

// 字段部分
func (this *Dao) Fields(fields ...interface{}) *Dao {
	fieldList := []string{}
	for _, field := range fields {
		// 考虑每个字段的各种情况
		switch field.(type) {
		case string:
			// 是否包含空格，有空格意味着有别名
			fieldSlice := strings.Split(field.(string), " ") // []string
			if fslen := len(fieldSlice); fslen > 1 {
				// 存在空格，别名问题，全部处理成 field as alisa
				fieldList = append(fieldList, fmt.Sprintf("%s AS `%s`", _fieldWrap(fieldSlice[0]), fieldSlice[fslen-1]))
			} else {
				// 没有别名，处理字段即可
				fieldList = append(fieldList, _fieldWrap(fieldSlice[0]))
			}
		case Raw:
			fieldList = append(fieldList, field.(Raw).String)
		}
	}

	this.queryFields = strings.Join(fieldList, ", ")
	return this
}


//获取单列
func (this *Dao) Column() ([]string, error) {
	rows, cols, err := this.rowsAndCols()
	if err != nil {
		return []string{}, err
	}
	firstCol := cols[0]
	result := []string{}
	for _, row := range rows {
		result = append(result, row[firstCol])
	}
	return result, nil
}

//fetchRow
func (this *Dao) Row() (map[string]string, error) {
	result, _, err := this.rowAndCols()
	return result, err
}

// 单个值
func (this *Dao) Value() (string , error) {
	row, cols, err := this.rowAndCols()
	if err != nil {
		return "", err
	}
	// 确定第一列的字段名
	firstCol := cols[0]

	if len(row) > 0 {
		// 返回第一个即可
		return row[firstCol], nil
	} else {
		return "", nil
	}

	return "", nil
}

// fetchAll
func (this *Dao) Rows() ([]map[string]string, error) {
	result, _, err := this.rowsAndCols()
	return result, err
}


//条件
func (this *Dao) Where(condition string, params []interface{}) *Dao {
	this.queryCondition, this.queryParams = condition, params
	return this
}

// 删除，delete from table where condition
func (this *Dao) Delete() (int64, error) {

	//拼凑SQL
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s",
		this.queryTable,
		this.queryCondition,
	)

	//执行SQL, 数据部分由 字段和条件参数部分组成
	//result, err := this.db.Exec(query, this.queryParams...)
	result, err := this.exec(query, this.queryParams...)
	// 执行过后，清理掉记录的query的每个部分
	this.clearQuery()
	if err != nil {
		return 0, err
	}
	// 执行成功，返回 AffectedRows
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

// 更新
func (this *Dao) Update(fields map[string]interface{}) (int64, error) {
	// 准备工作
	setList, valueList := []string{}, []interface{}{}
	for field, value := range fields {
		setList = append(setList, fmt.Sprintf("`%s` = ?", field))
		valueList = append(valueList, value)
	}

	// 拼凑SQL
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		this.queryTable,
		strings.Join(setList, ", "),
		this.queryCondition,
	)

	//执行SQL, 数据部分由 字段和条件参数部分组成
	result, err := this.exec(query, append(valueList, this.queryParams...)...)
	// 执行过后，清理掉记录的query的每个部分
	this.clearQuery()
	if err != nil {
		return 0, err
	}
	// 执行成功，返回 AffectedRows
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

// 插入
func (this *Dao) Insert(fields map[string]interface{}) (int64, error) {

	// 准备工作，遍历fields，得到全部的 key（就是字段） 和 value
	fieldList, valueList, valuePLList := []string{}, []interface{}{}, []string{}
	for field, value := range fields {
		// 字段名 使用反引号包裹
		fieldList = append(fieldList, "`" + field + "`")
		valueList = append(valueList, value)

		//值的占位符
		valuePLList = append(valuePLList, "?")
	}

	// 一 拼凑insert
	// 可以升级为 fmt.Sprintf(), 格式化字符串（不输出），返回
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		this.queryTable,
		strings.Join(fieldList, ", "),
		strings.Join(valuePLList, ", "),
	)

	// 二 执行
	result, err := this.exec(query, valueList...)
	// 执行过后，清理掉记录的query的每个部分
	this.clearQuery()
	if err != nil {
		return 0, err
	}
	// 执行成功，返回lastInsertID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return id, nil
}

// 获取SQL和参数
func (this *Dao) LastSQL () (string, []interface{}) {
	// 一次性获取
	sql, params := this.sql, this.sqlParams
	this.sql, this.sqlParams = "", []interface{}{}

	return sql, params
}

// 构造函数
func NewDao(config map[string]string) (*Dao, error) {
	// 一：拼凑DSN
	DSN := initDSN(config)

	//二 ：连接MySQL
	db, err := sql.Open("mysql", DSN) // sql. 抽象层
	if err != nil {
		return nil, err
	}
	// 测试连接
	if pingErr:=db.Ping(); pingErr!=nil {
		return nil,pingErr
	}

	// 三： 返回Dao对象
	this := new (Dao)
	this.db = db
	// 初始化 query 部分
	this.clearQuery()
	return this, nil
}

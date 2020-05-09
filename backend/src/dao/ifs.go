package dao
//Dao接口
type I_Dao interface {
	Distinct() *Dao
	Fields(fields ...interface{}) *Dao
	Table(string) *Dao
	From(string) *Dao
	Join (string, string, string) *Dao
	LeftJoin (string, string) *Dao
	RightJoin (string, string) *Dao
	InnerJoin (string, string) *Dao
	Where(string, []interface{}) *Dao
	GroupBy(...string) *Dao
	Having(string, []interface{}) *Dao
	OrderBy(...string) *Dao
	Limit(int) *Dao
	Offset(int) *Dao

	Column() ([]string, error)
	Row() (map[string]string, error)
	Value() (string , error)
	Rows() ([]map[string]string, error)
	Delete() (int64, error)
	Update(map[string]interface{}) (int64, error)
	Insert(map[string]interface{}) (int64, error)

	LastSQL() (string, []interface{})
}
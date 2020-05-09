package config
// config.App
var App map[string]string

func InitConfig() {
	App = map[string]string{
		"DB_DRIVER": "mysql",
		"DB_TABLE_PREFIX": "a_",
		"MYSQL_HOST": "localhost",
		"MYSQL_PORT": "3306",
		"MYSQL_USER": "projectAUser",
		"MYSQL_PASSWORD": "peiqingtao",
		"MYSQL_DBNAME": "projectA",
		"MYSQL_CHARSET": "utf8mb4",
		"MYSQL_LOC": "Local",
		"MYSQL_PARSETIME": "true", // false默认，不执行time的解析

		// redis 配置
		//"REDIS_HOST": "192.168.2.102",
		//"REDIS_PORT": "6379",
		//"REDIS_DB": "0",
		//"REDIS_PASSWORD": "",
		//
		//"UPLOAD_PATH": "D:\\projects\\class\\2th\\projectA\\upload\\", // 需要额外配置
		//"IMAGE_HOST": "http://localhost:8089/",
		//
		//"THUMB_SMALL_W": "146",
		//"THUMB_SMALL_H": "146",
		//"THUMB_BIG_W": "1460",
		//"THUMB_BIG_H": "1460",
		//
		"SERVER_ADDR": ":8088",
		//
		//"SECRET": "AAAAB3NzaC1yc2EAAAABJQAAACEAwchMWdcW8/z1lju6s/PFEbNokNWRgmzpuRGcx02wHn8",
	}
}
package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"model"
)

//添加分类
func CategoryAdd(c *gin.Context){
	//得到一个Category模型对象
	category := model.Category{}
	//从请求中解析数据
	err := c.ShouldBind(&category)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	//自动解析绑定成功
	log.Println(category)
	//数据入库
	orm.Create(&category)
	
	log.Println(category)
	// 完成响应
	c.JSON(200, gin.H{
		"error": "",
		"data": category,
	})
}
//分类树
func CategoryTree(c *gin.Context){
	//查询全部的分类
	var categories []model.Category
	//fmt.Println(db.HasTable(model.Category{}))
	orm.Unscoped().Find(&categories)
	// 遍历categories，得到每个分类，利用分类查询关联
	//for i,_:= range categories{
	//	orm.Model(&categories[i]).//确定使用的查询表
	//		Related(&categories[i].Products)//关联查询
	//}

	//响应json
	c.JSON(200,gin.H{
		"error":"",
		"data":categories,
	})


}

////分类树
//func CategoryTree(c *gin.Context){
//
//	//连接数据库，获取全部的分类内容
//	config:=map[string]string{
//		"username":"projectAUser",
//		"password":"peiqingtao",
//		"host":"127.0.0.1",
//		"port":"3306",
//		"dbname":"projectA",
//		"collation":"utf8mb4_general_ci",
//	}
//	db,err := dao.NewDao(config)
//	if err!= nil {
//		fmt.Println(err)
//		return
//	}
//	rows,err :=db.Table("a_categories").Rows()
//	if err==nil {
//		//响应json数据
//		c.JSON(200,gin.H{
//			"data":rows,
//			"error":"",
//		})
//	}else {
//		c.JSON(200,gin.H{
//			"error":err.Error(),
//		})
//	}
//}
func CategoryDelete(){

}

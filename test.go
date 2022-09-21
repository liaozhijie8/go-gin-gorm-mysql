package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// import "fmt"

func main() {
	dsn := "root:zxc43217@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	fmt.Println(db)
	fmt.Println(err)
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20);not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20);not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20);not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40);not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200);not null" json:"address" binding:"required"`
	}
	db.AutoMigrate(&List{})
	r := gin.Default()
	/* 增加数据 */
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)
		// 判断
		if err != nil {
			c.JSON(200, gin.H{
				"msg": "添加失败",
				"data":    gin.H{},
				"code":    400,
			})
		} else {
			// 添加到数据库
			db.Create(&data)
			c.JSON(200, gin.H{
				"msg": "添加成功",
				"data":    data,
				"code":    200,
			})
		}
	})
	/* 删除数据 */
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List
		id := c.Param("id")
		// 判断id是否存在
		db.Where("id=?", id).Find(&data)
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id没有找到,删除失败",
				"code": 400,
			})
		} else {
			// 操作数据删除
			db.Where("id=?", id).Delete(&data)
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}

	})
	/* 修改数据 */
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data List
		id := c.Param("id")
		// 判断id是否存在
		db.Select("id").Where("id=?", id).Find(&data)
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				db.Where("id=?", id).Updates(&data)
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})
	/* 查询数据 */
	r.GET("/user/list/:name", func(c *gin.Context) {
		// 获取路径参数
		name := c.Param("name")
		var dataList []List
		// 查询数据
		db.Where("name=?", name).Find(&dataList)
		// 判断是否查到数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": dataList,
			})
		}
	})
	// 查询全部
	r.GET("/user/list", func(c *gin.Context) {
		var dataList []List
		// 1.查询全部数据，查询分页数据
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))
		// 判断是否需要分页
		if pageNum == 0 {
			pageNum = -1
		}
		if pageSize == 0 {
			pageSize = -1
		}

		// 分页的固定写法
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		var total int64
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})
	// 端口号
	PORT := "3006"
	r.Run(":" + PORT) // 监听并在 0.0.0.0:8080 上启动服务

}

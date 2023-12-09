package dao

import (
	"github.com/CocaineCong/micro-todoList/app/user/repository/db/model"
)

// migration 数据库迁移，不用手动创建表，自动映射
// 类似于 Java中的 @Table 注解
func migration() {
	_db.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&model.User{})
}

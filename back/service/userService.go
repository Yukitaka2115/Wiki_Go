package service

import (
	"errors"
	"log"
	"wiki/dao"
	"wiki/model"
)

func AddUser(user model.User) {
	err := dao.Db.AutoMigrate(&model.User{})
	if err != nil {
		return
	}
	dao.Db.Create(&user)
} //注册添加用户方法，身份选择不开放，预先添加admin用户

func GetUsers(user []model.User) ([]model.User, error) {
	res := dao.Db.Find(&user)
	if res.Error != nil {
		log.Println("Failed to find users:", res.Error)
		return nil, res.Error
	}
	return user, nil
} //全部用户，仅针对管理员

func GetUserByUserNameAndPwd(username string, pwd string) model.User {
	var user model.User
	dao.Db.First(&user, username, pwd)
	return user
} //登录确定唯一

func UpdateUserInfo(id int, newInfo model.User) (model.User, error) {
	var user model.User
	if err := dao.Db.First(&user, id).Error; err != nil {
		return model.User{}, err
	}
	if user.ID == 0 {
		return model.User{}, errors.New("页面不存在")
	}
	user = newInfo
	if err := dao.Db.Model(&user).Save(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
} //用户/管理员修改字段，身份修改不开放

func DeleteUserByID(id int) error {
	var user model.User
	dao.Db.First(&user, id)
	dao.Db.Delete(&user)
	return nil
} //用户注销/管理员删除用户

func IsUsernameExists(username string) bool {
	// 这里需要根据您的数据库类型和使用的数据库操作包来执行查询操作
	// 假设您使用的是 GORM 这样的 ORM 框架
	var user model.User // 获取数据库连接
	if err := dao.Db.Where("username = ?", username).First(&user).Error; err != nil {
		return false // 没有找到用户，用户名可用
	}
	return true // 找到了用户，用户名已存在
}

func GetUserByID(id int) (model.User, error) {
	var user model.User
	// 根据用户ID从数据库中检索用户信息
	if err := dao.Db.First(&user, id).Error; err != nil {
		// 如果查询出错，返回错误
		return model.User{}, err
	}
	// 如果成功找到用户，返回用户信息和nil错误
	return user, nil
}

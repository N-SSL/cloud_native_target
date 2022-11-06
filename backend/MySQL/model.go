package MySQL

import "gorm.io/gorm"

type EnvStruct struct {
	ID           int64    `gorm:"primaryKey"` //自增查序列号
	EnvName      string   	  //镜像名
	Type         string						  //镜像类型
	Description  string                       //镜像描述
	DeployFile   []byte						  //k8s deployment
	ServiceFile  []byte                       //k8s service
	gorm.Model
}

type RunningStruct struct {
	ID           int64    `gorm:"primaryKey"` //自增序列号
	RunningName  string   `gorm:"unique"`	  //运行名，环境名 + 新增ID
	EnvName      string						  //环境名
	SessionID    string						  //用户浏览器cookie
	ClusterIP    string                       //绑定ClusterIP
	Port         int32						  //绑定端口号
	BindID       string						  //新增ID
	Url          string						  //新增ID
	gorm.Model
}

type Users struct {
	ID           int64    `gorm:"primaryKey"` //自增查序列号
	Username   	 string   `gorm:"type:varchar(50);not_null;unique" json:"username"`
	Password 	 []byte   `gorm:"type:varchar(200);not_null" json:"-"`
	Role  	     string   `json:"role"`
	gorm.Model
}


type Login struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
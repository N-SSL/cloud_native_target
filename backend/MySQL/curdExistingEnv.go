package MySQL

import (
	"github.com/N-SSL/container-target/controllers"
	"github.com/N-SSL/container-target/k8s"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListAllImages(c *gin.Context)  {
	var EnvList []EnvStruct
	SqlDB.Find(&EnvList)
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: EnvList,
	})
}

func GetImageByName(c *gin.Context)  {
	name := c.Param("imageName")
	var Env EnvStruct
	SqlDB.Where(&EnvStruct{EnvName: name}).First(&Env)
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: Env,
	})
}

func AddNewEnv(c *gin.Context)  {
	name              := c.PostForm("Name")
	Type              := c.PostForm("Type")
	Description       := c.PostForm("Description")
	DeployFile, derr  := c.FormFile("DeployFile")
	ServiceFile, serr := c.FormFile("ServiceFile")
	if name ==""{
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: "need Env Name",
			Data: nil,
		})
		return
	}
	if derr != nil {
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message:"cannot get deploy file",
			Data: nil,
		})
		return
	}
	if serr != nil {
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message:"cannot get service file",
			Data: nil,
		})
		return
	}
	deployData, derr1 := k8s.VerifyDeploy(DeployFile)
	serviceData, serr1 := k8s.VerifyService(ServiceFile)
	if derr1!= nil {
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: "cannot convert deploy file",
			Data: nil,
		})
		return
	}
	if serr1 != nil{
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: "cannot convert service file",
			Data: nil,
		})
		return
	}
	var newEnv = EnvStruct{
		EnvName: name,
		Type: Type,
		DeployFile: deployData,
		ServiceFile: serviceData,
		Description: Description,
	}
	if err := SqlDB.Create(&newEnv).Error; err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: `name " ` + name +` " already in use `,
			Data: nil,
		})
	}else {
		c.JSON(http.StatusOK,controllers.Response{
			Success: true,
			Code: http.StatusOK,
			Message: "OK",
			Data: nil,
		})
	}
}

func DeleteEnv(c *gin.Context)  {
	name := c.PostForm("Name")
	if name == "" {
		c.JSON(http.StatusBadRequest, controllers.Response{
				Success: false,
				Code: http.StatusBadRequest,
				Message: "name is required",
				Data: nil,
		})
		return
	}
	if SqlDB.Where("env_name = ?", name).Delete(&EnvStruct{}).Error != nil {
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: `no such environment: " `+name+` " or it has been deleted`,
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "delete successfully",
		Data: nil,
	})
}


func UpdateEnv(c *gin.Context)  {
	name              := c.PostForm("Name")
	Type              := c.PostForm("Type")
	Description       := c.PostForm("Description")
	DeployFile, derr  := c.FormFile("DeployFile")
	ServiceFile, serr := c.FormFile("ServiceFile")
	if name == "" {
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: false,
			Code: http.StatusBadRequest,
			Message: "need Env Name",
			Data: nil,
		})
		return
	}
	var updateEnv EnvStruct

	if Type != ""{
		updateEnv.Type = Type
	}
	if Description !=""{
		updateEnv.Description = Description
	}
	if derr == nil {
		deployData, derr1 := k8s.VerifyDeploy(DeployFile)
		if derr1 != nil{
			c.JSON(http.StatusBadRequest,controllers.Response{
				Success: false,
				Code: http.StatusBadRequest,
				Message: "cannot convert deploy file",
				Data: nil,
			})
			return
		}
		updateEnv.DeployFile = deployData

	}
	if serr == nil {
		serviceData, serr1 := k8s.VerifyService(ServiceFile)
		if serr1 != nil{
			c.JSON(http.StatusBadRequest,controllers.Response{
				Success: false,
				Code: http.StatusBadRequest,
				Message: "cannot convert service file",
				Data: nil,
			})
			return
		}
		updateEnv.ServiceFile = serviceData
	}
	result := SqlDB.Model(EnvStruct{}).Where("env_name = ?", name).Updates(updateEnv)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "server error",
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusCreated,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "update successfully",
		Data: nil,
	})
	return
}




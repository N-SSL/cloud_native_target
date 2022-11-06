package MySQL

import (
	"github.com/N-SSL/container-target/controllers"
	"github.com/N-SSL/container-target/k8s"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func StartEnv(c *gin.Context)  {
	name := c.PostForm("Name")
	var envStart EnvStruct
	err0 := SqlDB.Where(&EnvStruct{EnvName: name}).First(&envStart).Error
	if envStart.EnvName == "" || err0 != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "search error",
			Data: nil,
		})
		return
	}
	userID, _ := c.Get("id")
	var running RunningStruct
	SqlDB.Where(RunningStruct{SessionID: userID.(string) , EnvName: name}).First(&running)
	if running.SessionID == userID {
		c.JSON(http.StatusOK,controllers.Response{
			Success: true,
			Code: http.StatusOK,
			Message: "OK",
			Data: gin.H{
				"data": "You have already Started the env",
				"deployName": running.RunningName,
				"BindID": running.BindID,
				"EnvName":envStart.EnvName,
				"Description":envStart.Description,
				"url": "/app/" + running.BindID + "/",
			},
		})
		return
	}
	randomString := RandomString(16)
	extName := strings.Replace(strings.ToLower(name), "_", "", -1) + "-" + randomString
	ip, port, err := k8s.StartEnvNow(envStart.DeployFile,envStart.ServiceFile,extName)
	if err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "cannot establish env",
			Data: nil,
		})
		return
	}
	var newSVC = RunningStruct{
		RunningName: extName,
		EnvName: name,
		BindID: randomString,
		Url: "/app/" + randomString,
		SessionID: userID.(string),
		ClusterIP: ip,
		Port: port,
	}

	if err := SqlDB.Create(&newSVC).Error; err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "mysql error",
			Data: nil,
		})
		return
	}
	log.Println(ip)
	log.Println(port)
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"data": "success",
			"deployName":extName,
			"BindID": randomString,
			"EnvName":envStart.EnvName,
			"Description":envStart.Description,
			"url": "/app/" + randomString + "/",
		},
	})
}

func EndEnv(c *gin.Context)  {
	name := c.PostForm("Name")
	userID, _ := c.Get("id")
	err := endEnvBySessionIdAndEnvName(userID.(string),name)
	if err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: true,
			Code: http.StatusInternalServerError,
			Message: "cannot end that env or that env had been ended",
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: nil,
	})
	return
}

func RestartEnv(c *gin.Context)  {
	name := c.PostForm("Name")
	userID, _ := c.Get("id")
	endErr := endEnvBySessionIdAndEnvName(userID.(string),name)
	if endErr !=nil{
		log.Println(endErr)
	}
	var envStart EnvStruct
	err0 := SqlDB.Where(&EnvStruct{EnvName: name}).First(&envStart).Error
	if envStart.EnvName == "" || err0 != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "search error",
			Data: nil,
		})
		return
	}
	randomString := RandomString(16)
	extName := strings.Replace(strings.ToLower(name), "_", "", -1) + "-" + randomString
	ip, port, err := k8s.StartEnvNow(envStart.DeployFile,envStart.ServiceFile,extName)
	if err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "cannot establish env",
			Data: nil,
		})
		return
	}
	var newSVC = RunningStruct{
		RunningName: extName,
		EnvName: name,
		BindID: randomString,
		Url: "/app/" + randomString,
		SessionID: userID.(string),
		ClusterIP: ip,
		Port: port,
	}
	if err := SqlDB.Create(&newSVC).Error; err != nil{
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: false,
			Code: http.StatusInternalServerError,
			Message: "mysql error",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: gin.H{
			"data": "success",
			"deployName":extName,
			"BindID": randomString,
			"EnvName":envStart.EnvName,
			"Description":envStart.Description,
			"url": "/app/" + randomString + "/",
		},
	})
}

func endEnvBySessionIdAndEnvName(SessionId string, EnvName string) error {
	var running RunningStruct
	SqlDB.Where(&RunningStruct{EnvName: EnvName,SessionID: SessionId}).First(&running)
	err := k8s.EndEnv(running.RunningName)
	if err != nil{
		return err
	}
	SqlDB.Where(&RunningStruct{EnvName: EnvName,SessionID: SessionId}).Delete(&running)
	return nil
}


func EndEnvForce(c *gin.Context)  {
	name := c.PostForm("Name")
	var envRunning RunningStruct
	err := k8s.EndEnv(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError,controllers.Response{
			Success: true,
			Code: http.StatusInternalServerError,
			Message: "cannot end that env or that env had been ended",
			Data: nil,
		})
		return
	}
	SqlDB.Where(&RunningStruct{RunningName: name}).Delete(&envRunning)
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: nil,
	})
	return
}

func ListAllRunningEnv(c *gin.Context)  {
	var envRunningList []RunningStruct
	// SqlDB.Find(&envRunningList)
	// c.JSON(http.StatusOK,controllers.Response{
	// 	Success: true,
	// 	Code: http.StatusOK,
	// 	Message: "OK",
	// 	Data: envRunningList,
	// })
	// return
	err := SqlDB.Find(&envRunningList).Error
	if err != nil{
		c.JSON(http.StatusBadRequest,controllers.Response{
			Success: true,
			Code: http.StatusBadRequest,
			Message: "cannot find any running env",
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK,controllers.Response{
		Success: true,
		Code: http.StatusOK,
		Message: "OK",
		Data: envRunningList,
	})
	return

}
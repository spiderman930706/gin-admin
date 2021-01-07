package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spiderman930706/gin-admin/api"
)

func InitRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("table")
	{
		apiRouter.GET("", api.GetAdminTableList)
		apiRouter.GET("/:table", api.GetAdminTableData)
		apiRouter.GET("/:table/:data_id", api.GetAdminDataDetail)
		apiRouter.POST("/:table", api.NewAdminData)
		apiRouter.PUT("/:table/:data_id", api.ChangeAdminData)
		apiRouter.DELETE("/:table/:data_id", api.DeleteAdminData)
	}
}

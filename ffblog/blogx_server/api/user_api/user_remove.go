package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (UserApi) UserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.RemoveRequest](c)

	var list []models.UserModel
	global.DB.Find(&list, "id in ?", cr.IDList)

	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除失败", c)
			return
		}
	}

	res.OkWithMsg(fmt.Sprintf("删除成功 成功删除%d条", len(list)), c)
}

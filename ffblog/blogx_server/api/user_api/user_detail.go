// api/user_api/user_detail.go
package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type UserDetailResponse struct {
	ID             uint                    `json:"id"`
	CreatedAt      time.Time               `json:"createdAt"`
	Username       string                  `json:"username"`
	Nickname       string                  `json:"nickname"`
	Avatar         string                  `json:"avatar"`
	Abstract       string                  `json:"abstract"`
	RegisterSource enum.RegisterSourceType `json:"registerSource"` // 注册来源
	CodeAge        int                     `json:"codeAge"`        // 码龄
	Role           enum.RoleType           `json:"role"`           // 角色
	models.UserConfModel
	Email       string `json:"email"`
	UsePassword bool   `json:"usePassword"`
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var user models.UserModel
	err := global.DB.Preload("UserConfModel").Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	var data = UserDetailResponse{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Abstract:       user.Abstract,
		Role:           user.Role,
		RegisterSource: user.RegisterSource,
		CodeAge:        user.CodeAge(),
		Email:          user.Email,
	}
	if user.Password != "" {
		data.UsePassword = true
	}
	if user.UserConfModel != nil {
		data.UserConfModel = *user.UserConfModel
	}

	res.OkWithData(data, c)
}

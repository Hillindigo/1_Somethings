package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils"
	"blogx_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleTagOptionsView(c *gin.Context) {
	claims := jwts.GetClaims(c)

	var articleList []models.ArticleModel
	global.DB.Find(&articleList, "user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished)

	var tagList ctype.List
	for _, model := range articleList {
		tagList = append(tagList, model.TagList...)
	}
	tagList = utils.Unique(tagList)
	var list = make([]models.OptionsResponse[string], 0)
	for _, s := range tagList {
		list = append(list, models.OptionsResponse[string]{
			Label: s,
			Value: s,
		})
	}
	res.OkWithData(list, c)
}

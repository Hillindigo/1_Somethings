package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"blogx_server/utils/jwts"
	"blogx_server/utils/markdown"
	"blogx_server/utils/xss"
	"github.com/gin-gonic/gin"
)

type ArticleCreateRequest struct {
	Title       string             `json:"title" binding:"required"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content" binding:"required"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"required,oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateRequest](c)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	if global.Config.Site.SiteInfo.Mode == 2 {
		if user.Role != enum.AdminRole {
			res.FailWithMsg("博客模式下，普通用户不能发文章", c)
			return
		}
	}

	log := log_service.GetLog(c)
	log.ShowResponse()
	log.SetItem("文章标题", cr.Title)
	log.SetItem("发布用户", user.Username)

	// 判断分类id是不是自己创建的
	var category models.CategoryModel
	if cr.CategoryID != nil {
		err = global.DB.Take(&category, "id = ? and user_id = ?", *cr.CategoryID, user.ID).Error
		if err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}

	// 文章正文防xss注入
	cr.Content = xss.XSSFilter(cr.Content)

	// 如果不传简介，那么从正文中取前30个字符
	if cr.Abstract == "" {
		abs, err1 := markdown.ExtractContent(cr.Content, 100)
		if err1 != nil {
			log.SetError("正文解析错误", err1)
			res.FailWithMsg("正文解析错误", c)
			return
		}
		cr.Abstract = abs
	}
	// 正文内容图片转存
	// 1.图片过多，同步做，接口耗时高  异步做，

	var article = models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		UserID:      user.ID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		OpenComment: cr.OpenComment,
		CategoryID:  cr.CategoryID,
		Status:      cr.Status,
	}
	if cr.Status == enum.ArticleStatusExamine && global.Config.Site.Article.NoExamine {
		article.Status = enum.ArticleStatusPublished
	}

	err = global.DB.Create(&article).Error
	if err != nil {
		log.SetError("文章创建失败", err)
		res.FailWithMsg("文章创建失败", c)
		return
	}

	res.OkWithMsg("文章创建成功", c)
}

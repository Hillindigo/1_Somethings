// flags/flag_es.go
package flags

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/es_service"
	"github.com/sirupsen/logrus"
)

func EsIndex() {
	if global.ESClient == nil {
		logrus.Warnf("未开启es连接")
		return
	}
	article := models.ArticleModel{}
	es_service.CreateIndexV2(article.Index(), article.Mapping())
	text := models.TextModel{}
	es_service.CreateIndexV2(text.Index(), text.Mapping())
}

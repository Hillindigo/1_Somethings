package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
	"time"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()
	// 2 -> 8 -> 9
	//rootComment := GetRootComment(9)
	//fmt.Println(rootComment.ID)
	//rootComment = GetRootComment(8)
	//fmt.Println(rootComment.ID)
	//rootComment = GetRootComment(2)
	//fmt.Println(rootComment.ID)

	//model := models.CommentModel{
	//	Model: models.Model{ID: 2},
	//}
	//GetCommentTree(&model)

	//model := GetCommentTreeV3(2)
	//
	//fmt.Println(model.ID)
	//for _, c1 := range model.SubCommentList {
	//	fmt.Println("  ", c1.ID)
	//	for _, c2 := range c1.SubCommentList {
	//		fmt.Println("    ", c2.ID)
	//		for _, c3 := range c2.SubCommentList {
	//			fmt.Println("      ", c3.ID)
	//		}
	//	}
	//}

	//commentList := GetCommentOneDimensional(2)
	//for _, model := range commentList {
	//	fmt.Println(model.ID)
	//}
	//fmt.Println("======")
	//for _, model := range commentList[1:] {
	//	fmt.Println(model.ID)
	//}

	//res := GetCommentTreeV4(2)
	//byteData, _ := json.Marshal(res)
	//fmt.Println(string(byteData))

	lis := GetParents(10)
	for _, li := range lis {
		fmt.Println(li.ID)
	}
}

// GetRootComment 获取一个评论的根评论
func GetRootComment(commentID uint) (model *models.CommentModel) {
	var comment models.CommentModel
	err := global.DB.Take(&comment, commentID).Error
	if err != nil {
		return nil
	}
	if comment.ParentID == nil {
		// 没有父评论了，那么他就是根评论
		return &comment
	}
	return GetRootComment(*comment.ParentID)
}

// GetParents 获取一个评论的所有父评论
func GetParents(commentID uint) (list []models.CommentModel) {
	var comment models.CommentModel
	err := global.DB.Take(&comment, commentID).Error
	if err != nil {
		return
	}
	list = append(list, comment)
	if comment.ParentID != nil {
		// 没有父评论了，那么他就是根评论
		list = append(list, GetParents(*comment.ParentID)...)
	}
	return
}

// 查一个评论下的子评论
// 评论树
func GetCommentTree(model *models.CommentModel) {
	global.DB.Preload("SubCommentList").Take(model)
	for _, commentModel := range model.SubCommentList {
		GetCommentTree(commentModel)
	}
}

func GetCommentTreeV3(id uint) (model *models.CommentModel) {

	model = &models.CommentModel{
		Model: models.Model{ID: id},
	}

	global.DB.Preload("SubCommentList").Take(model)
	for i := 0; i < len(model.SubCommentList); i++ {
		commentModel := model.SubCommentList[i]
		item := GetCommentTreeV3(commentModel.ID)
		model.SubCommentList[i] = item
	}
	return
}

type CommentResponse struct {
	ID           uint               `json:"id"`
	CreatedAt    time.Time          `json:"createdAt"`
	Content      string             `json:"content"`
	UserID       uint               `json:"userID"`
	UserNickname string             `json:"userNickname"`
	UserAvatar   string             `json:"userAvatar"`
	ArticleID    uint               `json:"articleID"`
	ParentID     *uint              `json:"parentID"`
	DiggCount    int                `json:"diggCount"`
	ApplyCount   int                `json:"applyCount"`
	SubComments  []*CommentResponse `json:"subComments"`
}

func GetCommentTreeV4(id uint) (res *CommentResponse) {
	model := &models.CommentModel{
		Model: models.Model{ID: id},
	}

	global.DB.Preload("UserModel").Preload("SubCommentList").Take(model)
	res = &CommentResponse{
		ID:           model.ID,
		CreatedAt:    model.CreatedAt,
		Content:      model.Content,
		UserID:       model.UserID,
		UserNickname: model.UserModel.Nickname,
		UserAvatar:   model.UserModel.Avatar,
		ArticleID:    model.ArticleID,
		ParentID:     model.ParentID,
		DiggCount:    model.DiggCount,
		ApplyCount:   0,
		SubComments:  make([]*CommentResponse, 0),
	}
	for _, commentModel := range model.SubCommentList {
		res.SubComments = append(res.SubComments, GetCommentTreeV4(commentModel.ID))
	}
	return
}

// GetCommentOneDimensional 评论一维化
func GetCommentOneDimensional(id uint) (list []models.CommentModel) {
	model := models.CommentModel{
		Model: models.Model{ID: id},
	}

	global.DB.Preload("SubCommentList").Take(&model)
	list = append(list, model)
	for _, commentModel := range model.SubCommentList {
		subList := GetCommentOneDimensional(commentModel.ID)
		list = append(list, subList...)
	}
	return
}

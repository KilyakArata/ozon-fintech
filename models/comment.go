package models

type Comment struct {
	ID              string `json:"id"`
	Author          string `json:"author"`
	Body            string `json:"body"`
	PostId          string `json:"post_id"`
	ParentCommentId string `json:"parent_comment_Id"`
}

package model

type CommentHotModel struct {
	CommentId int64 `gorm:"column:comment_id" json:"comment_id"`
	PostID    int64 `gorm:"column:post_id" json:"post_id"`
	RootID    int64 `gorm:"column:root_id" json:"root_id"`
	ParentID  int64 `gorm:"column:parent_id" json:"parent_id"`
	UserID    int64 `gorm:"column:user_id" json:"user_id"`
	Score     int   `gorm:"column:score" json:"score"`
	DelFlag   int   `gorm:"column:del_flag" json:"del_flag"`
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *CommentHotModel) TableName() string {
	return "comment_hot"
}

package model

type UserCommentModel struct {
	ID        int64 `gorm:"column:id" json:"id"`
	CommentID int64 `gorm:"column:comment_id" json:"comment_id"`
	UserID    int64 `gorm:"column:user_id" json:"user_id"`
	DelFlag   int   `gorm:"column:del_flag" json:"del_flag"`
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserCommentModel) TableName() string {
	return "user_comments"
}

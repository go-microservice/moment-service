package model

type UserPostModel struct {
	ID        int64 `gorm:"column:id" json:"id"`
	PostID    int64 `gorm:"column:post_id" json:"post_id"`
	UserID    int64 `gorm:"column:user_id" json:"user_id"`
	DelFlag   int   `gorm:"column:del_flag" json:"del_flag"`
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserPostModel) TableName() string {
	return "user_posts"
}

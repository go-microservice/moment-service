package model

type CommentInfoModel struct {
	ID         int64 `gorm:"column:id" json:"id"`
	PostId     int64 `gorm:"column:post_id" json:"post_id"`
	Type       int   `gorm:"column:type" json:"type"`
	UserId     int64 `gorm:"column:user_id" json:"user_id"`
	RootId     int64 `gorm:"column:root_id" json:"root_id"`
	ParentId   int64 `gorm:"column:parent_id" json:"parent_id"`
	LikeCount  int   `gorm:"column:like_count" json:"like_count"`
	ReplyCount int   `gorm:"column:reply_count" json:"reply_count"`
	Score      int   `gorm:"column:score" json:"score"`
	ToUID      int64 `gorm:"column:to_uid" json:"to_uid"`
	DelFlag    int   `gorm:"column:del_flag" json:"del_flag"`
	CreatedAt  int64 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  int64 `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *CommentInfoModel) TableName() string {
	return "comment_info"
}

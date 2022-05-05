package model

type CommentContentModel struct {
	CommentId  int64  `gorm:"column:comment_id" json:"comment_id"`
	Content    string `gorm:"column:content" json:"content"`
	DeviceType string `gorm:"column:device_type" json:"device_type"`
	IP         string `gorm:"column:ip" json:"ip"`
	CreatedAt  int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *CommentContentModel) TableName() string {
	return "comment_content"
}

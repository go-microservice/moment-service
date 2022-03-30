package model

type PostInfoModel struct {
	ID           int     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	PostType     int     `gorm:"column:post_type" json:"post_type"`
	UserID       int64   `gorm:"column:user_id" json:"user_id"`
	Title        string  `gorm:"column:title" json:"title"`
	Content      string  `gorm:"column:content" json:"content"`
	ViewCount    int     `gorm:"column:view_count" json:"view_count"`
	LikeCount    int     `gorm:"column:like_count" json:"like_count"`
	CommentCount int     `gorm:"column:comment_count" json:"comment_count"`
	CollectCount int     `gorm:"column:collect_count" json:"collect_count"`
	ShareCount   int     `gorm:"column:share_count" json:"share_count"`
	Longitude    float64 `gorm:"column:longitude" json:"longitude"`
	Latitude     float64 `gorm:"column:latitude" json:"latitude"`
	Position     string  `gorm:"column:position" json:"position"`
	DelFlag      int     `gorm:"column:del_flag" json:"del_flag"`
	Visible      int     `gorm:"column:visible" json:"visible"`
	CreatedAt    int64   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    int64   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    int64   `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName sets the insert table name for this struct type
func (p *PostInfoModel) TableName() string {
	return "post_infos"
}

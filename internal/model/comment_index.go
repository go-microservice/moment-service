package model

type CommentIndexModel struct {
	Id        int64 `gorm:"column:id" json:"id"`
	ObjType   int   `gorm:"column:obj_type" json:"obj_type"`
	ObjId     int64 `gorm:"column:obj_id" json:"obj_id"`
	RootId    int64 `gorm:"column:root_id" json:"root_id"`
	ParentId  int64 `gorm:"column:parent_id" json:"parent_id"`
	Score     int   `gorm:"column:score" json:"score"`
	UserId    int64 `gorm:"column:user_id" json:"user_id"`
	DelFlag   int   `gorm:"column:del_flag" json:"del_flag"`
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *CommentIndexModel) TableName() string {
	return "comment_index"
}

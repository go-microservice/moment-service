package model

type UserLikeModel struct {
	ID        int64 `gorm:"column:id" json:"id"`
	ObjType   int64 `gorm:"column:obj_type" json:"obj_type"`
	ObjID     int64 `gorm:"column:obj_id" json:"obj_id"`
	UserID    int64 `gorm:"column:user_id" json:"user_id"`
	Status    int   `gorm:"column:status" json:"status"`
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserLikeModel) TableName() string {
	return "user_likes"
}

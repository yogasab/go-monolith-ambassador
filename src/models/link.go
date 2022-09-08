package models

type Link struct {
	ID       int       `json:"id"`
	Code     string    `json:"code"`
	UserID   uint      `json:"user_id"`
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Products []Product `json:"products" gorm:"many2many:link_products"`
	Orders   []Order   `json:"orders,omitempty" gorm:"-"`
}

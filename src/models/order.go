package models

type Order struct {
	ID              int         `json:"id"`
	TransactionID   string      `json:"transaction_id" gorm:"null"`
	UserID          uint        `json:"user_id"`
	Code            string      `json:"code"`
	AmbassadorEmail string      `json:"ambassador_email"`
	FirstName       string      `json:"-"`
	LastName        string      `json:"-"`
	Name            string      `json:"name" gorm:"-"`
	Email           string      `json:"email"`
	Address         string      `json:"address" gorm:"null"`
	Country         string      `json:"country" gorm:"null"`
	City            string      `json:"city" gorm:"null"`
	Zip             string      `json:"zip" gorm:"null"`
	Complete        bool        `json:"-" gorm:"default:false"`
	Total           float64     `json:"total" gorm:"-"`
	OrderItems      []OrderItem `json:"order_items" gorm:"order_id"`
}

type OrderItem struct {
	ID                int     `json:"id"`
	OrderID           uint    `json:"order_id"`
	ProductTitle      string  `json:"product_title"`
	Price             float64 `json:"price"`
	Quantity          uint    `json:"quantity"`
	AdminRevenue      float64 `json:"admin_revenue"`
	AmbassadorRevenue float64 `json:"ambassador_revenue"`
}

func (order *Order) GetFullName() string {
	order.Name = order.FirstName + " " + order.LastName
	return order.Name
}

func (order *Order) GetTotalPrice() float64 {
	var totalPrice float64
	for _, o := range order.OrderItems {
		totalPrice += o.Price * float64(o.Quantity)
	}
	return totalPrice
}

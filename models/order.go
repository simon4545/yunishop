package models

import (
	"time"

	"gorm.io/gorm"
)

// 订单
type Order struct {
	gorm.Model
	CustomerID uint        `json:"customer_id"`
	AgentID    uint        `json:"agent_id"` // 订单所属代理商
	OrderNo    string      `gorm:"size:50;unique" json:"order_no"`
	TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status     OrderStatus `gorm:"size:20" json:"status"`
	Address    string      `gorm:"size:255" json:"address"`
	Phone      string      `gorm:"size:20" json:"phone"`
	Remark     string      `gorm:"size:255" json:"remark"`
}

// 订单商品
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2)" json:"price"` // 下单时的价格
}

// 折扣
type Discount struct {
	gorm.Model
	Name       string    `gorm:"size:100;not null" json:"name"`
	Type       string    `gorm:"size:20" json:"type"` // percentage, fixed_amount
	Value      float64   `gorm:"type:decimal(10,2)" json:"value"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	MinOrder   float64   `gorm:"type:decimal(10,2)" json:"min_order"` // 最低订单金额
	ProductIDs string    `gorm:"type:text" json:"product_ids"`        // 适用的商品ID，多个用逗号分隔
}

// 商品分类
type Category struct {
	gorm.Model
	Name      string `gorm:"size:100;not null" json:"name"`
	ParentID  uint   `json:"parent_id"` // 父分类ID
	SortOrder int    `json:"sort_order"`
}

// 订单状态
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "pending_payment" // 待付款
	OrderStatusPendingReceipt OrderStatus = "pending_receipt" // 待收货
	OrderStatusCompleted      OrderStatus = "completed"       // 已签收
	OrderStatusRefunded       OrderStatus = "refunded"        // 已退款
)

// 售后状态
type AfterSalesStatus string

const (
	AfterSalesStatusPending    AfterSalesStatus = "pending"    // 待处理
	AfterSalesStatusProcessing AfterSalesStatus = "processing" // 处理中
	AfterSalesStatusCompleted  AfterSalesStatus = "completed"  // 已完成
	AfterSalesStatusRejected   AfterSalesStatus = "rejected"   // 已拒绝
)

// 售后类型
type AfterSalesType string

const (
	AfterSalesTypeRefund   AfterSalesType = "refund"   // 退款
	AfterSalesTypeExchange AfterSalesType = "exchange" // 换货
	AfterSalesTypeReturn   AfterSalesType = "return"   // 退货
)

// 售后
type AfterSales struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	OrderID     uint             `json:"order_id"`
	OrderItemID uint             `json:"order_item_id"`
	CustomerID  uint             `json:"customer_id"`
	Type        AfterSalesType   `gorm:"size:20" json:"type"`
	Status      AfterSalesStatus `gorm:"size:20" json:"status"`
	Reason      string           `gorm:"size:255" json:"reason"`
	Description string           `gorm:"type:text" json:"description"`
	Images      string           `gorm:"type:text" json:"images"` // 凭证图片，多个用逗号分隔
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

package model

import (
	"encoding/json"
	"fmt"
	"time"
)



type Order struct {

	OrderID     uint       `json:"order_id"`
    Image       string     `json:"image"`
	CustomerID  uint       `json:"customer_id"`
	LineItems   []LineItem `json:"line_items"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemID   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
	Price    uint `json:"price"`
}

func (o *Order) MarshalLineItems() string {
    js, err := json.Marshal(o.LineItems[0])
    if err != nil {
        panic(err)
    }
    items := string(js)
    for _, item := range o.LineItems[1:] {
        js, err := json.Marshal(item)
        if err != nil {
            panic(err)
        }
        items += fmt.Sprintf(", %v", string(js))
    }
    return fmt.Sprintf(`{"line_items":[%v]}`, items)
}

func (o *Order) UnmarshalLineItems(data string) {
    json.Unmarshal([]byte(data), o)
}



package order

import (
	"database/sql"
	"fmt"
	"math/rand"
	"vs/model"
)

type SqlRepo struct {
	DB *sql.DB
}

func RandomImagePath() string {
	min := 0
	max := 5
	img := rand.Intn(max-min) + min
	return fmt.Sprintf("http://localhost:3000/assets/images/%v.png", img)
}

func OrderIDKey(id uint) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *SqlRepo) Insert(order model.Order) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`order`" + ` (
        image,
		customer_id,
        line_items,
        created_at,
        shipped_at,
        completed_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert order: %w", err)
	}
	_, err = statement.Exec(RandomImagePath(), order.CustomerID, order.MarshalLineItems(), order.CreatedAt, order.ShippedAt, order.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)

	}

	return nil
}

func (r *SqlRepo) FindById(order_id uint) (model.Order, error) {
	model_order := model.Order{}
	line_item := ""
	rows, err := r.DB.Query(`SELECT * FROM `+"`order`"+` WHERE order_id = ?`, order_id)
	if err != nil {
		return model_order, fmt.Errorf("failed to prepare find order: %w", err)
	}

	for rows.Next() {
		rows.Scan(&model_order.OrderID, &model_order.Image, &model_order.CustomerID, &line_item, &model_order.CreatedAt, &model_order.ShippedAt, &model_order.CompletedAt)

	}
	model_order.UnmarshalLineItems(line_item)
	return model_order, nil
}

func (r *SqlRepo) DeleteById(order_id uint) error {
	//_, err := r.DB.Query(`DELETE FROM `+"`order`"+` WHERE order_id = ?`, order_id)
	statement, err := r.DB.Prepare(`DELETE FROM ` + "`order`" + ` WHERE order_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete order: %w", err)
	}

	_, err = statement.Exec(order_id)
	if err != nil {
		return fmt.Errorf("failed to prepare delete EXEC order: %w", err)
	}

	fmt.Println("Сделано")
	return nil

}

func (r *SqlRepo) Update(model_order model.Order) error {

	statement, err := r.DB.Prepare(`UPDATE ` + "`order`" + ` SET
		customer_id = ?,
		image = ?,
        line_items = ?,
        created_at = ?,
        shipped_at = ?,
        completed_at = ?
		WHERE order_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update order: %w", err)
	}
	_, err = statement.Exec(model_order.CustomerID, model_order.Image, model_order.MarshalLineItems(), model_order.CreatedAt, model_order.ShippedAt, model_order.CompletedAt, model_order.OrderID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (r *SqlRepo) 	FindAll() ([]model.Order, error) {

	array := []model.Order{}

	rows, err := r.DB.Query(`SELECT * FROM ` + "`order`" + ``)
	if err != nil {
		return []model.Order{}, fmt.Errorf("failed to prepare FindAll order: %w", err)
	}

	for rows.Next() {
		model_order := model.Order{}
		line_item := ""
		rows.Scan(&model_order.OrderID, &model_order.Image, &model_order.CustomerID, &line_item, &model_order.CreatedAt, &model_order.ShippedAt, &model_order.CompletedAt)
		model_order.UnmarshalLineItems(line_item)
		array = append(array, model_order)
	}
	return array, nil
}

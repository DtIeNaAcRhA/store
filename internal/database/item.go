package database

import (
	"fmt"
	"store/internal/model"
)

func GetItemByID(id int) (*model.Item, error) {
	row := DB.QueryRow(`
		SELECT id, id_user, title, description_item, image_path, price, created_at
		FROM item
		WHERE id = ?`, id)

	var item model.Item
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.DescriptionItem,
		&item.ImagePath,
		&item.Price,
		&item.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateItem(item *model.Item) error {
	stmt := `INSERT INTO item (id_user, title, description_item, image_path, price) VALUES (?, ?, ?, ?, ?)`
	res, err := DB.Exec(stmt, item.UserID, item.Title, item.DescriptionItem, item.ImagePath, item.Price)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	item.ID = int(id)
	return nil
}

// Получить товары, отсортированные по дате (новые первыми)
func GetItemsByDate(minPrice, maxPrice float64, order string, limit, offset int) ([]model.Item, error) {
	return getItemsWithFilters("item.created_at", minPrice, maxPrice, order, limit, offset)
}

// Получить товары, отсортированные по цене
func GetItemsByPrice(minPrice, maxPrice float64, order string, limit, offset int) ([]model.Item, error) {
	return getItemsWithFilters("item.price", minPrice, maxPrice, order, limit, offset)
}

// Внутренняя общая функция для выборки товаров
func getItemsWithFilters(sortField string, minPrice, maxPrice float64, order string, limit, offset int) ([]model.Item, error) {
	items := []model.Item{}

	query := `
	SELECT 
		item.id, 
		item.id_user,
		item.title, 
		item.description_item, 
		item.image_path, 
		item.price, 
		item.created_at, 
		user.username AS author_login
	FROM item
	JOIN user ON user.id = item.id_user
	WHERE 1=1
	`

	if minPrice > 0 {
		query += fmt.Sprintf(" AND item.price >= %f", minPrice)
	}
	if maxPrice > 0 {
		query += fmt.Sprintf(" AND item.price <= %f", maxPrice)
	}

	// Проверка корректности сортировки
	// if order != "ASC" {
	// 	order = "DESC"
	// }

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", sortField, order, limit, offset)

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Title,
			&item.DescriptionItem,
			&item.ImagePath,
			&item.Price,
			&item.CreatedAt,
			&item.AuthorLogin,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

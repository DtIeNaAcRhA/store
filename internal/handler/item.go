package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"store/internal/database"
	"store/internal/model"
	"strconv"
)

// создание товара
func CreateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	item := &model.Item{
		UserID:          userID,
		Title:           req.Title,
		DescriptionItem: req.Description,
		ImagePath:       req.ImagePath,
		Price:           req.Price,
	}

	err = database.CreateItem(item)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	item, err = database.GetItemByID(item.ID)
	if err != nil {
		fmt.Print(err)
	}

	username := ""
	user, err := database.GetUserByID(userID)
	if err == nil {
		username = user.Username

	}

	response := model.ItemResponse{
		ID:              item.ID,
		Title:           item.Title,
		DescriptionItem: item.DescriptionItem,
		ImagePath:       item.ImagePath,
		Price:           item.Price,
		CreatedAt:       item.CreatedAt.String(),
		AuthorLogin:     username,
		IsOwner:         userID == item.UserID,
	}

	JSON(w, response, http.StatusCreated)
}

func ParseListItemsRequest(r *http.Request) model.ListItemsRequest {
	q := r.URL.Query()

	// Значения по умолчанию
	req := model.ListItemsRequest{
		Sort:     "date",
		Order:    "desc",
		MinPrice: -1,
		MaxPrice: -1,
		Limit:    20,
		Offset:   0,
	}

	if v := q.Get("sort"); v == "price" {
		req.Sort = v
		log.Print("получили sort")
	}
	if v := q.Get("order"); v == "asc" {
		req.Order = v
		log.Print("получили order")
	}

	if v := q.Get("min"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.MinPrice = f
			log.Print("получили min")
		}
	}
	if v := q.Get("max"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.MaxPrice = f
			log.Print("получили max")
		}
	}
	if v := q.Get("limit"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			req.Limit = i
			log.Print("получили limit")

		}
	}
	if v := q.Get("offset"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			req.Offset = i
			log.Print("получили offset")
		}
	}

	return req
}

func ListItems(w http.ResponseWriter, r *http.Request) {
	req := ParseListItemsRequest(r)

	var items []model.Item
	var err error

	// Пытаемся получить userID (если пользователь авторизован)
	userID, err := getUserIDFromToken(r)
	if err != nil {
		userID = 0
	}

	switch req.Sort {
	case "price":
		items, err = database.GetItemsByPrice(req.MinPrice, req.MaxPrice, req.Order, req.Limit, req.Offset)
	default:
		items, err = database.GetItemsByDate(req.MinPrice, req.MaxPrice, req.Order, req.Limit, req.Offset)
	}

	if err != nil {
		http.Error(w, "Ошибка при получении товаров", http.StatusInternalServerError)
		return
	}

	resp := model.ListItemsResponse{
		Items: convertItemsToResponse(items, userID),
		Total: len(items),
	}

	JSON(w, resp, http.StatusOK)
}

func convertItemsToResponse(items []model.Item, currentUserID int) []model.ItemResponse {
	var result []model.ItemResponse

	for _, item := range items {
		user, err := database.GetUserByID(item.UserID)
		if err != nil {
			log.Println("не получилось достать юзера")
		}

		resp := model.ItemResponse{
			ID:              item.ID,
			Title:           item.Title,
			DescriptionItem: item.DescriptionItem,
			ImagePath:       item.ImagePath,
			Price:           item.Price,
			CreatedAt:       item.CreatedAt.String(),
			AuthorLogin:     user.Username,
			IsOwner:         currentUserID == item.UserID,
		}
		log.Println("добавили", resp)
		result = append(result, resp)
	}

	return result
}

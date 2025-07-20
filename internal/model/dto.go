package model

// Регистрация пользователя
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

// Авторизация
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Объявление
type ItemResponse struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	DescriptionItem string  `json:"description_item"`
	ImagePath       string  `json:"imagePath"`
	Price           float64 `json:"price"`
	CreatedAt       string  `json:"createdAt"`
	AuthorLogin     string  `json:"author"`
	IsOwner         bool    `json:"isOwner"`
}

type CreateItemRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImagePath   string  `json:"imagePath"`
	Price       float64 `json:"price"`
}

type ListItemsRequest struct {
	Sort     string  `json:"sort"`   // "price" или "date"
	Order    string  `json:"order"`  // "asc" или "desc"
	MinPrice float64 `json:"min"`    // фильтр: от
	MaxPrice float64 `json:"max"`    // фильтр: до
	Limit    int     `json:"limit"`  // количество на странице
	Offset   int     `json:"offset"` // смещение
}

// Список товаров
type ListItemsResponse struct {
	Items []ItemResponse `json:"items"`
	Total int            `json:"total"`
}

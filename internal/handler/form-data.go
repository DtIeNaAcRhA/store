package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	// Ограничим размер файла
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Сохраняем файл
	dstPath := fmt.Sprintf("uploads/%s", handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Ошибка копирования", http.StatusInternalServerError)
		return
	}

	// Возвращаем путь
	JSON(w, map[string]string{"imagePath": dstPath}, http.StatusOK)
}

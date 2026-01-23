package static

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServeHtml(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("internal", "web", "templates", "html", filename)

		// Проверяем существование файла
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		// Отдаем файл
		http.ServeFile(w, r, filePath)
	}
}

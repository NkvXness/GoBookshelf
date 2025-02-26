package api

import (
	"net/http"
	"strings"
)

// Router представляет собой простой маршрутизатор для API
type Router struct {
	routes      map[string]map[string]http.HandlerFunc
	middlewares []Middleware
}

// Middleware представляет собой функцию промежуточного ПО
type Middleware func(http.HandlerFunc) http.HandlerFunc

// NewRouter создает новый экземпляр маршрутизатора
func NewRouter() *Router {
	return &Router{
		routes:      make(map[string]map[string]http.HandlerFunc),
		middlewares: []Middleware{},
	}
}

// Use добавляет middleware в цепочку обработки
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// HandleFunc регистрирует обработчик для указанного пути и метода
func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	if _, exists := r.routes[path]; !exists {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}

// GET регистрирует обработчик для GET запросов
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodGet, path, handler)
}

// POST регистрирует обработчик для POST запросов
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodPost, path, handler)
}

// PUT регистрирует обработчик для PUT запросов
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodPut, path, handler)
}

// DELETE регистрирует обработчик для DELETE запросов
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.HandleFunc(http.MethodDelete, path, handler)
}

// ServeHTTP реализует интерфейс http.Handler и обрабатывает все запросы
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Извлекаем базовый путь без query параметров
	path := req.URL.Path

	// Проверяем точное совпадение пути
	handlers, exists := r.routes[path]
	if !exists {
		// Если точного совпадения нет, пробуем найти обработчик
		// который соответствует более общему пути (без ID в конце)
		found := false
		for routePath, routeHandlers := range r.routes {
			// Проверяем, оканчивается ли путь на {id} (шаблон)
			if strings.HasSuffix(routePath, "/{id}") {
				// Удаляем /{id} из конца
				basePath := routePath[:len(routePath)-4]
				// Проверяем, начинается ли запрашиваемый путь с этого базового пути
				if strings.HasPrefix(path, basePath) && len(path) > len(basePath) {
					handlers = routeHandlers
					found = true
					break
				}
			}
		}

		if !found {
			http.NotFound(w, req)
			return
		}
	}

	handler, exists := handlers[req.Method]
	if !exists {
		w.Header().Set("Allow", getAllowedMethods(handlers))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Применяем middleware в обратном порядке
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	handler(w, req)
}

// getAllowedMethods возвращает строку с разрешенными методами
func getAllowedMethods(handlers map[string]http.HandlerFunc) string {
	methods := ""
	for method := range handlers {
		if methods != "" {
			methods += ", "
		}
		methods += method
	}
	return methods
}

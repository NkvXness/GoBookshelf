package api

import (
	"net/http"
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
	path := req.URL.Path
	handlers, exists := r.routes[path]
	if !exists {
		http.NotFound(w, req)
		return
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

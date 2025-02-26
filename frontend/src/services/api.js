import axios from "axios";

export const api = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// Перехватчики
api.interceptors.request.use(
  (config) => {
    console.debug("🚀 [API]", config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Обработка ошибок
    const message = error.response?.data?.message || "Произошла ошибка";
    console.error("🚨 [API]", message);
    return Promise.reject(error);
  },
);

// API endpoints
export const booksApi = {
  getBooks: (page = 1, pageSize = 10) =>
    api.get(`/books?page=${page}&page_size=${pageSize}`),
  createBook: (data) => api.post("/books", data),
  updateBook: (id, data) => api.put(`/books?id=${id}`, data),
  deleteBook: (id) => api.delete(`/books?id=${id}`),
};

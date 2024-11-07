import axios from "axios";

// Создаем экземпляр axios с базовой конфигурацией
const api = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// Добавляем перехватчик для запросов
api.interceptors.request.use(
  (config) => {
    console.log("Making request to:", config.url, config.method?.toUpperCase());
    return config;
  },
  (error) => {
    console.error("Request error:", error);
    return Promise.reject(error);
  }
);

// Добавляем перехватчик для ответов
api.interceptors.response.use(
  (response) => {
    console.log(
      "Received response from:",
      response.config.url,
      response.status
    );
    return response;
  },
  (error) => {
    console.error("Response error:", error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default api;

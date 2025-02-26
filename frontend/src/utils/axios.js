import axios from "axios";

// Создаем экземпляр axios с базовой конфигурацией
const api = axios.create({
  baseURL: "", // Пустая строка, запросы будут идти относительно текущего домена
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000,
});

// Добавляем перехватчик для запросов
api.interceptors.request.use(
  (config) => {
    // Логирование запросов
    console.log(`${config.method?.toUpperCase()} ${config.url}`);
    
    // Не добавляем префикс /api, так как в консоли видно, что запросы уже идут с /api
    
    return config;
  },
  (error) => {
    console.error("Ошибка запроса:", error);
    return Promise.reject(error);
  }
);

// Добавляем перехватчик для ответов
api.interceptors.response.use(
  (response) => {
    console.log(`Ответ от ${response.config.url}: ${response.status}`);
    return response;
  },
  (error) => {
    console.error("Ошибка ответа:", error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default api;
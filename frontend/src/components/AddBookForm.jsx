import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { X, Book, User, Hash, Calendar } from "lucide-react";
import api from "../utils/axios";

// Утилита для валидации ISBN
const validateISBN = (isbn) => {
  // Удаляем все не цифровые символы
  const cleanISBN = isbn.replace(/[^0-9]/g, "");

  if (cleanISBN.length !== 13) {
    return "ISBN должен содержать 13 цифр";
  }

  // Проверка контрольной суммы ISBN-13
  let sum = 0;
  for (let i = 0; i < 12; i++) {
    sum += parseInt(cleanISBN[i]) * (i % 2 === 0 ? 1 : 3);
  }
  
  const checkDigit = (10 - (sum % 10)) % 10;
  
  if (parseInt(cleanISBN[12]) !== checkDigit) {
    return "Неверная контрольная сумма ISBN";
  }

  return ""; // Пустая строка означает отсутствие ошибок
};

// Форматирование ISBN в читаемый вид
const formatISBN = (isbn) => {
  // Удаляем все не цифровые символы
  const cleanISBN = isbn.replace(/[^0-9]/g, "");
  // Форматируем ISBN в формат XXX-X-XXX-XXXXX-X
  if (cleanISBN.length === 13) {
    return `${cleanISBN.slice(0, 3)}-${cleanISBN.slice(
      3,
      4
    )}-${cleanISBN.slice(4, 7)}-${cleanISBN.slice(7, 12)}-${cleanISBN.slice(
      12
    )}`;
  }
  return isbn;
};

const AddBookForm = ({ onClose }) => {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({
    title: "",
    author: "",
    isbn: "",
    published: new Date().toISOString().split("T")[0],
  });
  const [errors, setErrors] = useState({});

  // Обработка изменения ISBN с валидацией
  const handleISBNChange = (e) => {
    const value = e.target.value;
    // Разрешаем вводить только цифры и дефисы
    const isbn = value.replace(/[^0-9-]/g, "");

    setFormData({ ...formData, isbn });

    // Валидация при вводе
    if (isbn.length > 0) {
      const error = validateISBN(isbn);
      setErrors({ ...errors, isbn: error });
    } else {
      setErrors({ ...errors, isbn: "" });
    }
  };

  // Мутация для добавления книги
  const addBookMutation = useMutation({
    mutationFn: (newBook) => {
      console.log("Отправка данных:", newBook);
      return api.post("/api/books", newBook);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["books"]);
      onClose();
    },
    onError: (error) => {
      console.error("Ошибка добавления книги:", error);
      alert("Ошибка добавления книги: " + (error.response?.data?.message || "Проверьте формат ISBN"));
    },
  });

  // Обработчик отправки формы
  const handleSubmit = (e) => {
    e.preventDefault();
    
    // Валидация формы
    const newErrors = {};
    if (!formData.title.trim()) {
      newErrors.title = "Название книги обязательно";
    }
    if (!formData.author.trim()) {
      newErrors.author = "Автор обязателен";
    }
    
    // Проверяем ISBN
    const isbnError = validateISBN(formData.isbn);
    if (isbnError) {
      newErrors.isbn = isbnError;
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    // Форматируем ISBN перед отправкой
    const formattedData = {
      ...formData,
      isbn: formatISBN(formData.isbn),
      published: new Date(formData.published).toISOString(),
    };

    addBookMutation.mutate(formattedData);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full">
        {/* Заголовок */}
        <div className="flex justify-between items-center p-5 border-b dark:border-gray-700">
          <div className="flex items-center">
            <Book className="w-5 h-5 text-blue-600 dark:text-blue-400 mr-2" />
            <h3 className="text-lg font-bold">Добавить новую книгу</h3>
          </div>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Форма */}
        <form onSubmit={handleSubmit} className="p-5">
          <div className="space-y-4">
            <div>
              <label className="block mb-1 flex items-center text-gray-700 dark:text-gray-300">
                <Book className="w-4 h-4 mr-1" />
                Название
              </label>
              <input
                type="text"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                className={`w-full p-2 border rounded ${errors.title ? "border-red-500" : ""}`}
                required
              />
              {errors.title && <p className="text-red-500 text-sm mt-1">{errors.title}</p>}
            </div>

            <div>
              <label className="block mb-1 flex items-center text-gray-700 dark:text-gray-300">
                <User className="w-4 h-4 mr-1" />
                Автор
              </label>
              <input
                type="text"
                value={formData.author}
                onChange={(e) => setFormData({ ...formData, author: e.target.value })}
                className={`w-full p-2 border rounded ${errors.author ? "border-red-500" : ""}`}
                required
              />
              {errors.author && <p className="text-red-500 text-sm mt-1">{errors.author}</p>}
            </div>

            <div>
              <label className="block mb-1 flex items-center text-gray-700 dark:text-gray-300">
                <Hash className="w-4 h-4 mr-1" />
                ISBN
              </label>
              <input
                type="text"
                value={formatISBN(formData.isbn)}
                onChange={handleISBNChange}
                className={`w-full p-2 border rounded ${errors.isbn ? "border-red-500" : ""}`}
                placeholder="978-3-16-148410-0"
                required
              />
              {errors.isbn ? (
                <p className="text-red-500 text-sm mt-1">{errors.isbn}</p>
              ) : (
                <p className="text-xs text-gray-500 mt-1">
                  Формат: 978-3-16-148410-0 (должен быть валидный ISBN-13)
                </p>
              )}
            </div>

            <div>
              <label className="block mb-1 flex items-center text-gray-700 dark:text-gray-300">
                <Calendar className="w-4 h-4 mr-1" />
                Дата публикации
              </label>
              <input
                type="date"
                value={formData.published}
                onChange={(e) => setFormData({ ...formData, published: e.target.value })}
                className="w-full p-2 border rounded"
                required
              />
            </div>
          </div>

          {/* Кнопки */}
          <div className="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded"
            >
              Отмена
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded"
              disabled={addBookMutation.isLoading || Object.keys(errors).some(key => errors[key])}
            >
              {addBookMutation.isLoading ? "Добавление..." : "Добавить книгу"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddBookForm;
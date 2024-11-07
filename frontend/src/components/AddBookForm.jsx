import React, { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { X, AlertCircle } from "lucide-react";
import api from "../utils/axios";

const AddBookForm = ({ onClose }) => {
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({
    title: "",
    author: "",
    isbn: "",
    published: new Date().toISOString().split("T")[0],
  });
  const [errors, setErrors] = useState({});

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

    return "";
  };

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

  const addBookMutation = useMutation({
    mutationFn: (newBook) => api.post("/books", newBook),
    onError: (error) => {
      console.error("Error adding book:", error);
      alert("Failed to add book. Please check the console for details.");
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["books"]);
      onClose();
    },
  });

  const handleSubmit = (e) => {
    e.preventDefault();

    // Проверяем ISBN перед отправкой
    const isbnError = validateISBN(formData.isbn);
    if (isbnError) {
      setErrors({ ...errors, isbn: isbnError });
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
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div className="bg-white rounded-lg p-8 max-w-md w-full">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">Add New Book</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Title
            </label>
            <input
              type="text"
              required
              className="mt-1 block w-full border rounded-md shadow-sm p-2"
              value={formData.title}
              onChange={(e) =>
                setFormData({ ...formData, title: e.target.value })
              }
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Author
            </label>
            <input
              type="text"
              required
              className="mt-1 block w-full border rounded-md shadow-sm p-2"
              value={formData.author}
              onChange={(e) =>
                setFormData({ ...formData, author: e.target.value })
              }
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              ISBN-13
              <span className="text-gray-500 text-xs ml-1">
                (13 цифр, например: 978-3-16-148410-0)
              </span>
            </label>
            <div className="relative">
              <input
                type="text"
                required
                placeholder="978-X-XXX-XXXXX-X"
                className={`mt-1 block w-full border rounded-md shadow-sm p-2 ${
                  errors.isbn ? "border-red-500" : "border-gray-300"
                }`}
                value={formatISBN(formData.isbn)}
                onChange={handleISBNChange}
              />
              {errors.isbn && (
                <div className="absolute right-2 top-1/2 transform -translate-y-1/2">
                  <AlertCircle className="h-5 w-5 text-red-500" />
                </div>
              )}
            </div>
            {errors.isbn && (
              <p className="mt-1 text-sm text-red-600">{errors.isbn}</p>
            )}
            <p className="mt-1 text-xs text-gray-500">
              Введите 13 цифр ISBN. Дефисы будут добавлены автоматически.
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Published Date
            </label>
            <input
              type="date"
              required
              className="mt-1 block w-full border rounded-md shadow-sm p-2"
              value={formData.published}
              onChange={(e) =>
                setFormData({ ...formData, published: e.target.value })
              }
            />
          </div>

          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border rounded-md text-gray-700 hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={
                addBookMutation.isLoading ||
                Object.keys(errors).some((key) => errors[key])
              }
            >
              {addBookMutation.isLoading ? "Adding..." : "Add Book"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddBookForm;

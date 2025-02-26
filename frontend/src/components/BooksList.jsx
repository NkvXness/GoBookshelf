import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Edit, Trash2, User, Hash, Calendar } from "lucide-react";
import api from "../utils/axios";

// Компонент карточки книги
const BookCard = ({ book, onEdit, onDelete }) => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden border border-gray-200 dark:border-gray-700">
      <div className="p-5">
        <h3 className="text-lg font-bold mb-3 line-clamp-2">{book.title}</h3>
        
        <div className="space-y-2 text-sm">
          <div className="flex items-center">
            <User size={16} className="mr-2 text-gray-500 dark:text-gray-400" />
            <span>{book.author}</span>
          </div>
          
          <div className="flex items-center">
            <Hash size={16} className="mr-2 text-gray-500 dark:text-gray-400" />
            <span className="font-mono">{book.isbn}</span>
          </div>
          
          <div className="flex items-center">
            <Calendar size={16} className="mr-2 text-gray-500 dark:text-gray-400" />
            <span>{new Date(book.published).toLocaleDateString()}</span>
          </div>
        </div>
      </div>
      
      <div className="flex border-t border-gray-200 dark:border-gray-700">
        <button
          onClick={() => onEdit(book)}
          className="flex-1 py-2 text-blue-600 dark:text-blue-400 hover:bg-gray-50 dark:hover:bg-gray-700 flex items-center justify-center"
        >
          <Edit size={18} className="mr-1" />
          <span>Изменить</span>
        </button>
        
        <div className="w-px bg-gray-200 dark:bg-gray-700"></div>
        
        <button
          onClick={() => onDelete(book)}
          className="flex-1 py-2 text-red-600 dark:text-red-400 hover:bg-gray-50 dark:hover:bg-gray-700 flex items-center justify-center"
        >
          <Trash2 size={18} className="mr-1" />
          <span>Удалить</span>
        </button>
      </div>
    </div>
  );
};

// Основной компонент списка книг
const BooksList = ({ searchQuery }) => {
  const [page, setPage] = useState(1);
  const [editingBook, setEditingBook] = useState(null);
  const [bookToDelete, setBookToDelete] = useState(null);
  const queryClient = useQueryClient();

  // Запрос на получение книг
  const { data, isLoading, isError } = useQuery({
    queryKey: ["books", page, searchQuery],
    queryFn: async () => {
      // Для поиска используем формат, который точно работает в бэкенде:
      // Просто получаем все книги и фильтруем на клиенте
      const response = await api.get(`/api/books?page=${page}`);
      return response.data;
    },
  });

  // Мутация для удаления книги
const deleteMutation = useMutation({
  mutationFn: (id) => {
    // Используем POST с action=delete в параметрах вместо метода DELETE
    return api.post(`/api/books?id=${id}&action=delete`);
  },
  onSuccess: () => {
    queryClient.invalidateQueries(["books"]);
    setBookToDelete(null);
    alert("Книга успешно удалена");
  },
  onError: (error) => {
    console.error("Delete error:", error);
    alert("Ошибка при удалении книги: " + (error.response?.data?.message || "Неизвестная ошибка"));
  }
});

  // Мутация для обновления книги
  // ВАЖНО: Используем POST вместо PUT, так как на сервере ожидается именно POST
  const updateMutation = useMutation({
    mutationFn: (book) => {
      // Используем тот же формат URL, что и при удалении
      return api.post(`/api/books?id=${book.id}`, {
        title: book.title,
        author: book.author,
        isbn: book.isbn,
        published: book.published instanceof Date
          ? book.published.toISOString()
          : new Date(book.published).toISOString()
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["books"]);
      setEditingBook(null);
      alert("Книга успешно обновлена");
    },
    onError: (error) => {
      console.error("Update error:", error);
      alert("Ошибка при обновлении книги: " + (error.response?.data?.message || "Неизвестная ошибка"));
    }
  });

  // Обработчик начала редактирования книги
  const handleEdit = (book) => {
    setEditingBook({
      ...book,
      published: new Date(book.published).toISOString().split("T")[0],
    });
  };

  // Обработчик удаления книги
  const handleDelete = (book) => {
    if (confirm(`Вы уверены, что хотите удалить книгу "${book.title}"?`)) {
      deleteMutation.mutate(book.id);
    }
  };

  // Обработчик отправки формы редактирования
  const handleUpdateSubmit = (e) => {
    e.preventDefault();
    updateMutation.mutate(editingBook);
  };

  // Фильтрация книг на стороне клиента
  const filterBooks = (books, query) => {
    if (!query || query.trim() === "") return books;
    
    const lowercaseQuery = query.toLowerCase().trim();
    return books.filter(book => 
      book.title.toLowerCase().includes(lowercaseQuery) ||
      book.author.toLowerCase().includes(lowercaseQuery) ||
      book.isbn.includes(lowercaseQuery)
    );
  };

  // Отображение состояния загрузки
  if (isLoading) {
    return <div className="text-center py-10">Загрузка...</div>;
  }

  // Отображение ошибки
  if (isError) {
    return <div className="text-center py-10 text-red-500">Ошибка при загрузке данных</div>;
  }

  // Фильтруем книги на стороне клиента
  const filteredBooks = searchQuery ? filterBooks(data.books, searchQuery) : data.books;

  // Если нет книг
  if (!filteredBooks || filteredBooks.length === 0) {
    return (
      <div className="text-center py-10">
        <p className="text-gray-600 dark:text-gray-400">
          {searchQuery ? "По вашему запросу ничего не найдено" : "Нет доступных книг"}
        </p>
      </div>
    );
  }

  // Модальное окно редактирования книги
  if (editingBook) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full">
          <h3 className="text-xl font-bold mb-4">Редактирование книги</h3>
          
          <form onSubmit={handleUpdateSubmit}>
            <div className="space-y-4">
              <div>
                <label className="block mb-1">Название</label>
                <input
                  type="text"
                  className="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
                  value={editingBook.title}
                  onChange={(e) => setEditingBook({ ...editingBook, title: e.target.value })}
                  required
                />
              </div>
              
              <div>
                <label className="block mb-1">Автор</label>
                <input
                  type="text"
                  className="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
                  value={editingBook.author}
                  onChange={(e) => setEditingBook({ ...editingBook, author: e.target.value })}
                  required
                />
              </div>
              
              <div>
                <label className="block mb-1">ISBN</label>
                <input
                  type="text"
                  className="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
                  value={editingBook.isbn}
                  onChange={(e) => setEditingBook({ ...editingBook, isbn: e.target.value })}
                  required
                />
                <p className="text-xs text-gray-500 mt-1">
                  Формат: 978-3-16-148410-0 (ISBN-13)
                </p>
              </div>
              
              <div>
                <label className="block mb-1">Дата публикации</label>
                <input
                  type="date"
                  className="w-full p-2 border rounded dark:bg-gray-700 dark:border-gray-600"
                  value={editingBook.published}
                  onChange={(e) => setEditingBook({ ...editingBook, published: e.target.value })}
                  required
                />
              </div>
            </div>
            
            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                className="px-4 py-2 border rounded"
                onClick={() => setEditingBook(null)}
              >
                Отмена
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-blue-600 text-white rounded"
                disabled={updateMutation.isLoading}
              >
                {updateMutation.isLoading ? "Сохранение..." : "Сохранить"}
              </button>
            </div>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div>
      {/* Сетка книг */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredBooks.map((book) => (
          <BookCard
            key={book.id}
            book={book}
            onEdit={handleEdit}
            onDelete={handleDelete}
          />
        ))}
      </div>
      
      {/* Пагинация */}
      {!searchQuery && (
        <div className="mt-6 flex justify-between items-center">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-4 py-2 border rounded disabled:opacity-50"
          >
            Назад
          </button>
          
          <span>
            Страница {page} из {Math.ceil(data.total_books / data.page_size) || 1}
          </span>
          
          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={!data.total_books || page >= Math.ceil(data.total_books / data.page_size)}
            className="px-4 py-2 border rounded disabled:opacity-50"
          >
            Вперед
          </button>
        </div>
      )}
    </div>
  );
};

export default BooksList;
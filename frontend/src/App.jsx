import { useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useTheme } from "./hooks/useTheme";
import { Sun, Moon, Book, Plus } from "lucide-react";
import BooksList from "./components/BooksList";
import AddBookForm from "./components/AddBookForm";
import { ToastProvider } from "./contexts/ToastContext";

// Создаем экземпляр QueryClient
const queryClient = new QueryClient();

function App() {
  const [showAddForm, setShowAddForm] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const { isDarkMode, toggleTheme } = useTheme();

  return (
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <div className={isDarkMode ? "dark" : ""}>
          <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-white">
            {/* Хедер */}
            <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
              <div className="container mx-auto px-4 py-4 flex justify-between items-center">
                <div className="flex items-center space-x-2">
                  <Book className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                  <h1 className="text-2xl font-bold">GoBookshelf</h1>
                </div>
                
                <button 
                  onClick={toggleTheme}
                  className="p-2 rounded-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                  aria-label="Переключить тему"
                >
                  {isDarkMode ? (
                    <Sun className="h-5 w-5 text-yellow-500" />
                  ) : (
                    <Moon className="h-5 w-5 text-gray-700" />
                  )}
                </button>
              </div>
            </header>

            {/* Основной контент */}
            <main className="container mx-auto px-4 py-6">
              <h2 className="text-xl font-bold mb-4">Список книг</h2>
              
              {/* Поиск и добавление */}
              <div className="flex flex-col sm:flex-row justify-between mb-6 gap-4">
                <div className="relative w-full sm:w-64">
                  <input
                    type="text"
                    placeholder="Поиск книг..."
                    className="w-full px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
                
                <button
                  onClick={() => setShowAddForm(true)}
                  className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg"
                >
                  <Plus size={20} />
                  <span>Добавить книгу</span>
                </button>
              </div>
              
              {/* Список книг */}
              <BooksList searchQuery={searchQuery} />
              
              {/* Модальное окно добавления книги */}
              {showAddForm && <AddBookForm onClose={() => setShowAddForm(false)} />}
            </main>
          </div>
        </div>
      </ToastProvider>
    </QueryClientProvider>
  );
}

export default App;
import { useState } from "react";
import { Search } from "lucide-react";

const SearchBar = ({ onSearch, darkMode }) => {
  const [query, setQuery] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    onSearch(query.trim());
  };

  return (
    <form onSubmit={handleSubmit} className="relative w-full sm:max-w-md">
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search className={`h-5 w-5 ${darkMode ? 'text-gray-400' : 'text-gray-500'}`} />
        </div>
        
        <input
          type="text"
          placeholder="Поиск книг..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="block w-full pl-10 pr-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg 
                   shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500
                   bg-white dark:bg-gray-700 text-gray-900 dark:text-white
                   transition-colors duration-200"
        />
        
        <button
          type="submit"
          className="absolute inset-y-0 right-0 flex items-center px-4 text-gray-700 dark:text-gray-300 
                     hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
        >
          Найти
        </button>
      </div>
    </form>
  );
};

export default SearchBar;
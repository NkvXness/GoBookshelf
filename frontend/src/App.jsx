import { useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import BooksList from "./components/BooksList";
import AddBookForm from "./components/AddBookForm";

const queryClient = new QueryClient();

function App() {
  const [showAddForm, setShowAddForm] = useState(false);

  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-gray-100">
        <nav className="bg-white shadow-sm">
          <div className="max-w-7xl mx-auto px-4 py-3">
            <h1 className="text-2xl font-bold text-gray-800">GoBookshelf</h1>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold text-gray-700">Books List</h2>
            <button
              onClick={() => setShowAddForm(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Add New Book
            </button>
          </div>

          <BooksList />
          {showAddForm && <AddBookForm onClose={() => setShowAddForm(false)} />}
        </main>
      </div>
    </QueryClientProvider>
  );
}

export default App;

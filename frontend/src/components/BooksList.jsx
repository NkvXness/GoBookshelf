import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Trash2, Pencil } from "lucide-react";
import api from "../utils/axios";

const BooksList = () => {
  const [page, setPage] = useState(1);
  const [editingBook, setEditingBook] = useState(null);
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ["books", page],
    queryFn: () =>
      api.get(`/books?page=${page}&page_size=10`).then((res) => res.data),
  });

  const deleteMutation = useMutation({
    mutationFn: (id) => api.delete(`/books?id=${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries(["books"]);
    },
    onError: (error) => {
      console.error("Delete error:", error.response?.data || error.message);
      alert("Failed to delete book. Please try again.");
    },
  });

  const updateMutation = useMutation({
    mutationFn: (book) => {
      const updatedBook = {
        id: book.id,
        title: book.title,
        author: book.author,
        isbn: book.isbn,
        published:
          book.published instanceof Date
            ? book.published.toISOString()
            : new Date(book.published).toISOString(),
      };

      // Добавляем отладочный вывод
      console.log("Sending update request with data:", updatedBook);

      return api.put(`/books?id=${book.id}`, updatedBook);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["books"]);
      setEditingBook(null);
    },
    onError: (error) => {
      console.error("Update error details:", {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
      });
      alert("Failed to update book. Please check the console for details.");
    },
  });

  const handleUpdateClick = (book) => {
    // Добавляем отладочный вывод
    console.log("Original book data:", book);

    setEditingBook({
      id: book.id,
      title: book.title,
      author: book.author,
      isbn: book.isbn,
      published: new Date(book.published).toISOString().split("T")[0],
    });
  };

  if (isLoading) {
    return <div className="flex justify-center p-8">Loading...</div>;
  }

  if (!data?.books) {
    return <div className="text-center p-8">No books found</div>;
  }

  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Title
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Author
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              ISBN
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Published
            </th>
            <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {data.books.map((book) => (
            <tr key={book.id}>
              <td className="px-6 py-4 whitespace-nowrap">
                {editingBook?.id === book.id ? (
                  <input
                    type="text"
                    className="border rounded p-1 w-full"
                    value={editingBook.title}
                    onChange={(e) =>
                      setEditingBook({ ...editingBook, title: e.target.value })
                    }
                  />
                ) : (
                  book.title
                )}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                {editingBook?.id === book.id ? (
                  <input
                    type="text"
                    className="border rounded p-1 w-full"
                    value={editingBook.author}
                    onChange={(e) =>
                      setEditingBook({ ...editingBook, author: e.target.value })
                    }
                  />
                ) : (
                  book.author
                )}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                {editingBook?.id === book.id ? (
                  <input
                    type="text"
                    className="border rounded p-1 w-full"
                    value={editingBook.isbn}
                    onChange={(e) =>
                      setEditingBook({ ...editingBook, isbn: e.target.value })
                    }
                  />
                ) : (
                  book.isbn
                )}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                {editingBook?.id === book.id ? (
                  <input
                    type="date"
                    className="border rounded p-1 w-full"
                    value={editingBook.published.split("T")[0]}
                    onChange={(e) =>
                      setEditingBook({
                        ...editingBook,
                        published: new Date(e.target.value).toISOString(),
                      })
                    }
                  />
                ) : (
                  new Date(book.published).toLocaleDateString()
                )}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-center">
                {editingBook?.id === book.id ? (
                  <div className="flex justify-center space-x-2">
                    <button
                      className="text-green-600 hover:text-green-900"
                      onClick={() => updateMutation.mutate(editingBook)}
                    >
                      Save
                    </button>
                    <button
                      className="text-gray-600 hover:text-gray-900"
                      onClick={() => setEditingBook(null)}
                    >
                      Cancel
                    </button>
                  </div>
                ) : (
                  <div className="flex justify-center space-x-2">
                    <button
                      className="text-blue-600 hover:text-blue-900"
                      onClick={() => handleUpdateClick(book)}
                    >
                      <Pencil className="h-5 w-5" />
                    </button>
                    <button
                      className="text-red-600 hover:text-red-900"
                      onClick={() => deleteMutation.mutate(book.id)}
                    >
                      <Trash2 className="h-5 w-5" />
                    </button>
                  </div>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="px-6 py-4 flex justify-between items-center bg-gray-50">
        <button
          className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          disabled={page === 1}
          onClick={() => setPage((p) => p - 1)}
        >
          Previous
        </button>
        <span>
          Page {page} of {Math.ceil(data.total_books / data.page_size) || 1}
        </span>
        <button
          className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          disabled={
            !data.total_books ||
            page >= Math.ceil(data.total_books / data.page_size)
          }
          onClick={() => setPage((p) => p + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default BooksList;

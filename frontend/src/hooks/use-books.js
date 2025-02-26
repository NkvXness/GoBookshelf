import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { booksApi } from "@/services/api";

export function useBooks(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ["books", page],
    queryFn: () => booksApi.getBooks(page, pageSize).then((res) => res.data),
  });
}

export function useCreateBook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data) => booksApi.createBook(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books"] });
    },
  });
}

export function useUpdateBook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => booksApi.updateBook(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books"] });
    },
  });
}

export function useDeleteBook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id) => booksApi.deleteBook(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["books"] });
    },
  });
}

import { AlertTriangle, Loader2 } from "lucide-react";

const DeleteConfirmModal = ({ book, onCancel, onConfirm, isDeleting }) => {
  if (!book) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 backdrop-blur-sm">
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 max-w-md w-full shadow-xl transform transition-all animate-fadeIn">
        <div className="flex items-center justify-center mb-6">
          <div className="bg-red-100 dark:bg-red-900/30 p-3 rounded-full">
            <AlertTriangle size={36} className="text-red-600 dark:text-red-500" />
          </div>
        </div>
        
        <h3 className="text-xl font-bold text-center mb-3 text-gray-900 dark:text-white">
          Подтверждение удаления
        </h3>
        
        <p className="text-center mb-2 text-gray-600 dark:text-gray-300">
          Вы уверены, что хотите удалить книгу:
        </p>
        
        <p className="text-center font-medium mb-6 text-gray-900 dark:text-white">
          "{book.title}"
        </p>
        
        <p className="text-center text-sm mb-6 text-red-600 dark:text-red-400">
          Это действие нельзя отменить.
        </p>
        
        <div className="flex justify-center space-x-4">
          <button
            onClick={onCancel}
            disabled={isDeleting}
            className="px-5 py-2.5 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors disabled:opacity-50 font-medium"
          >
            Отмена
          </button>
          
          <button
            onClick={onConfirm}
            disabled={isDeleting}
            className="px-5 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 font-medium flex items-center"
          >
            {isDeleting ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Удаление...
              </>
            ) : (
              "Удалить"
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default DeleteConfirmModal;
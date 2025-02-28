import { createContext, useContext, useState, useCallback } from 'react';
import Toast from '../components/ui/Toast';

// Создаем контекст
const ToastContext = createContext(undefined);

// ID для новых тостов
let toastId = 1;

// Провайдер контекста тостов
export const ToastProvider = ({ children }) => {
  const [toasts, setToasts] = useState([]);

  // Добавить новый тост
  const addToast = useCallback((message, type = 'info', duration = 5000) => {
    const id = toastId++;
    setToasts(prevToasts => [...prevToasts, { id, message, type, duration }]);
    return id; // Возвращаем ID для возможности удаления
  }, []);

  // Удалить тост по ID
  const removeToast = useCallback((id) => {
    setToasts(prevToasts => prevToasts.filter(toast => toast.id !== id));
  }, []);

  // Вспомогательные функции для различных типов уведомлений
  const success = useCallback((message, duration) => 
    addToast(message, 'success', duration), [addToast]);
  
  const error = useCallback((message, duration) => 
    addToast(message, 'error', duration), [addToast]);
  
  const warning = useCallback((message, duration) => 
    addToast(message, 'warning', duration), [addToast]);
  
  const info = useCallback((message, duration) => 
    addToast(message, 'info', duration), [addToast]);

  return (
    <ToastContext.Provider value={{ addToast, removeToast, success, error, warning, info }}>
      {children}
      
      {/* Контейнер для тостов */}
      <div className="fixed top-4 right-4 z-50 flex flex-col items-end">
        {toasts.map(toast => (
          <Toast
            key={toast.id}
            id={toast.id}
            message={toast.message}
            type={toast.type}
            duration={toast.duration}
            onClose={removeToast}
          />
        ))}
      </div>
    </ToastContext.Provider>
  );
};

// Хук для использования тостов в компонентах
export const useToast = () => {
  const context = useContext(ToastContext);
  
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  
  return context;
};
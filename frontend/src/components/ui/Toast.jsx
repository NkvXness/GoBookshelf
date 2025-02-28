import { useEffect, useState } from 'react';
import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react';

// Компонент отдельного уведомления
const Toast = ({ id, message, type = 'info', duration = 5000, onClose }) => {
  const [isExiting, setIsExiting] = useState(false);
  
  // Иконки для разных типов уведомлений
  const icons = {
    success: <CheckCircle className="w-5 h-5" />,
    error: <AlertCircle className="w-5 h-5" />,
    warning: <AlertTriangle className="w-5 h-5" />,
    info: <Info className="w-5 h-5" />
  };

  // Цвета для разных типов уведомлений
  const colors = {
    success: "bg-green-50 text-green-800 border-green-200 dark:bg-green-900/30 dark:text-green-300 dark:border-green-800",
    error: "bg-red-50 text-red-800 border-red-200 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800",
    warning: "bg-amber-50 text-amber-800 border-amber-200 dark:bg-amber-900/30 dark:text-amber-300 dark:border-amber-800",
    info: "bg-blue-50 text-blue-800 border-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-800"
  };

  // Цвета иконок
  const iconColors = {
    success: "text-green-500 dark:text-green-400",
    error: "text-red-500 dark:text-red-400",
    warning: "text-amber-500 dark:text-amber-400",
    info: "text-blue-500 dark:text-blue-400"
  };

  // Автоматическое закрытие уведомления через заданное время
  useEffect(() => {
    if (duration) {
      const timer = setTimeout(() => {
        setIsExiting(true);
        setTimeout(() => onClose(id), 300); // Задержка для анимации
      }, duration);
      
      return () => clearTimeout(timer);
    }
  }, [duration, id, onClose]);

  // Функция закрытия с анимацией
  const handleClose = () => {
    setIsExiting(true);
    setTimeout(() => onClose(id), 300);
  };

  return (
    <div 
      className={`flex items-center max-w-md w-full p-4 mb-3 border rounded-lg shadow-lg 
                transition-all duration-300 transform 
                ${isExiting ? 'opacity-0 translate-x-full' : 'opacity-100 translate-x-0'} 
                ${colors[type]}`}
    >
      <div className={`flex-shrink-0 ${iconColors[type]}`}>
        {icons[type]}
      </div>
      <div className="ml-3 mr-2 flex-grow text-sm font-medium">
        {message}
      </div>
      <button 
        onClick={handleClose}
        className="ml-auto -mx-1.5 -my-1.5 rounded-lg p-1.5 inline-flex items-center justify-center h-8 w-8 
                 hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-700 dark:hover:text-white"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
};

export default Toast;
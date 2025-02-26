import { useState, useEffect } from "react";

// Хук для управления темой приложения
export function useTheme() {
  // Проверяем localStorage и системные настройки при инициализации
  const getInitialTheme = () => {
    // Проверяем сохраненную тему
    const savedTheme = localStorage.getItem("theme");

    if (savedTheme) {
      return savedTheme === "dark";
    }

    // Если нет сохраненной темы, смотрим системные настройки
    return window.matchMedia("(prefers-color-scheme: dark)").matches;
  };

  const [isDarkMode, setIsDarkMode] = useState(getInitialTheme);

  // Применяем тему к документу при изменении isDarkMode
  useEffect(() => {
    const html = document.documentElement;

    if (isDarkMode) {
      html.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      html.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  }, [isDarkMode]);

  // Функция переключения темы
  const toggleTheme = () => {
    setIsDarkMode((prev) => !prev);
  };

  return { isDarkMode, toggleTheme };
}

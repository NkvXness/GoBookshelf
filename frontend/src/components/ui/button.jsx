import { forwardRef } from "react";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

const Button = forwardRef(
  ({ className, variant = "default", size = "default", ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={twMerge(
          clsx(
            "inline-flex items-center justify-center rounded-md font-medium transition-colors",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2",
            "disabled:pointer-events-none disabled:opacity-50",
            {
              "bg-blue-600 text-white hover:bg-blue-700": variant === "default",
              "bg-red-600 text-white hover:bg-red-700":
                variant === "destructive",
              "border border-gray-300 bg-white hover:bg-gray-100":
                variant === "outline",
              "underline-offset-4 hover:underline": variant === "link",
              "h-10 px-4 py-2": size === "default",
              "h-8 px-3 text-sm": size === "sm",
              "h-12 px-8": size === "lg",
            },
            className,
          ),
        )}
        {...props}
      />
    );
  },
);

Button.displayName = "Button";

export { Button };

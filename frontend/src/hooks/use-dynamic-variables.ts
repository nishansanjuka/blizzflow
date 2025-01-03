"use client";
import { useEffect, useState } from "react";

type CSSVariables = {
  [key: string]: string | number;
};

function useDynamicVariables() {
  const [cssVariables, setCSSVariables] = useState<CSSVariables>({}); // Initialize empty

  useEffect(() => {
    const root = document.documentElement;

    // Apply CSS variables globally to :root when they are updated
    Object.entries(cssVariables).forEach(([key, value]) => {
      root.style.setProperty(`--${key}`, String(value));
    });

    // No cleanup is necessary for setting variables
  }, [cssVariables]);

  // Merge existing CSS variables with new ones
  const updateCSSVariables = (newVariables: CSSVariables) => {
    setCSSVariables((prevVariables) => ({
      ...prevVariables,
      ...newVariables,
    }));
  };

  return updateCSSVariables; // Return the merged setter for dynamic updates
}

export default useDynamicVariables;

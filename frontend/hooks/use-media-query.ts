"use client";
import { useState, useEffect } from "react";

interface MediaQueryParams {
  maxWidth?: number;
  minWidth?: number;
  orientation?: "portrait" | "landscape";
}

function useMediaQuery({ maxWidth, minWidth, orientation }: MediaQueryParams) {
  const [matches, setMatches] = useState(false);

  useEffect(() => {
    const mediaQueries = [];

    if (maxWidth) mediaQueries.push(`(max-width: ${maxWidth}px)`);
    if (minWidth) mediaQueries.push(`(min-width: ${minWidth}px)`);
    if (orientation) mediaQueries.push(`(orientation: ${orientation})`);

    const mediaQuery = mediaQueries.join(" and ");
    const mediaQueryList = window.matchMedia(mediaQuery);

    const updateMatch = () => setMatches(mediaQueryList.matches);

    updateMatch(); // Check initial state
    mediaQueryList.addEventListener("change", updateMatch); // Listen for changes

    return () => mediaQueryList.removeEventListener("change", updateMatch); // Cleanup
  }, [maxWidth, minWidth, orientation]);

  return matches;
}

export default useMediaQuery;

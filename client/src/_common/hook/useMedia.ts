import { useState, useEffect, useMemo, useCallback } from "react";

// source : https://github.com/craig1123/react-recipes/blob/master/src/useMedia.js

// sample usage: see https://github.com/craig1123/react-recipes/blob/master/src/useMedia.js

function useMedia<T>(queries: string[], values: T[], defaultValue: T) {
  // Array containing a media query list for each query
  const mediaQueryLists = useMemo(
    () => queries.map((q) => window.matchMedia(q)),
    [queries]
  );

  // Function that gets value based on matching media query
  const getValue = useCallback(() => {
    // Get index of first media query that matches
    const index = mediaQueryLists.findIndex((mql) => mql.matches);
    // Return related value or defaultValue if none
    return typeof values[index] !== "undefined" ? values[index] : defaultValue;
  }, [mediaQueryLists, defaultValue, values]);

  const [value, setValue] = useState(getValue);

  useEffect(() => {
    // Event listener callback
    // NOTE: By defining getValue outside of useEffect we ensure that it has ...
    // ... current values of hook args (as this hook callback is created once on mount).
    const handler = () => setValue(getValue);

    // Set a listener for each media query with above handler as callback.
    mediaQueryLists.forEach((mql) => mql.addEventListener("change", handler));
    return () => mediaQueryLists.forEach((mql) => mql.removeEventListener("change", handler));
  }, [mediaQueryLists, getValue]);

  return value;
}

export default useMedia;

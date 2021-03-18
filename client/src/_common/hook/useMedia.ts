import { useState, useEffect, useMemo, useCallback } from "react";

/**
 * useMedia
 * 
 * Similar as css media queries, but in React js component.
 * 
 * source : https://github.com/craig1123/react-recipes/blob/master/src/useMedia.js
 * 
 * @param queries media breaks conditions. ex : ["(max-width: 1250px)", "(min-width: 1250px)"]
 * @param values return values for each valid query
 * @param defaultValue default value
 * @returns one value of param "values" or else default value
 */
function useMedia<T>(queries: string[], values: T[], defaultValue: T) {
  const mediaQueries = useMemo(() => queries.map((q) => window.matchMedia(q)), [queries]);
  const getValue = useCallback(() => {
    const index = mediaQueries.findIndex((mql) => mql.matches);
    return typeof values[index] !== "undefined" ? values[index] : defaultValue;
  }, [mediaQueries, defaultValue, values]);

  const [value, setValue] = useState(getValue);

  useEffect(() => {
    const handler = () => setValue(getValue);
    mediaQueries.forEach((mql) => mql.addEventListener("change", handler));
    return () => mediaQueries.forEach((mql) => mql.removeEventListener("change", handler));
  }, [mediaQueries, getValue]);

  return value;
}

export default useMedia;

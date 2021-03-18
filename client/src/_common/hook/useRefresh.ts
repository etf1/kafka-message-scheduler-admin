import { useState, useCallback } from "react";

/**
 * useRefresh
 *
 * Force Component redraw
 */
const useRefresh = (): [() => void, number] => {
  const [count, setCount] = useState(0);
  const refresh = useCallback(() => {
    setCount((prevCount) => {
      return prevCount + 1;
    });
  }, []);

  return [refresh, count];
};

export default useRefresh;

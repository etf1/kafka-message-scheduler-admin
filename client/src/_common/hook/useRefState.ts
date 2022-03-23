import { useRef, useMemo } from "react";
import { isFunction } from "_common/type/utils";
import useRefresh from "./useRefresh";

export default function useRefState<S>(
  initialValueOrInitializer: S | (() => S)
): [() => S, (value: S | ((old: S) => S)) => void] {
  const [refresh] = useRefresh();
  const valueRef = useRef<S>(
    isFunction(initialValueOrInitializer)
      ? initialValueOrInitializer()
      : initialValueOrInitializer
  );
  return useMemo(
    () => [
      (): S => valueRef.current,
      (value: S | ((old: S) => S)): void => {
        valueRef.current = isFunction(value) ? value(valueRef.current) : value;
        refresh();
      },
    ],
    // eslint-disable-next-line
    []
  );
}

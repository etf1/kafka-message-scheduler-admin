/* eslint-disable @typescript-eslint/no-explicit-any */

import { useRef } from "react";
import { later } from "_common/service/FunUtil";

const useFocus = (): [React.MutableRefObject<any>, () => void] => {
  const htmlElRef = useRef<any>(null);
  const setFocus = () => {
    later().then(() => htmlElRef.current && htmlElRef.current.focus());
  };

  return [htmlElRef, setFocus];
};

export default useFocus;

import { later } from '_common/service/FunUtil';
import { useRef } from 'react';

const useFocus = (): [React.MutableRefObject<any>, () => void] => {
  const htmlElRef = useRef<any>(null);
  const setFocus = () => {
    later().then(() => htmlElRef.current && htmlElRef.current.focus());
  };

  return [htmlElRef, setFocus];
};

export default useFocus;

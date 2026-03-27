import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import React, { useEffect, useRef } from 'react';

hljs.registerLanguage('json', json);

export type ScheduleValueProps = {
  value: React.ReactNode;
};
const ScheduleValue: React.FC<ScheduleValueProps> = ({ value }) => {
  const codeRef = useRef<HTMLElement>(null);

  useEffect(() => {
    if (codeRef.current) {
      hljs.highlightElement(codeRef.current);
    }
  }, []);

  return (
    <pre>
      <code
        ref={codeRef}
        className='json'
      >
        {value}
      </code>
    </pre>
  );
};

export default ScheduleValue;

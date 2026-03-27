import React, { Suspense } from 'react';
import { createRoot } from 'react-dom/client';
import 'bulma/css/bulma.css';
import './bulma-tooltip.css';
import './index.css';
import App from '_core/app/App';

import '_core/i18n';

import init from '_core/service/config';

init().then(() => {
  const container = document.getElementById('root');
  if (!container) throw new Error('Root element not found');

  const root = createRoot(container);
  root.render(
    <React.StrictMode>
      <Suspense fallback={<div></div>}>
        <App />
      </Suspense>
    </React.StrictMode>,
  );
});

// These must be the first lines in src/index.js
import "react-app-polyfill/ie11";
import "react-app-polyfill/stable";

import React, { Suspense } from "react";
import ReactDOM from "react-dom";
import "bulma/css/bulma.css";
import "./index.css";
import App from "_core/app/App";

import "_core/i18n";

import init from "_core/service/config";

init().then(() => {
  /*if (process.env.NODE_ENV === "development") {
    const { worker } = require("./mocks/browser");
    worker.start();
  }*/

  ReactDOM.render(
    <React.StrictMode>
      <Suspense fallback={<div></div>}>
        <App />
      </Suspense>
    </React.StrictMode>,
    document.getElementById("root")
  );
});

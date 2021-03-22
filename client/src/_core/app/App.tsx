import React, { Suspense } from "react";
import "./App.css";

import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import routes from "../router/routes";
import AppNavbar from "./app-navbar/AppNavbar";

function App() {
  return (
    <Router>
      <div className="App">
        <section className="hero is-info is-fullheight app-hero">
          <div className="hero-head">
            <AppNavbar />
          </div>
          <div className="hero-body">
            <Switch>
              {routes.map((route) => {
                return (
                  <Route key={route.key} exact={route.exact} path={route.path}>
                    <Suspense fallback={<div></div>}>
                      <route.component />
                    </Suspense>
                  </Route>
                );
              })}
            </Switch>
          </div>
        </section>
      </div>
    </Router>
  );
}

export default App;

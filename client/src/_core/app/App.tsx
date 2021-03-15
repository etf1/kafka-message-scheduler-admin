import React, { Suspense, useEffect } from "react";
import "./App.css";
import Navbar, { MenuItem } from "_common/component/navbar/Navbar";

import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import routes, { ROUTE_ABOUT, ROUTE_HOME } from "../router/routes";
import { useTranslation } from "react-i18next";
import { changeLanguage } from "_core/i18n";
import AppNavbar from "./AppNavbar";

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

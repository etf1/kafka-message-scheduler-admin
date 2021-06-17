import React, { Suspense } from "react";
import Styles from "./App.module.css";

import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import routes from "../router/routes";
import AppNavbar from "./app-navbar/AppNavbar";
import AppLeftSidebar from "./app-left-sidebar/AppLeftSidebar";
import ModalProvider from "_common/component/modal/ModalProvider";

function App() {
  return (
    <>
      <Router>
        <div className={Styles.App}>
          <Switch>
            {routes.map((route) => {
              return (
                <Route key={route.key} exact={route.exact} path={route.path}>
                  <AppLeftSidebar />
                  <AppNavbar />
                  <div className={Styles.AppContainer}>
                    {" "}
                    <Suspense fallback={<div></div>}>
                      <route.component />
                    </Suspense>
                  </div>
                </Route>
              );
            })}
          </Switch>
        </div>
      </Router>

      <ModalProvider />
    </>
  );
}

export default App;

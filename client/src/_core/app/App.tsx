import { Suspense } from "react";
import Style from "./App.module.css";

import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import routes from "../router/routes";
import AppNavbar from "./app-navbar/AppNavbar";
import AppLeftSidebar from "./app-left-sidebar/AppLeftSidebar";
import ModalProvider from "_common/component/modal/ModalProvider";

function App() {
  /*useEffect( ()=> {
    changeLanguage("en-US")
  }, [])*/

  return (
    <>
      <Router>
        <Switch>
          {routes.map((route) => {
            return (
              <Route key={route.key} exact={route.exact} path={route.path}>
                <AppLeftSidebar />
                <AppNavbar />
                <main role="main" className={Style.AppContainer}>
                  <Suspense fallback={<div></div>}>
                    <route.component />
                  </Suspense>
                </main>
              </Route>
            );
          })}
        </Switch>
      </Router>
      <ModalProvider />
    </>
  );
}

export default App;

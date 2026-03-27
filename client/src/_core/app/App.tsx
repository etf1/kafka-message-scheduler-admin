import ModalProvider from '_common/component/modal/ModalProvider';
import { Suspense } from 'react';

import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';

import routes from '../router/routes';
import Style from './App.module.css';
import AppLeftSidebar from './app-left-sidebar/AppLeftSidebar';
import AppNavbar from './app-navbar/AppNavbar';

function App() {
  return (
    <>
      <Router>
        <Routes>
          {routes.map((route) => {
            const Component = route.component;
            return (
              <Route
                key={route.key}
                path={route.path}
                element={
                  <>
                    <AppLeftSidebar />
                    <AppNavbar />
                    <main
                      role='main'
                      className={Style.AppContainer}
                    >
                      <Suspense fallback={<div></div>}>
                        <Component />
                      </Suspense>
                    </main>
                  </>
                }
              />
            );
          })}
        </Routes>
      </Router>
      <ModalProvider />
    </>
  );
}

export default App;

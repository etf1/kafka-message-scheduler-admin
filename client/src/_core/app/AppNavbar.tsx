import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { ROUTE_ALL_SCHEDULES, ROUTE_HOME, ROUTE_LIVE_SCHEDULES } from "_core/router/routes";

const AppNavbar = () => {
  const { t } = useTranslation();
  
  return (
    <nav className="navbar">
      <div className="container">
        <div className="navbar-brand">
          <span className="navbar-burger burger white" data-target="navbarMenu">
            <span></span>
            <span></span>
            <span></span>
          </span>
        </div>
        <div id="navbarMenu" className="navbar-menu">
          <div className="navbar-start">
            <span className="navbar-item">
              <a className="button is-white is-outlined" href={ROUTE_HOME}>
                <span className="icon">
                  <i className="fa fa-home"></i>
                </span>
                <span>{t("Menu-home")}</span>
              </a>
            </span>
          </div>
          <div className="navbar-end">
            <span className="navbar-item">
              <Link className="button is-white is-outlined" to={ROUTE_LIVE_SCHEDULES}>
                <span className="icon">
                  <i className="fa fa-calendar"></i>
                </span>
                <span>{t("Menu-schedules-live")}</span>
              </Link>
            </span>
            <span className="navbar-item">
              <Link className="button is-white is-outlined" to={ROUTE_ALL_SCHEDULES}>
                <span className="icon">
                  <i className="fa fa-calendar-alt"></i>
                </span>
                <span>{t("Menu-schedules-all")}</span>
              </Link>
            </span>
            <span className="navbar-item">
              <a
                className="button is-white is-outlined"
                target="_blank" rel="noreferrer"
                href="https://github.com/etf1/kafka-message-scheduler-admin"
              >
                <span className="icon">
                  <i className="fab fa-github"></i>
                </span>
                <span>{t("Menu-Source")}</span>
              </a>
            </span>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default AppNavbar;

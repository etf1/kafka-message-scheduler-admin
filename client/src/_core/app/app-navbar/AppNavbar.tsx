import clsx from "clsx";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import {
  ROUTE_ALL_SCHEDULES,
  ROUTE_HOME,
  ROUTE_LIVE_SCHEDULES,
} from "_core/router/routes";
import Styles from "./AppNavbar.module.css";

const AppNavbar = () => {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  const handleBurgerClick = () => setIsOpen((isOpen) => !isOpen);

  return (
    <nav className={clsx("navbar", Styles.Navbar)}>
      <div className="container">
        <div className="navbar-brand">
          <span
            role="button"
            className={clsx(
              "navbar-burger burger white",
              isOpen ? "is-active" : null
            )}
            aria-label="menu"
            aria-expanded="false"
            data-target="navbarMenu"
            onClick={handleBurgerClick}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </span>
        </div>
        <div
          id="navbarMenu"
          className={clsx(
            "navbar-menu",
            Styles.NavbarMenu,
            isOpen ? "is-active" : null
          )}
        >
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
              <Link
                className="button is-white is-outlined"
                to={ROUTE_LIVE_SCHEDULES}
              >
                <span className="icon">
                  <i className="fa fa-calendar"></i>
                </span>
                <span>{t("Menu-schedules-live")}</span>
              </Link>
            </span>
            <span className="navbar-item">
              <Link
                className="button is-white is-outlined"
                to={ROUTE_ALL_SCHEDULES}
              >
                <span className="icon">
                  <i className="fa fa-calendar-alt"></i>
                </span>
                <span>{t("Menu-schedules-all")}</span>
              </Link>
            </span>
            <span className="navbar-item">
              <a
                className="button is-white is-outlined"
                target="_blank"
                rel="noreferrer"
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

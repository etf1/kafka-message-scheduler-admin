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
    <nav className={clsx("navbar", Styles.Nav )}>
      <div className="container">
        <div className="navbar-brand">
          <span
            role="button"
            className={clsx(
              "navbar-burger burger white",
              isOpen ? "is-active" : null, 
              Styles.NavbarMenu
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
            <span className={clsx("navbar-item", Styles.Brand)}>
              <a className={clsx("button is-white",Styles.NavbarLink)} href={ROUTE_HOME}>
              
                <span className={Styles.BrandTitle}><span style={{color:"#00b89c"}}>K</span>afka Scheduler</span> 
              </a>
            </span>
          </div>
          <div className="navbar-end">
            
            <span className={clsx("navbar-item")}>
              <a
                className={clsx("button is-white is-outlined", Styles.NavbarLink)}
                target="_blank"
                rel="noreferrer"
                href="https://github.com/etf1/kafka-message-scheduler-admin"
                style={{color:"gray"}}
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

import clsx from "clsx";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { ROUTE_HOME } from "_core/router/routes";
import Styles from "./AppNavbar.module.css";

const highlightFirstLetter = (text: string, color: string) => {
  const first = text.charAt(0);
  const rest = text.substring(1);
  return (
    <>
      <span style={{ color }}>{first}</span>
      {rest}
    </>
  );
};

const AppNavbar = () => {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  const handleBurgerClick = () => setIsOpen((isOpen) => !isOpen);

  return (
    <nav className={clsx("navbar", Styles.Nav)}>
      <div className="container">
        <div className="navbar-brand">
          <span
            role="button"
            className={clsx("navbar-burger burger white", isOpen ? "is-active" : null, Styles.NavbarMenu)}
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
        <div id="navbarMenu" className={clsx("navbar-menu", Styles.NavbarMenu, isOpen ? "is-active" : null)}>
          <div className="navbar-start">
            <span className={clsx("navbar-item", Styles.Brand)}>
              <a className={clsx("button is-white", Styles.NavbarLink)} href={ROUTE_HOME}>
                <span className={Styles.BrandTitle}>{highlightFirstLetter(t("App-title"), "#00b89c")}</span>
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
                style={{ color: "gray" }}
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

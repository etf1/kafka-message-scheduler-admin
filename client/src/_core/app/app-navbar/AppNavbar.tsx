import clsx from "clsx";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { changeLanguage, Lang } from "_core/i18n";
import { ROUTE_HOME } from "_core/router/routes";
import Style from "./AppNavbar.module.css";

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

  const setLang = (lang: Lang) => {
    changeLanguage(lang);
  };

  const getLangLabel = (lang: Lang) => {
    switch (lang) {
      case "en-US": {
        return "Menu-Display-In-English";
      }
      case "fr-FR": {
        return "Menu-Display-In-French";
      }
      default:
        return "Menu-Display-In-English";
    }
  };

  return (
    <nav className={clsx("navbar", Style.Nav)}>
      <div className="container">
        <div className="navbar-brand">
          <span
            role="button"
            className={clsx("navbar-burger burger white", isOpen ? "is-active" : null, Style.NavbarMenu)}
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
        <div id="navbarMenu" className={clsx("navbar-menu", Style.NavbarMenu, isOpen ? "is-active" : null)}>
          <div className="navbar-start">
            <span className={clsx("navbar-item", Style.Brand)}>
              <a className={clsx("button is-white", Style.NavbarLink)} href={ROUTE_HOME}>
                <span className={Style.BrandTitle}>{highlightFirstLetter(t("App-title"), "#00b89c")}</span>
              </a>
            </span>
          </div>
          <div className="navbar-end">
            <span className={clsx("navbar-item")}>
              <a
                className={clsx("button is-white is-outlined", Style.NavbarLink)}
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
            <div className={clsx("navbar-item has-dropdown is-hoverable", Style.NavbarDropdown)}>
              <label className={clsx("navbar-link", Style.NavbarLink)} style={{ color: "#5d5d5d !important" }}>
                <span className="icon">
                  <i className="fa fa-flag"></i>
                </span>
              </label>

              <div className="navbar-dropdown">
                <span
                  onClick={() => setLang("en-US")}
                  className={clsx("navbar-item", "has-tooltip-left")}
                  style={{ cursor: "pointer", paddingRight: 30 }}
                  data-tooltip={t(getLangLabel("en-US"))}
                >
                  <img src="/asset/english_flag.svg" width="32" alt={t(getLangLabel("en-US"))}/>
                </span>
                <span
                  onClick={() => setLang("fr-FR")}
                  className={clsx("navbar-item", "has-tooltip-left")}
                  style={{ cursor: "pointer", paddingRight: 30 }}
                  data-tooltip={t(getLangLabel("fr-FR"))}
                >
                  <img src="/asset/french_flag.svg" width="32" alt={t(getLangLabel("fr-FR"))}/>
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default AppNavbar;

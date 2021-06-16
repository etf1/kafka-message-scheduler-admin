import clsx from "clsx";
import React from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import Style from "./AppLeftSidebar.module.css";

const AppLeftSidebar = () => {
  const { t } = useTranslation();
  const pathname = window.location.pathname;

  return (
    <div className={clsx("menu", Style.Container)}>
      <Link data-key="menu-item" to="/home">
        <div className={clsx(Style.MenuItem, pathname === "/" ? Style.MenuItemSelected : null)}>
          <span className="icon has-tooltip-right" data-tooltip={t("Menu-home")}>
            <i className="fa fa-home fas fa-lg"></i>
          </span>
        </div>
      </Link>

      <Link data-key="menu-item" to="/schedulers">
        <div className={clsx(Style.MenuItem, pathname === "/schedulers" ? Style.MenuItemSelected : null)}>
          <span className="icon has-tooltip-right" data-tooltip={t("Menu-schedulers")}>
            <i className="fa fa-stopwatch fas fa-lg"></i>
          </span>
        </div>
      </Link>
 
      <Link data-key="menu-item" to="/live">
        <div className={clsx(Style.MenuItem, pathname === "/live" ? Style.MenuItemSelected : null)}>
          <span className="icon has-tooltip-right" data-tooltip={t("Menu-schedules-live")}>
            <i className="fa fa-calendar fas fa-lg"></i>
          </span>
        </div>
      </Link>

      <Link data-key="menu-item" to="/all">
        <div className={clsx(Style.MenuItem, pathname === "/all" ? Style.MenuItemSelected : null)}>
          <span className="icon has-tooltip-right" data-tooltip={t("Menu-schedules-all")}>
            <i className="fa fa-calendar-alt fas fa-lg"></i>
          </span>
        </div>
      </Link>
    </div>
  );
};

export default AppLeftSidebar;
/*
"Menu-home": "Accueil",
  "Menu-bout": "A propos",
  "Menu-schedules-all": "Toutes les planifications",
  "Menu-schedules-live": "Planifications actives",
*/
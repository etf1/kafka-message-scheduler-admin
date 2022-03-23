import clsx from "clsx";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import Icon from "_common/component/element/icon/Icon";
import { routesWithMenu } from "_core/router/routes";
import Style from "./AppLeftSidebar.module.css";

const bestStartPath = (paths: string[], path: string) => {
  const pathname = window.location.pathname;
  let len = 0;
  if (pathname.startsWith(path)) {
    len = path.length;
  } else {
    return false;
  }
  return !paths.find((p) => pathname.startsWith(p) && p.length > len);
};

const AppLeftSidebar = () => {
  const { t } = useTranslation();

  const allPaths = routesWithMenu.map(({ path }) => path);

  return (
    <div className={clsx("menu", Style.Container)}>
      {routesWithMenu.map(({ path, key, menu }) => {
        return (
          <Link key={key} data-key="menu-item" to={path}>
            <div
              className={clsx(
                Style.MenuItem,
                bestStartPath(allPaths, path) ? Style.MenuItemSelected : null
              )}
            >
              <Icon
                name={menu?.icon || ""}
                size="lg"
                className={clsx("has-tooltip-right", Style.Icon)}
                data-tooltip={t(menu?.label || "")}
              />
            </div>
          </Link>
        );
      })}
    </div>
  );
};

export default AppLeftSidebar;

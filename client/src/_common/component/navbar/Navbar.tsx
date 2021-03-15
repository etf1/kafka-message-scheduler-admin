import clsx from "clsx";
import React, { useState } from "react";
import { Link, useLocation } from "react-router-dom";

export type Spacer = "-";
export type MenuDropDown = {
  label: React.ReactNode;
  children: MenuItem[];
};
export type MenuItemLink = {
  label: React.ReactNode;
  href: string;
};
export type MenuItemDefault = {
  label: React.ReactNode;
  onClick?: (e: any) => void;
};

export type MenuItem = Spacer | MenuItemDefault | MenuItemLink | MenuDropDown | React.ReactNode;

export const isMenuDropDown = (value: any): value is MenuDropDown => {
  return value !== null && value.hasOwnProperty("children");
};
export const isMenuItemLink = (value: any): value is MenuItemLink => {
  return value !== null && value.hasOwnProperty("href");
};
export const isMenuDefault = (value: any): value is MenuItemDefault => {
  return value !== null && value.hasOwnProperty("label");
};

const renderItems = (data: MenuItem[], location:{pathname:string}) => {
 
  return data.map((item, index) => {
    if (item === "-") {
      return <hr key={"hr"+index} className="navbar-divider" />;
    } else if (isMenuDropDown(item)) {
      return (
        <div key={"dd"+index} className="navbar-item has-dropdown is-hoverable">
          <span className="navbar-link">{item.label}</span>

          <div className="navbar-dropdown">{renderItems(item.children, location)}</div>
        </div>
      );
    } else if (isMenuItemLink(item)) {
      return (
        <Link key={"a"+index}  to={item.href} className={clsx("navbar-item", (location.pathname === item.href) && "is-active")}>
          {item.label}
        </Link>
      );
    } else if (isMenuDefault(item)) {
      return (
        <span  key={"def"+index}  className="navbar-item as-link" onClick={item.onClick}>
          {item.label}
        </span>
      );
    } else {
      return <div  key={"custom"+index}  className="navbar-item">{item}</div>;
    }
  });
};

export type NavbarProps = {
  brand?: React.ReactNode;
  items?: MenuItem[];
  rightItems?: MenuItem[];
};
/**
 * 
 * Main Application Navbar, uses react-router-dom (Link and location)
 * 
 * @param props 
 * @returns The main application navbar
 */
const Navbar: React.FC<NavbarProps> = ({ brand, items, rightItems }) => {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);

  const handleBurgerClick = () => setIsOpen((isOpen) => !isOpen);

  return (
    <nav className="navbar" role="navigation" aria-label="main navigation">
      <div className="navbar-brand">
        {brand}

        <span
          role="button"
          className={clsx("navbar-burger", isOpen ? "is-active" : null)}
          aria-label="menu"
          aria-expanded="false"
          data-target="navbarBasicExample"
          onClick={handleBurgerClick}
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </span>
      </div>

      <div id="navbarBasicExample" className={clsx("navbar-menu", isOpen ? "is-active" : null)}>
        <div className="navbar-start">{items && renderItems(items, location)}</div>
        <div className="navbar-end">{rightItems && renderItems(rightItems, location)}</div>
      </div>
    </nav>
  );
};

export default Navbar;

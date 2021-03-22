import React from "react";
import Style from "./Calendar.module.css";
import { CalendarTheme } from "./Calendar";
import { format } from "date-fns";
import clsx from "clsx";

type CalendarNavProps = {
  date: Date;
  theme: CalendarTheme;
  locale: Locale;
  onAddMonth: () => void;
  onSubMonth: () => void;
};

const CalendarNav = ({
  date,
  theme,
  locale,
  onAddMonth,
  onSubMonth,
}: CalendarNavProps) => {
  const width = `${theme.cellsWidth * 7 + 2}px`;

  const formatDate = (date: Date) => {
    try {
      return format(date, "MMMM yyyy", { locale });
    } catch (err) {
      return "";
    }
  };

  return (
    <div
      className={clsx("calendar-nav", Style.CalendarNav)}
      style={{
        width,
        lineHeight: theme.cellsWidth - theme.cellsPadding * 2 + "px",
        backgroundColor: theme.primaryColor,
      }}
    >
      <button
        onClick={onSubMonth}
        className="calendar-nav-previous button is-small is-text"
        style={{
          backgroundColor: "transparent",
          marginLeft: 5,
          boxShadow: "none",
          textDecoration: "none",
        }}
      >
        <span className="icon " style={{ color: "white" }}>
          <i className="fas fa-chevron-left" aria-hidden="true"></i>
        </span>
      </button>
      <div className="calendar-nav-month-year" style={{ display: "flex" }}>
        {formatDate(date)}
      </div>
      <button
        onClick={onAddMonth}
        className="calendar-nav-next button is-small is-text"
        style={{
          backgroundColor: "transparent",
          marginRight: 5,
          boxShadow: "none",
          textDecoration: "none",
        }}
      >
        <span className="icon " style={{ color: "white" }}>
          <i className="fas fa-chevron-right" aria-hidden="true"></i>
        </span>
      </button>
    </div>
  );
};

export default CalendarNav;

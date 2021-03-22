import React, { useState } from "react";
import { DayOfMonth } from "_common/service/DateUtil";
import clsx from "clsx";
import { CalendarTheme } from "./Calendar";
import { isSameDay } from "date-fns";

type CalendarDayProps = {
  day: DayOfMonth;
  theme: CalendarTheme;
  onClick?: (day: DayOfMonth) => void;
  selection: Date[];
};

const CalendarDay = ({ day, theme, onClick, selection }: CalendarDayProps) => {
  const [isHover, setIsOver] = useState(false);

  const toggleHover = () => setIsOver((isHover) => !isHover);
  const isSelectedDay = selection.find((d) => isSameDay(d, day.date));

  return (
    <div
      className={clsx("calendar-day", day.isToday && "is-today")}
      style={{
        textAlign: "center",
        padding: theme.cellsPadding,
        width: theme.cellsWidth + "px",
        backgroundColor: day.isThisMonth ? "#fff" : "#f5f5f5",
      }}
    >
      <button
        className="button is-white"
        onMouseOver={toggleHover}
        onMouseOut={toggleHover}
        onClick={() => onClick && onClick(day)}
        style={{
          backgroundColor: isSelectedDay ? theme.primaryColor : "transparent",
          borderRadius:
            day.isThisMonth || day.isToday || isHover || isSelectedDay
              ? theme.cellsBorderRadius
              : 0,
          width: theme.cellsWidth - theme.cellsPadding * 2 + "px",
          height: theme.cellsWidth - theme.cellsPadding * 2 + "px",
          border: isHover || day.isToday ? theme.border : "none",
          fontSize: theme.fontSize,
          textDecoration: "none",
          textAlign: "center",
          fontWeight: day.isToday ? "bold" : "normal",
          color: isSelectedDay
            ? "#fff"
            : day.isToday
            ? theme.primaryColor
            : "#333",
        }}
      >
        {day.date.getDate()}
      </button>
    </div>
  );
};

export default CalendarDay;

import React, { CSSProperties } from "react";
import Style from "./Calendar.module.css";
import { Locale, subMonths, addMonths } from "date-fns";
import {
  getDayLabelsOfWeek,
  getDaysOfMonth,
  DayOfMonth,
} from "_common/service/DateUtil";
import CalendarDay from "./CalendarDay";
import CalendarNav from "./CalendarNav";
import clsx from "clsx";
import useStateWithUpdate from "_common/hook/useStateWithUpdate";

// sources : https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e

export type CalendarTheme = {
  fontSize: string;
  primaryColor: string;
  border: string;
  cellsWidth: number;
  cellsPadding: number;
  cellsBorderRadius: number;
};
type CalendarThemeProps = Partial<CalendarTheme>; // see https://www.typescriptlang.org/docs/handbook/release-notes/typescript-2-1.html
type CalendarProps = /*HTMLAttributes<HTMLDivElement> &*/ {
  date: Date;
  locale: Locale;
  todayLabel?: string;
  theme?: CalendarThemeProps;
  onDayClick?: (day: DayOfMonth) => void;
  position?: { top: number; left: number };
};

const defaultTheme: CalendarTheme = {
  fontSize: "11px",
  primaryColor: "rgb(0, 209, 178)",
  border: "#ddd thin solid",
  cellsPadding: 2,
  cellsWidth: 36,
  cellsBorderRadius: 36,
};

const Calendar = React.forwardRef<HTMLDivElement, CalendarProps>(
  (
    {
      date,
      locale,
      theme: inputTheme,
      onDayClick,
      position,
      todayLabel,
    }: CalendarProps,
    ref
  ) => {
    const [currentDate, setCurrentDate] = useStateWithUpdate(date);

    const theme: CalendarTheme = Object.assign(defaultTheme, inputTheme || {});

    const days = getDaysOfMonth(currentDate, locale);
    const labels = getDayLabelsOfWeek(locale);

    const width = `${theme.cellsWidth * 7 + 2}px`;
    const gridTemplateColumns = `${theme.cellsWidth}px ${theme.cellsWidth}px ${theme.cellsWidth}px ${theme.cellsWidth}px ${theme.cellsWidth}px ${theme.cellsWidth}px ${theme.cellsWidth}px`;

    const handleSubMonth = () => {
      setCurrentDate((currentDate) => subMonths(currentDate, 1));
    };
    const handleAddMonth = () => {
      setCurrentDate((currentDate) => addMonths(currentDate, 1));
    };

    const handleTodayClick = () => {
      onDayClick &&
        onDayClick({
          date: new Date(),
          isToday: true,
          isThisMonth: true,
        });
    };

    let style: CSSProperties = { width };
    if (position) {
      style = {
        ...style,
        position: "absolute",
        top: position.top,
        left: position.left,
      };
    }
    return (
      <div
        className={clsx("calendar-container", Style.CalendarContainer)}
        style={style}
        ref={ref}
      >
        <CalendarNav
          date={currentDate}
          onAddMonth={handleAddMonth}
          onSubMonth={handleSubMonth}
          locale={locale}
          theme={theme}
        />
        <div
          className={clsx("calendar-header", Style.CalendarHeader)}
          style={{
            width,
            gridTemplateColumns,
            border: theme.border,
          }}
        >
          {labels.map((day) => (
            <div
              key={day}
              className="calendar-date"
              style={{
                textAlign: "center",
                padding: theme.cellsPadding,
                fontSize: theme.fontSize,
                textDecoration: "none",
                color: theme.primaryColor,
                lineHeight: `${theme.cellsWidth - 8}px`,
              }}
            >
              {day}
            </div>
          ))}
        </div>
        <div
          className={clsx("calendar-body", Style.CalendarBody)}
          style={{
            width,
            gridTemplateColumns,
            border: theme.border,
          }}
        >
          {days.map((day) => (
            <CalendarDay
              key={day.date.toString()}
              day={day}
              theme={theme}
              onClick={onDayClick}
              selection={[date]}
            />
          ))}
        </div>
        <div className={Style.TodayLinkButton} onClick={handleTodayClick}>
          {todayLabel ? todayLabel : "Today"}
        </div>
      </div>
    );
  }
);

export default Calendar;

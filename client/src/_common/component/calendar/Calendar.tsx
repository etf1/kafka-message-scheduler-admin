import React, { useEffect, useRef } from "react";
import bulmaCalendar from "bulma-calendar";
import "bulma-calendar/dist/css/bulma-calendar.min.css";
import "./Calendar.css";
import { useTranslation } from "react-i18next";
import { getShortLanguageFromLS } from "_core/i18n";
import parse from 'date-fns/parse'
const withoutTime = function (date:Date) {
  var d = new Date(date);
  d.setHours(0, 0, 0, 0);
  return d;
}

export type CalendarProps = {
  uid:string;
  className: string;
  onChange: (d:Date) => void;
  value: Date | undefined
};
function Calendar({ uid, className, onChange, value }: CalendarProps) {
  const refValue = useRef(value);
  useEffect( ()=>{
    refValue.current = value;
  },[value])
  const { t } = useTranslation();

  const lang = getShortLanguageFromLS();

  useEffect(() => {
   
    // Initialize all input of date type.
    bulmaCalendar.attach(`#${"cal"+uid}[type="date"]`, {type:"date", startDate: refValue.current});

    // eslint-disable-next-line no-undef
    const element = document.querySelector(`#${"cal"+uid}`);
    if (element) {
      (element as any).bulmaCalendar.on("select", (datepicker: any) => {
        const dt = parse(datepicker.data.value(), t("Calendar-date-format"), new Date());
        if (!refValue.current || (withoutTime(dt).getTime() !== withoutTime(refValue.current).getTime())){
          onChange(dt);
        }
      });
    }
  }, [t, uid, onChange]);

  return (
    <div className={className}>
      <input
        id={"cal"+uid}
        type="date"
        data-display-mode="default"
        data-date-format={t('Calendar-date-format').toUpperCase()}
        data-show-header="false"
        data-lang={(lang && lang.substring(0, 2)) || "en"}
        data-today-label={t("Calendar-Today-btn-label")}
        data-clear-label={t("Calendar-Clear-btn-label")}
        data-cancel-label={t("Calendar-Cancel-btn-label")}
      />
    </div>
  );
}

export default Calendar;

import React  from "react";
import Calendar from "./Calendar";
import clsx from "clsx";
import { enGB } from "date-fns/locale";
import { format } from "date-fns";
import Style from "./DatePicker.module.css";
import usePopup from "_common/hook/usePopup";
import Control from "_common/component/element/Control";
export type DatePickerHandler = (date: Date | undefined) => void;

type DatePickerProps = {
  locale?: Locale;
  value: Date | undefined;
  dateFormat?: string;
  todayLabel?: string;
  className?: string;
  onChange?: DatePickerHandler;
  isError?: boolean;
  isUp?: boolean;
  isRight?: boolean;
  isSmall?: boolean;
  placeholder?: string;
  disabled?: boolean;
};

function DatePicker({
  locale,
  value,
  dateFormat,
  todayLabel,
  isSmall,
  className,
  onChange,
  isError,
  placeholder,
  isUp,
  isRight,
  disabled,
}: DatePickerProps) {
  const { popupVisible, setPopupVisible, popupRef } = usePopup<HTMLDivElement>(false);

  const handleItemClick = (item: Date) => {
    setPopupVisible(false);
    onChange && onChange(item);
  };
  const toogleOpen = () => {
    if (!disabled) {
      setPopupVisible(!popupVisible);
    }
  };

  const btnStyle = isError ? { borderColor: "#f14668" } : {};

  const formatDate = (value: any) => {
    try {
      return (value && format(value, dateFormat || "MM/dd/yyyy")) || "";
    } catch (err) {
      return "";
    }
  };

  const deleteIconProps = disabled
    ? {}
    : {
        rightIcon: (
          <span className="icon" style={{ height: 34, color: "#dc8080" }}>
            <i className="fas fa-times" aria-hidden="true"></i>
          </span>
        ),
        rightIconClassName: Style.DeleteIcon,
        onRightIconClick: (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
          e.stopPropagation();
          onChange && onChange(undefined);
        },
        leftIcon: (
          <span className="icon" style={{ height: 34 }}>
            <i className="fas fa-calendar-alt" aria-hidden="true"></i>
          </span>
        ),
      };

  return (
    <div className={clsx("dropdown", popupVisible && "is-active", className, isRight && "is-right", isUp && "is-up")}>
      <div className="dropdown-trigger">
        <div aria-haspopup="true" aria-controls="dropdown-menu" style={btnStyle}>
          <div className="field is-grouped is-grouped-multiline has-addons" style={{ minWidth: 160, minHeight: 30 }}>
            <Control style={{ marginRight: 0 }} onClick={toogleOpen} {...deleteIconProps}>
              <input
                placeholder={placeholder}
                value={formatDate(value)}
                className={clsx("input", className, isError && "is-danger", isSmall && "is-small", Style.Input)}
                style={{
                  backgroundColor: disabled ? "rgb(245, 245, 245)" : "#fff",
                  cursor: disabled ? "not-allowed" : "pointer",
                }}
                readOnly
              />
            </Control>
          </div>
        </div>
      </div>
      {!disabled && (
        <div className="dropdown-menu" role="menu" ref={popupRef} style={{ paddingTop: 0 }}>
          <div className={clsx("dropdown-content", Style.DropDownContent)}>
            <Calendar
              ref={popupRef}
              date={value || new Date()}
              locale={locale || enGB}
              todayLabel={todayLabel}
              onDayClick={(d) => handleItemClick(d.date)}
            />
          </div>
        </div>
      )}
    </div>
  );
}

export default DatePicker;

import clsx from "clsx";
import React, { useState } from "react";

let uid: number = 0;

export type DropdownProps<T> = {
  placeholder: string;
  value: T | undefined;
  options: T[];
  getKey: (option: T, index: number) => string;
  renderOption: (option: T) => React.ReactElement;
  onChange: (option: T) => void;
};

function Dropdown<T>({ placeholder, value, options, getKey, renderOption, onChange }: DropdownProps<T>) {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const compId = `dropdown-menu${uid++}`;

  const handleTriggerBtnClick = () => {
    setIsOpen(true);
  };

  const handleOptionClick = (option: T) => {
    onChange(option);
    setIsOpen(false);
  };

  return (
    <div className={clsx("dropdown", isOpen && "is-active")}>
      <div className="dropdown-trigger">
        <button className="button" aria-haspopup="true" aria-controls={compId} onClick={handleTriggerBtnClick}>
          <span style={{ minWidth: 100, maxWidth:280, display:"block", overflow:"hidden" }}>{(value && renderOption(value)) || placeholder}</span>
          <span className="icon is-small" style={{    position: "inherit"}}>
            <i className="fas fa-angle-down" aria-hidden="true"></i>
          </span>
        </button>
      </div>
      <div className="dropdown-menu" id={compId} role="menu">
        <div className="dropdown-content">
          {options.map((option, index) => {
            return (
              <React.Fragment key={getKey(option, index)}>
                {index > 0 && <hr className="dropdown-divider" />}
                <div className="dropdown-item pointer" onClick={() => handleOptionClick(option)}>
                  {renderOption(option)}
                </div>
              </React.Fragment>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default Dropdown;

import React from "react";

export type IconLabelItem = {
  icon: string;
  label: React.ReactNode;
};

export type IconLabelProps = {
  data: IconLabelItem[];
};

const IconLabel: React.FC<IconLabelProps> = ({ data }) => {
  return (
    <span className="icon-text">
      {data.map(({ icon, label }: IconLabelItem, index: number) => {
        return (
          <React.Fragment key={index}>
            <span className="icon">
              <i className={`fas fa-${icon}`}></i>
            </span>
            <span>{label}</span>
          </React.Fragment>
        );
      })}
    </span>
  );
};

export default IconLabel;

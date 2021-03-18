import React from "react";

export type IconLabelItem = {
  icon: string;
  label: React.ReactNode;
};
export type IconLabelItems = {
  data: IconLabelItem[];
};

export type IconLabelProps = IconLabelItems | IconLabelItem;
export function isIconLabelItems(
  value: IconLabelItems | IconLabelItem
): value is IconLabelItems {
  return value.hasOwnProperty("data");
}

const IconLabel: React.FC<IconLabelProps> = (props) => {
  let items: IconLabelItem[] = isIconLabelItems(props) ? props.data : [props];
  return (
    <span className="icon-text">
      {items.map(({ icon, label }: IconLabelItem, index: number) => {
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

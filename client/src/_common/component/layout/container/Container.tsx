import clsx from "clsx";
import React from "react";
import Style from "./Container.module.css";

export type ContainerProps = {
  title?: React.ReactNode;
  size?: 12 | 10 | 8;
};
const Container: React.FC<ContainerProps> = ({
  title,
  size = 12,
  children
}) => {
  let offset = "";
  if (size === 8) {
    offset = "is-offset-2";
  } else if (size === 10) {
    offset = "is-offset-1";
  }
  return (
    <div className="container">
      <div
        className={clsx(
          "column",
          size ? "is-" + size : null,
          offset,
          Style.Column
        )}
      >
        {title && <h3 className={clsx("title is-5", Style.Title)} >{title}</h3>}
        {children}
      </div>
    </div>
  );
};

export default Container;

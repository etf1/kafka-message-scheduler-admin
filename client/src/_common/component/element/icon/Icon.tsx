import clsx from "clsx";
import React from "react";
import { slsx } from "_common/service/FunUtil";
import { Dictionary } from "_common/type/utils";
import "./Icon.css";

export type IconProps = React.HTMLAttributes<HTMLSpanElement> & {
  name: string;
  rotated?: "0" | "45" | "90" | "180";
  size?: "" | "lg" | "2x" | "3x";
  isLeft?: boolean;
  isRight?: boolean;
  isSmall?: boolean;
  marginRight?: string | number;
  marginLeft?: string | number;
};

const Icon = ({
  name,
  isLeft,
  isRight,
  isSmall,
  className,
  rotated,
  size,
  style,
  marginRight,
  marginLeft,
  ...otherProps
}: IconProps): React.ReactElement => {
  const dataTransform = {} as Dictionary<string>;
  if (rotated) {
    dataTransform["data-fa-transform"] = `rotate-${rotated}`;
  }
  return (
    <span
      key={name + className + rotated + size}
      className={clsx(
        "icon defaultSize",
        isLeft ? "is-left" : "",
        isRight ? "is-right" : "",
        isSmall ? "is-small" : "",

        className
      )}
      style={slsx({borderBottomStyle:"none"}, style, { marginLeft }, { marginRight })}
      {...otherProps}
    >
      <i
        className={clsx(`fas fa-${name}`, size ? `fa-${size}` : "")}
        {...dataTransform}
      ></i>
    </span>
  );
};

export default Icon;

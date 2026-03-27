import clsx from 'clsx';
import React, { CSSProperties } from 'react';

export type ControlProps = {
  children?: React.ReactNode;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  isSmall?: boolean;
  style?: CSSProperties;
  className?: string;
  leftIconStyle?: CSSProperties;
  leftIconClassName?: string;
  rightIconStyle?: CSSProperties;
  rightIconClassName?: string;
  onRightIconClick?: (event: React.MouseEvent<HTMLSpanElement, MouseEvent>) => void;
  onClick?: ((event: React.MouseEvent<HTMLElement, MouseEvent>) => void) | undefined;
  highlighted?: boolean;
};
const Control: React.FC<ControlProps> = ({
  children,
  leftIcon: LeftIcon,
  rightIcon: RightIcon,
  isSmall,
  style,
  className,
  leftIconStyle,
  leftIconClassName,
  rightIconStyle,
  rightIconClassName,
  onRightIconClick,
  onClick,
}) => {
  return (
    <div
      className={clsx(
        'control',
        LeftIcon && 'has-icons-left',
        RightIcon && 'has-icons-right',
        className,
        isSmall && 'is-small',
      )}
      style={style}
      onClick={onClick}
    >
      {children}
      {LeftIcon && (
        <span
          className={clsx('icon is-small is-left', leftIconClassName)}
          style={leftIconStyle}
        >
          {LeftIcon}
        </span>
      )}
      {RightIcon && (
        <span
          className={clsx('icon is-small is-right', rightIconClassName)}
          style={rightIconStyle}
          onClick={(e) => {
            e.persist();
            onRightIconClick && onRightIconClick(e);
          }}
        >
          {RightIcon}
        </span>
      )}
    </div>
  );
};

export default Control;

import clsx from "clsx";
import React, { useState } from "react";
import Appear from "_common/component/transition/Appear";
import Icon from "_common/component/element/icon/Icon";
import Style from "./Panel.module.css";

type PanelProps = {
  title?: React.ReactNode;
  icon?: string;
  iconStyle?: React.CSSProperties;
  rightHeader?: React.ReactNode;
  allowCollapse?: boolean;
  className?: string;
};

const Panel: React.FC<PanelProps> = ({
  title,
  icon,
  iconStyle,
  rightHeader,
  className,
  allowCollapse = false,
  children,
  ...restProps
}) => {
  const [isDown, setIsDown] = useState<boolean>(true);
  const handleClick = () => {
    allowCollapse && setIsDown((isDown) => !isDown);
  };

  return (
    <div className={clsx("box", Style.Panel, className)} {...restProps}>
      <div className="columns">
        <div className="column" onClick={handleClick}>
          <p className={clsx("title is-4", Style.Title)}>
            {icon && (
              <Icon
                name={icon}
                className={Style.TitleIcon}
                size="lg"
                style={iconStyle}
              />
            )}
            <Appear visible={!!title}>
              {(nodeRef) => (
                <span ref={nodeRef} className="ml5">
                  {title}
                </span>
              )}
            </Appear>
          </p>
        </div>
        {rightHeader && <div className="column is-narrow">{rightHeader}</div>}
        {allowCollapse && (
          <div
            className={clsx("column is-narrow", Style.CollapseIcon)}
            onClick={handleClick}
          >
            <Icon name={isDown ? "chevron-up" : "chevron-down"} />
          </div>
        )}
      </div>
      <Appear visible={!!(isDown && React.Children.count(children) > 0)}>
        {(nodeRef) => <div ref={nodeRef}>{children}</div>}
      </Appear>
    </div>
  );
};

export default Panel;

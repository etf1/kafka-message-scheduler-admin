import React from "react";
import { CSSTransition } from "react-transition-group";
import Style from "./Appear.module.css";

export type AppearProps = {
  visible: boolean | undefined;
  timeout?: number;
  fadeMore?: boolean;
  children?: (nodeRef: React.MutableRefObject<null>) => React.ReactNode;
};
const Appear = ({ visible, timeout, fadeMore, children }: AppearProps) => {
  const nodeRef = React.useRef(null);
  return (
    <CSSTransition
      in={visible}
      timeout={timeout || 2000}
      nodeRef={nodeRef}
      classNames={{
        enter: Style.enter,
        enterActive: fadeMore ? Style.enterMoreActive : Style.enterActive,
        exit: Style.exit,
        exitActive: fadeMore ? Style.exitMoreActive : Style.exitActive
      }}
    >
      {children && children(nodeRef)}
    </CSSTransition>
  );
};

export default Appear;

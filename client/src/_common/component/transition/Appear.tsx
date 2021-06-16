import React from "react";
import { CSSTransition } from "react-transition-group";
import Styles from "./Appear.module.css";

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
        enter: Styles.enter,
        enterActive: fadeMore ? Styles.enterMoreActive : Styles.enterActive,
        exit: Styles.exit,
        exitActive: fadeMore ? Styles.exitMoreActive : Styles.exitActive
      }}
    >
      {children && children(nodeRef)}
    </CSSTransition>
  );
};

export default Appear;

import clsx from "clsx";
import Styles from "./Container.module.css";

export type ContainerProps = {
  title: string;
  size?: 12 | 10 | 8;
};
const Container: React.FC<ContainerProps> = ({
  title,
  size = 12,
  children,
}) => {
  let offset = "";
  if (size === 8) {
    offset = "is-offset-2";
  } else if (size === 10) {
    offset = "is-offset-1";
  }
  return (
    <div className="container has-text-centered">
      <div
        className={clsx(
          "column",
          size ? "is-" + size : null,
          offset,
          Styles.Column
        )}
      >
        <h1 className="title">{title}</h1>
        {children}
      </div>
    </div>
  );
};

export default Container;

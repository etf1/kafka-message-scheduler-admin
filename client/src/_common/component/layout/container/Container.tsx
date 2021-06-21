import clsx from "clsx";
import Style from "./Container.module.css";

export type ContainerProps = {
  title: string;
  size?: 12 | 10 | 8;
  ref?: React.LegacyRef<HTMLTableElement> | undefined;
};
const Container: React.FC<ContainerProps> = ({
  title,
  size = 12,
  children,
  ref,
}) => {
  let offset = "";
  if (size === 8) {
    offset = "is-offset-2";
  } else if (size === 10) {
    offset = "is-offset-1";
  }
  return (
    <div className="container" ref={ref}>
      <div
        className={clsx(
          "column",
          size ? "is-" + size : null,
          offset,
          Style.Column
        )}
      >
        <h3 className="title is-5">{title}</h3>
        {children}
      </div>
    </div>
  );
};

export default Container;

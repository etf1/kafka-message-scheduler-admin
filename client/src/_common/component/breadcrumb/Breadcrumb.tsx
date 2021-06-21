import { Link } from "react-router-dom";

export type BreadcrumbProps = {
  data: { label: string; linkTo: string }[];
};

const Breadcrumb: React.FC<BreadcrumbProps> = ({ data }) => {
  const len = data.length;

  return (
    <nav className="breadcrumb" aria-label="breadcrumbs" style={{ marginLeft: "1rem" }}>
      <ul>
        {data.map((u, index) => {
          return index < len - 1 ? (
            <li key={u.linkTo}>
              <Link to={u.linkTo}>{u.label}</Link>
            </li>
          ) : (
            <li key={u.linkTo} className="is-active">
              <Link to={u.linkTo}>{u.label}</Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
};


export default Breadcrumb;
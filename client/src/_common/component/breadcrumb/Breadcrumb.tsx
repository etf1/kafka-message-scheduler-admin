export type BreadcrumbProps = {
  data: { label: string; url: string }[];
};

const Breadcrumb: React.FC<BreadcrumbProps> = ({ data }) => {
  const len = data.length;

  return (
    <nav className="breadcrumb" aria-label="breadcrumbs" style={{ marginLeft: "1rem" }}>
      <ul>
        {data.map((u, index) => {
          return index < len - 1 ? (
            <li>
              <a href={u.url}>{u.label}</a>
            </li>
          ) : (
            <li className="is-active">
              <a href={u.url}>{u.label}</a>
            </li>
          );
        })}
      </ul>
    </nav>
  );
};


export default Breadcrumb;
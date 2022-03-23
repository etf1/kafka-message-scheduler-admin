import { useTranslation } from "react-i18next";
const Loader = () => {
  const { t } = useTranslation();
  return (
    <strong className="animate-opacity gray italic more-space-top text_center block min-width-100">
      {t("Loading")}
    </strong>
  );
};

export default Loader;

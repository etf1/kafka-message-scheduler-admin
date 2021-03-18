import { useTranslation } from "react-i18next";

const Home = () => {
  const { t } = useTranslation();

  return <div>{t("About-page-title")}</div>;
};

export default Home;

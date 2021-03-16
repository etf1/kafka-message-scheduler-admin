import SearchScheduler from "business/scheduler/component/SearchScheduler";
import { useTranslation } from "react-i18next";

const SchedulesAll = () => {
  const { t } = useTranslation();
  return (
    <div className="container has-text-centered">
      <div className="column is-10 is-offset-1">
        <h1 className="title">{t("Page-title-schedules-all")}</h1>

        <SearchScheduler live={false} />
      </div>
    </div>
  );
};

export default SchedulesAll;

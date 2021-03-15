import SearchScheduler from "business/scheduler/component/SearchScheduler";
import React from "react";
import { useTranslation } from "react-i18next";

const AllSchedules = () => {
    const { t } = useTranslation();
    return <div className="container has-text-centered">
    <div className="column is-10 is-offset-1">
      <h1 className="title">{t("All-schedules-page-title")}</h1>
      <h2 className="subtitle">Fill some criteria to view live schedules...</h2>
      <div className="app-box">
        <SearchScheduler live={false} />
      </div>
    </div>
  </div>
}

export default AllSchedules;
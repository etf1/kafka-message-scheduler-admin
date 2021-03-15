import SearchScheduler from "business/scheduler/component/SearchScheduler";
import React from "react";
import { useTranslation } from "react-i18next";

const LiveSchedules = () => {
    const { t } = useTranslation();
    return <div className="container has-text-centered">
    <div className="column is-10 is-offset-1">
      <h1 className="title">{t("Live-schedules-page-title")}</h1>
      <h2 className="subtitle">Fill some criteria to view live schedules...</h2>
      <div className="app-box">
        <SearchScheduler live={true} />
      </div>
    </div>
  </div>
}

export default LiveSchedules;
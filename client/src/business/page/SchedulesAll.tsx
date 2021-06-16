import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import React from "react";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";

const SchedulesAll = () => {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const schedulerName = urlParams.get("schedulerName") || undefined;
  const scheduleId = urlParams.get("scheduleId") || undefined;
  const epochFrom = urlParams.get("epochFrom");
  const epochTo = urlParams.get("epochTo");

  return (
    <Panel icon={"calendar-alt"} title={t("Page-title-schedules-all")}>
      <SearchScheduler
        live={false}
        schedulerName={schedulerName}
        scheduleId={scheduleId}
        epochFrom={
          (epochFrom &&
            parse(epochFrom, t("Calendar-date-format"), new Date())) ||
          undefined
        }
        epochTo={
          (epochTo && parse(epochTo, t("Calendar-date-format"), new Date())) ||
          undefined
        }
      />
    </Panel>
  );
};

export default SchedulesAll;

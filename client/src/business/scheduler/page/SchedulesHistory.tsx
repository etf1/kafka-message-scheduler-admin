import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";
import endOfDay from "date-fns/endOfDay";
import { clear, load } from "_common/service/SessionStorageService";

export type SchedulesUrlParams = {
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: string;
  epochTo?: string;
};

const SchedulesHistory = () => {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const schedulerName =
    urlParams.get("schedulerName") || load("historySchedulerName", undefined);
  const scheduleId =
    urlParams.get("scheduleId") || load("historyScheduleId", undefined);
  const epochFrom =
    urlParams.get("epochFrom") || load("historyEpochFrom", undefined);
  const epochTo = urlParams.get("epochTo") || load("historyEpochTo", undefined);
  clear((key) => {
    return key.indexOf("history") === 0;
  });

  return (
    <Panel icon={"history"} title={t("Page-title-schedules-history")}>
      <SearchScheduler
        scheduleType={"history"}
        schedulerName={schedulerName}
        scheduleId={scheduleId}
        epochFrom={
          (epochFrom &&
            parse(epochFrom, t("Calendar-date-format"), new Date())) ||
          undefined
        }
        epochTo={
          (epochTo &&
            endOfDay(parse(epochTo, t("Calendar-date-format"), new Date()))) ||
          undefined
        }
      />
    </Panel>
  );
};

export default SchedulesHistory;

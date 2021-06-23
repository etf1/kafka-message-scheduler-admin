import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";
import endOfDay from "date-fns/endOfDay";
import { clear, load } from "_common/service/SessionStorageService";
import startOfDay from "date-fns/startOfDay";
import format from "date-fns/format";

const SchedulesAll = () => {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const schedulerName = urlParams.get("schedulerName") || load("AllSchedulerName", undefined);
  const scheduleId = urlParams.get("scheduleId") || load("AllScheduleId", undefined);
  const epochFrom = urlParams.get("epochFrom") || load("AllEpochFrom",  format(startOfDay(new Date()), t("Calendar-date-format")));
  const epochTo = urlParams.get("epochTo") || load("AllEpochTo", format(endOfDay(new Date()), t("Calendar-date-format")));
  clear( (key) => {
    return key.indexOf("All") ===0;
  });

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
          (epochTo && endOfDay(parse(epochTo, t("Calendar-date-format"), new Date()))) ||
          undefined
        }
      />
    </Panel>
  );
};

export default SchedulesAll;

import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";
import endOfDay from "date-fns/endOfDay";
import { load } from "_common/service/SessionStorageService";

export type SchedulesUrlParams = {
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: string;
  epochTo?: string;
};

const SchedulesLive = () => {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const schedulerName = urlParams.get("schedulerName") || load("schedulerNameLive", undefined);
  const scheduleId = urlParams.get("scheduleId") || load("scheduleIdLive", undefined);
  const epochFrom = urlParams.get("epochFrom") || load("epochFromLive", undefined);
  const epochTo = urlParams.get("epochTo") || load("epochToLive", undefined);

  return (
    <Panel icon={"calendar"} title={t("Page-title-schedules-live")}>
      <SearchScheduler
        live={true}
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

export default SchedulesLive;

import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";
import endOfDay from "date-fns/endOfDay";

export type SchedulesUrlParams = {
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: string;
  epochTo?: string;
};

const SchedulesLive = () => {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const schedulerName = urlParams.get("schedulerName") || undefined;
  const scheduleId = urlParams.get("scheduleId") || undefined;
  const epochFrom = urlParams.get("epochFrom");
  const epochTo = urlParams.get("epochTo");

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

import React from "react";
import SearchScheduler from "business/scheduler/component/SearchScheduler";
import parse from "date-fns/parse";
import { useTranslation } from "react-i18next";
import Container from "_common/component/layout/Container";

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
    <Container title={t("Page-title-schedules-live")}>
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
          (epochTo && parse(epochTo, t("Calendar-date-format"), new Date())) ||
          undefined
        }
      />
    </Container>
  );
};

export default SchedulesLive;

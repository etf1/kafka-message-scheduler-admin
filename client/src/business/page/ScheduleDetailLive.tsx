import ScheduleForm from "business/scheduler/component/ScheduleForm";
import React from "react";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import Container from "_common/component/layout/Container";

type ScheduleDetailLiveUrlParams = { schedulerName: string; scheduleId: string };

const ScheduleDetailLive = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const handleClose = () => {
    history.goBack();
  };

  const { schedulerName, scheduleId } = useParams<ScheduleDetailLiveUrlParams>();

  return (
    <Container
      size={8}
      title={t("Page-title-schedule-detail", { id: scheduleId })}
    >
      <ScheduleForm
        schedulerName={schedulerName}
        scheduleId={scheduleId}
        onClose={handleClose}
        live={true}
      />
    </Container>
  );
};

export default ScheduleDetailLive;

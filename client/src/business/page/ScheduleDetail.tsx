import ScheduleForm from "business/scheduler/component/ScheduleForm";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import Breadcrumb from "_common/component/breadcrumb/Breadcrumb";
import Panel from "_common/component/layout/panel/Panel";
import { resolvePath, ROUTE_ALL_SCHEDULES, ROUTE_SCHEDULE_ALL_DETAIL } from "_core/router/routes";

type ScheduleDetailUrlParams = { schedulerName: string; scheduleId: string };

const ScheduleDetail = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const handleClose = () => {
    history.goBack();
  };

  const { schedulerName, scheduleId } = useParams<ScheduleDetailUrlParams>();

  return (
    <>
      <Breadcrumb
        data={
          [
            { url: ROUTE_ALL_SCHEDULES,  label: t("Menu-schedules-all") },
            {
              url: resolvePath(ROUTE_SCHEDULE_ALL_DETAIL, {
                schedulerName: schedulerName,
                scheduleId:scheduleId
              }),
              label: scheduleId,
            },
          ]
        }
      />
    <Panel
      icon={"calendar-alt"}
      title={t("Page-title-schedule-detail", { id: scheduleId })}
    >
      <ScheduleForm
        schedulerName={schedulerName}
        scheduleId={scheduleId}
        onClose={handleClose}
      />
    </Panel>
    </>
  );
};

export default ScheduleDetail;

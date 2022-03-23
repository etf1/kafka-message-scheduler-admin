import ScheduleForm from "business/scheduler/component/ScheduleForm";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import Breadcrumb from "_common/component/breadcrumb/Breadcrumb";
import Panel from "_common/component/layout/panel/Panel";
import {
  resolvePath,
  ROUTE_HISTORY_SCHEDULES,
  ROUTE_SCHEDULE_HISTORY_DETAIL,
} from "_core/router/routes";

type ScheduleDetailHistoryUrlParams = {
  schedulerName: string;
  scheduleId: string;
};

const ScheduleDetailHistory = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const handleClose = () => {
    history.goBack();
  };

  const { schedulerName, scheduleId } =
    useParams<ScheduleDetailHistoryUrlParams>();

  return (
    <>
      <Breadcrumb
        data={[
          {
            linkTo: ROUTE_HISTORY_SCHEDULES,
            label: t("Menu-schedules-history"),
          },
          {
            linkTo: resolvePath(ROUTE_SCHEDULE_HISTORY_DETAIL, {
              schedulerName: schedulerName,
              scheduleId: scheduleId,
            }),
            label: scheduleId,
          },
        ]}
      />

      <Panel
        icon={"history"}
        title={t("Page-title-schedule-detail", { id: scheduleId })}
      >
        <ScheduleForm
          schedulerName={schedulerName}
          scheduleId={scheduleId}
          onClose={handleClose}
          scheduleType="history"
        />
      </Panel>
    </>
  );
};

export default ScheduleDetailHistory;

import ScheduleForm from "business/scheduler/component/ScheduleForm";
import { useTranslation } from "react-i18next";
import { useHistory, useParams } from "react-router-dom";
import Breadcrumb from "_common/component/breadcrumb/Breadcrumb";
import Panel from "_common/component/layout/panel/Panel";
import {
  resolvePath,
  ROUTE_LIVE_SCHEDULES,
  ROUTE_SCHEDULE_LIVE_DETAIL,
} from "_core/router/routes";

type ScheduleDetailLiveUrlParams = {
  schedulerName: string;
  scheduleId: string;
};

const ScheduleDetailLive = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const handleClose = () => {
    history.goBack();
  };

  const { schedulerName, scheduleId } =
    useParams<ScheduleDetailLiveUrlParams>();

  return (
    <>
      <Breadcrumb
        data={[
          { linkTo: ROUTE_LIVE_SCHEDULES, label: t("Menu-schedules-live") },
          {
            linkTo: resolvePath(ROUTE_SCHEDULE_LIVE_DETAIL, {
              schedulerName: schedulerName,
              scheduleId: scheduleId,
            }),
            label: scheduleId,
          },
        ]}
      />

      <Panel
        icon={"bolt"}
        title={t("Page-title-schedule-detail", { id: scheduleId })}
      >
        <ScheduleForm
          schedulerName={schedulerName}
          scheduleId={scheduleId}
          onClose={handleClose}
          scheduleType="live"
        />
      </Panel>
    </>
  );
};

export default ScheduleDetailLive;

import ScheduleForm from "business/scheduler/component/ScheduleForm";
import { useTranslation } from "react-i18next";
import { useHistory, useLocation, useParams } from "react-router-dom";
import { ROUTE_ALL_SCHEDULES, ROUTE_LIVE_SCHEDULES } from "_core/router/routes";

type ScheduleDetailUrlParams = { schedulerName: string; scheduleId: string };

const ScheduleDetail = () => {
  const { t } = useTranslation();

  const location = useLocation();
  const history = useHistory();
  const handleClose = () => {
    if (location.pathname.indexOf("live/detail") > -1) {
      history.push(ROUTE_LIVE_SCHEDULES);
    } else {
      history.push(ROUTE_ALL_SCHEDULES);
    }
  };

  const { schedulerName, scheduleId } = useParams<ScheduleDetailUrlParams>();

  return (
    <div className="container has-text-centered">
      <div className="column is-8 is-offset-2">
        <h1 className="title">
          {t("Page-title-schedule-detail", { id: scheduleId })}
        </h1>

        <ScheduleForm
          schedulerName={schedulerName}
          scheduleId={scheduleId}
          onClose={handleClose}
        />
      </div>
    </div>
  );
};

export default ScheduleDetail;

import ScheduleForm from "business/scheduler/component/ScheduleForm";
import { useTranslation } from "react-i18next";
import { useHistory, useLocation, useParams } from "react-router-dom";

type ScheduleDetailUrlParams = { schedulerName: string; scheduleId: string };

const ScheduleDetail = () => {
  const { t } = useTranslation();
  const history = useHistory();
  const handleClose = () => {
    history.goBack();
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

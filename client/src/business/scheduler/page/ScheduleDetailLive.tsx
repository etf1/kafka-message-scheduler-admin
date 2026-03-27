import Breadcrumb from '_common/component/breadcrumb/Breadcrumb';
import Panel from '_common/component/layout/panel/Panel';
import { ROUTE_LIVE_SCHEDULES, ROUTE_SCHEDULE_LIVE_DETAIL, resolvePath } from '_core/router/routes';
import ScheduleForm from 'business/scheduler/component/ScheduleForm';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';

type ScheduleDetailLiveUrlParams = {
  schedulerName: string;
  scheduleId: string;
};

const ScheduleDetailLive = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const handleClose = () => {
    navigate(-1);
  };

  const { schedulerName, scheduleId } = useParams<ScheduleDetailLiveUrlParams>();

  return (
    <>
      <Breadcrumb
        data={[
          { linkTo: ROUTE_LIVE_SCHEDULES, label: t('Menu-schedules-live') },
          {
            linkTo: resolvePath(ROUTE_SCHEDULE_LIVE_DETAIL, {
              schedulerName: schedulerName,
              scheduleId: scheduleId,
            }),
            label: scheduleId ?? '',
          },
        ]}
      />

      <Panel
        icon={'bolt'}
        title={t('Page-title-schedule-detail', { id: scheduleId })}
      >
        <ScheduleForm
          schedulerName={schedulerName!}
          scheduleId={scheduleId!}
          onClose={handleClose}
          scheduleType='live'
        />
      </Panel>
    </>
  );
};

export default ScheduleDetailLive;

import Breadcrumb from '_common/component/breadcrumb/Breadcrumb';
import Panel from '_common/component/layout/panel/Panel';
import { ROUTE_ALL_SCHEDULES, ROUTE_SCHEDULE_ALL_DETAIL, resolvePath } from '_core/router/routes';
import ScheduleForm from 'business/scheduler/component/ScheduleForm';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router-dom';

type ScheduleDetailUrlParams = { schedulerName: string; scheduleId: string };

const ScheduleDetail = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const handleClose = () => {
    navigate(-1);
  };

  const { schedulerName, scheduleId } = useParams<ScheduleDetailUrlParams>();

  return (
    <>
      <Breadcrumb
        data={[
          { linkTo: ROUTE_ALL_SCHEDULES, label: t('Menu-schedules-all') },
          {
            linkTo: resolvePath(ROUTE_SCHEDULE_ALL_DETAIL, {
              schedulerName: schedulerName,
              scheduleId: scheduleId,
            }),
            label: scheduleId ?? '',
          },
        ]}
      />
      <Panel
        icon={'calendar-alt'}
        title={t('Page-title-schedule-detail', { id: scheduleId })}
      >
        <ScheduleForm
          schedulerName={schedulerName!}
          scheduleId={scheduleId!}
          onClose={handleClose}
          scheduleType='all'
        />
      </Panel>
    </>
  );
};

export default ScheduleDetail;

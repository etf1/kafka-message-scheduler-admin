import { useTranslation } from "react-i18next";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";

import { useEffect, useState } from "react";
import { Schedule } from "../type";
import { getScheduleDetail } from "../service/SchedulerService";

export type ScheduleFormProps = {
  schedulerName: string;
  scheduleId: string;
};

const ScheduleForm: React.FC<ScheduleFormProps> = ({ schedulerName, scheduleId }) => {
    const { t } = useTranslation();
  const [schedule, setSchedule] = useState<Schedule>();  
  useEffect(() => {

    if (schedulerName && scheduleId) {
        getScheduleDetail(schedulerName, scheduleId).then((result) => {
            setSchedule(result);
      });
    }
  }, [schedulerName, scheduleId]);
 
  return (
    <div className="box">
  
   { schedule &&
     <fieldset disabled style={{textAlign:"left", margin:"2rem"}}>
        <div className="field">
          <label className="label">{t('Schedule-field-id')}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={schedule.id} />
          </div>
        </div>
        <div className="field">
          <label className="label">{t("Schedule-field-scheduler")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={schedule.scheduler} />
          </div>
        </div>
        <div className="field">
          <label className="label">{t("Schedule-field-creation-date")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={format(fromUnixTime(schedule.timestamp), t("Calendar-date-format"))} />
          </div>
        </div>
        <div className="field">
          <label className="label">{t("Schedule-field-trigger-date")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={format(fromUnixTime(schedule.epoch), t("Calendar-date-format"))} />
          </div>
        </div>
        <div className="field">
          <label className="label">{t("Schedule-field-target-topic")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={schedule.targetTopic} />
          </div>
        </div>
        <div className="field">
          <label className="label">{t("Schedule-field-target-id")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={schedule.targetId} />
          </div>
        </div>
        <div className="field more-space-bottom">
          <label className="label">{t("Schedule-field-headers")}</label>
          <div className="control">
            <input className="input" type="tex"  defaultValue={schedule.headers.map ( h => `${h.name}=${h.value}`).join(', ')} />
          </div>
        </div>
        <div className="field is-grouped " style={{justifyContent:"center"}}>
          <div className="control">
            <button className="button is-link">{t("Close-button")}</button>
          </div>
        </div>
      </fieldset> }   
    </div>
  );
};

export default ScheduleForm;

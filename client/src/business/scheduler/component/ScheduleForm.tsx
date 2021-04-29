import { useTranslation } from "react-i18next";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";

import { useEffect, useState } from "react";
import { Schedule } from "../type";
import { getLiveScheduleDetail, getScheduleDetail } from "../service/SchedulerService";
import IconLabel from "_common/component/element/icon-label/IconLabel";

const formatUnixTime = (time: number, fmt: string) => {
  if (time) {
    const dt = fromUnixTime(time);
    return format(dt, fmt);
  }
  return "";
};

const getScheduleValue = (value:string) => {
  try {
    return atob(value);
  } catch(err) {
    console.error(err);
  }
  return value;
}

export type ScheduleFormProps = {
  schedulerName: string;
  scheduleId: string;
  onClose: () => void;
  live?:boolean;
};

const ScheduleForm: React.FC<ScheduleFormProps> = ({ schedulerName, scheduleId, onClose, live }) => {
  const { t } = useTranslation();
  const [schedule, setSchedule] = useState<Schedule>();
  useEffect(() => {
    if (schedulerName && scheduleId) {
      live ? getLiveScheduleDetail(schedulerName, scheduleId).then((result) => {
        setSchedule(result);
      }): getScheduleDetail(schedulerName, scheduleId).then((result) => {
        setSchedule(result);
      });
    }
  }, [schedulerName, scheduleId]);

  return (
    <div className="box" style={{ padding: "3rem" }}>
      {schedule && (
        <div className="columns">
          <div className="column">
            <fieldset disabled style={{ textAlign: "left" }}>
              <div className="field">
                <label className="label">{t("Schedule-field-id")}</label>
                <div className="control">
                  <input className="input" type="text" defaultValue={schedule.id} />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-scheduler")}</label>
                <div className="control">
                  <input className="input" type="text" defaultValue={schedule.scheduler} />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-creation-date")}</label>
                <div className="control">
                  <input
                    className="input"
                    type="tex"
                    defaultValue={formatUnixTime(schedule.timestamp, t("Calendar-date-hour-format"))}
                  />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-trigger-date")}</label>
                <div className="control">
                  <input
                    className="input"
                    type="tex"
                    defaultValue={formatUnixTime(schedule.epoch, t("Calendar-date-hour-format"))}
                  />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-source-topic")}</label>
                <div className="control">
                  <input className="input" type="text" defaultValue={schedule.topic} />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-target-topic")}</label>
                <div className="control">
                  <input className="input" type="text" defaultValue={schedule.targetTopic} />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Schedule-field-target-id")}</label>
                <div className="control">
                  <input className="input" type="text" defaultValue={schedule.targetId} />
                </div>
              </div>
            </fieldset>
          </div>
          <div className="column">
          <fieldset disabled style={{ textAlign: "left" }}>
            <div className="field">
              <label className="label">{t("Schedule-field-value")}</label>
              <div className="control">
                <textarea rows={12} className="textarea" defaultValue={getScheduleValue(schedule.value)} />
              </div>
            </div>
            </fieldset>
          </div>
        </div>
      )}
      <div className="field is-grouped " style={{ justifyContent: "center", marginTop: "2rem" }}>
        <div className="control">
          <button className="button is-link" onClick={onClose}>
            <IconLabel icon="times" label={t("Close-button")} />
          </button>
        </div>
      </div>
    </div>
  );
};

export default ScheduleForm;

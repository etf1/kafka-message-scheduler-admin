import { useTranslation } from "react-i18next";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";

import { useEffect, useState } from "react";
import { Schedule } from "../type";
import { getScheduleDetail } from "../service/SchedulerService";
import IconLabel from "_common/component/element/IconLabel";

export type ScheduleFormProps = {
  schedulerName: string;
  scheduleId: string;
  onClose: ()=>void;
};

const ScheduleForm: React.FC<ScheduleFormProps> = ({ schedulerName, scheduleId, onClose }) => {
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
    <div className="box" style={{ padding: "3rem" }}>
      {schedule && (
        <fieldset disabled style={{ textAlign: "left" }}>
          <div className="field">
            <label className="label">{t("Schedule-field-id")}</label>
            <div className="control">
              <input className="input" type="tex" defaultValue={schedule.id} />
            </div>
          </div>
          <div className="field">
            <label className="label">{t("Schedule-field-scheduler")}</label>
            <div className="control">
              <input className="input" type="tex" defaultValue={schedule.scheduler} />
            </div>
          </div>
          <div className="field">
            <label className="label">{t("Schedule-field-creation-date")}</label>
            <div className="control">
              <input
                className="input"
                type="tex"
                defaultValue={format(fromUnixTime(schedule.timestamp), t("Calendar-date-format"))}
              />
            </div>
          </div>
          <div className="field">
            <label className="label">{t("Schedule-field-trigger-date")}</label>
            <div className="control">
              <input
                className="input"
                type="tex"
                defaultValue={format(fromUnixTime(schedule.epoch), t("Calendar-date-format"))}
              />
            </div>
          </div>
          <div className="field">
            <label className="label">{t("Schedule-field-target-topic")}</label>
            <div className="control">
              <input className="input" type="tex" defaultValue={schedule.targetTopic} />
            </div>
          </div>
          <div className="field">
            <label className="label">{t("Schedule-field-target-id")}</label>
            <div className="control">
              <input className="input" type="tex" defaultValue={schedule.targetId} />
            </div>
          </div>
          <div className="field more-space-bottom">
            <label className="label">{t("Schedule-field-headers")}</label>
            <div className="control">
              <input
                className="input"
                type="tex"
                defaultValue={schedule.headers.map((h) => `${h.name}=${h.value}`).join(", ")}
              />
            </div>
          </div>
        </fieldset>
      )}
      <div className="field is-grouped " style={{ justifyContent: "center" }}>
        <div className="control">
          <button className="button is-link" onClick={onClose}><IconLabel icon="times" label={t("Close-button")} /></button>
        </div>
      </div>
    </div>
  );
};

export default ScheduleForm;

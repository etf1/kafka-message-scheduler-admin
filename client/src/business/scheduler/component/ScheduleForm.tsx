import { useTranslation } from "react-i18next";

import { useEffect, useState } from "react";
import { Schedule } from "../type";
import { getLiveScheduleDetail, getScheduleDetail } from "../service/SchedulerService";
import Container from "_common/component/layout/container/Container";
import ScheduleVersionTable from "./ScheduleVersionTable";
import useMedia from "_common/hook/useMedia";
import { pluralizeIf } from "_core/i18n";
import Icon from "_common/component/element/icon/Icon";


export type ScheduleFormProps = {
  schedulerName: string;
  scheduleId: string;
  onClose: () => void;
  live?: boolean;
};

const ScheduleForm: React.FC<ScheduleFormProps> = ({ schedulerName, scheduleId, onClose, live }) => {
  const { t } = useTranslation();
  const [schedule, setSchedule] = useState<Schedule[]>();
  const smallScreen = useMedia(["(max-width: 1250px)", "(min-width: 1250px)"], [true, false], true);

  useEffect(() => {
    if (schedulerName && scheduleId) {
      live
        ? getLiveScheduleDetail(schedulerName, scheduleId).then((result) => {
            setSchedule(result);
          })
        : getScheduleDetail(schedulerName, scheduleId).then((result) => {
            setSchedule(result);
          });
    }
  }, [schedulerName, scheduleId, live]);

  const firstSchedule = schedule && schedule[0]

  return (
    <>
      <Container title={<><Icon name="cog"/> {t("Schedule-field-main")}</>}>
        <div className="box" style={{ padding: "3rem" }}>
          {firstSchedule && (
            <div className="columns is-desktop">
              <div className="column is-6">
                <fieldset disabled style={{ textAlign: "left" }}>
                  <div className="field">
                    <label className="label">{t("Schedule-field-id")}</label>
                    <div className="control">
                      <input className="input" type="text" defaultValue={firstSchedule.id} />
                    </div>
                  </div>
                  <div className="field">
                    <label className="label">{t("Schedule-field-scheduler")}</label>
                    <div className="control">
                      <input className="input" type="text" defaultValue={firstSchedule.scheduler} />
                    </div>
                  </div>
                </fieldset>
              </div>
              <div className="column is-6">
              <fieldset disabled style={{ textAlign: "left" }}>
                <div className="field">
                  <label className="label">{t("Schedule-field-source-topic")}</label>
                  <div className="control">
                    <input className="input" type="text" defaultValue={firstSchedule.topic} />
                  </div>
                </div>
                </fieldset>
              </div>
            </div>
          )}
        </div>
      </Container>
      <Container title={<><Icon name="copy"/> {(schedule?.length || 0)+" "+pluralizeIf((schedule?.length || 0), t("Version"), t("Versions")) || ""}</>}>
      <div className="box" style={{ padding: "3rem" }}>
            <ScheduleVersionTable data={schedule || []} showAsTable={!smallScreen}/>
      </div>
      </Container>
    </>
  );
};

export default ScheduleForm;

import { useTranslation } from "react-i18next";

import { useEffect, useState } from "react";
import { Schedule, ScheduleType } from "../type";
import {
  getScheduleDetailByType,
} from "../service/SchedulerService";
import Container from "_common/component/layout/container/Container";
import ScheduleVersionTable from "./ScheduleVersionTable";
import useMedia from "_common/hook/useMedia";
import { pluralizeIf } from "_core/i18n";
import Icon from "_common/component/element/icon/Icon";
import Appear from "_common/component/transition/Appear";

export type ScheduleFormProps = {
  schedulerName: string;
  scheduleId: string;
  onClose: () => void;
  scheduleType: ScheduleType;
};

const ScheduleForm: React.FC<ScheduleFormProps> = ({ schedulerName, scheduleId, onClose, scheduleType }) => {
  const { t } = useTranslation();
  const [schedule, setSchedule] = useState<Schedule[]>();
  const smallScreen = useMedia(["(max-width: 1250px)", "(min-width: 1250px)"], [true, false], true);
  const [error, setError] = useState<Error>();

  useEffect(() => {
    if (schedulerName && scheduleId) {
      getScheduleDetailByType(scheduleType)(schedulerName, scheduleId)
        .then((result) => {
          setSchedule(result);
          setError(undefined);
        })
        .catch((err: Error) => {
          console.error(err);
          setError(err);
        });
    }
  }, [schedulerName, scheduleId, scheduleType]);

  const firstSchedule = schedule && schedule[0];

  return (
    <Appear visible={!!firstSchedule}>
      {(nodeRef) => (
        <div ref={nodeRef}>
          <Container
            title={
              <>
                <Icon name="cog" /> {t("Schedule-field-main")}
              </>
            }
          >
            <div style={{ padding: "2rem" }}>
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
          <Container
            title={
              <>
                <Icon name="copy" />{" "}
                {(schedule?.length || 0) + " " + pluralizeIf(schedule?.length || 0, t("Version"), t("Versions")) || ""}
              </>
            }
          >
            <div style={{ padding: "2rem" }}>
              {error && (
                <div className="animate-opacity" style={{ fontWeight: 800, color: "red" }}>
                  <Icon name="exclamation-triangle" /> {t("LoadingError")}
                </div>
              )}
              {!error && <ScheduleVersionTable data={schedule || []} showAsTable={!smallScreen} />}
            </div>
          </Container>
        </div>
      )}
    </Appear>
  );
};

export default ScheduleForm;

import SchedulerInstanceTable from "business/scheduler/component/SchedulerInstanceTable";
import useSchedulers from "business/scheduler/hook/useSchedulers";
import { SchedulerInstance } from "business/scheduler/type";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import Container from "_common/component/layout/container/Container";
import Panel from "_common/component/layout/panel/Panel";

const SchedulerDetail = () => {
  const { t } = useTranslation();
  const { schedulerName } = useParams<{ schedulerName: string }>();
  const { schedulers } = useSchedulers();

  const scheduler = schedulers.find((sch) => sch.name === schedulerName);

  const instances: SchedulerInstance[] = scheduler?.instances || [];

  return (
    <Panel icon={"stopwatch"} title={t("Page-title-scheduler-detail")}>
      <div className="box" style={{ padding: "3rem" }}>
        {scheduler && (
          <div className="columns">
            <div className="column">
              <fieldset disabled style={{ textAlign: "left" }}>
                <div className="field">
                  <label className="label">{t("Scheduler-field-name")}</label>
                  <div className="control">
                    <input className="input" type="text" defaultValue={scheduler.name} />
                  </div>
                </div>
                <div className="field">
                  <label className="label">{t("Scheduler-field-port")}</label>
                  <div className="control">
                    <input className="input" type="text" defaultValue={scheduler.http_port} />
                  </div>
                </div>
              </fieldset>
            </div>
          </div>
        )}
      </div>
      <Container title={t("Scheduler-field-instances")}>
        <SchedulerInstanceTable schedulerInstances={instances} />
      </Container>
    </Panel>
  );
};

export default SchedulerDetail;

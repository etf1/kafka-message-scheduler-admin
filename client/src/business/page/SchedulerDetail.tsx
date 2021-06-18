import SchedulerInstanceTable from "business/scheduler/component/SchedulerInstanceTable";
import useSchedulers from "business/scheduler/hook/useSchedulers";
import { SchedulerInstance } from "business/scheduler/type";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import Breadcrumb from "_common/component/breadcrumb/Breadcrumb";
import Container from "_common/component/layout/container/Container";
import Panel from "_common/component/layout/panel/Panel";
import { pluralizeIf } from "_core/i18n";
import { resolvePath, ROUTE_SCHEDULER_DETAIL } from "_core/router/routes";

const SchedulerDetail = () => {
  const { t } = useTranslation();
  const { schedulerName } = useParams<{ schedulerName: string }>();
  const { schedulers } = useSchedulers();

  const scheduler = schedulers.find((sch) => sch.name === schedulerName);

  const instances: SchedulerInstance[] = scheduler?.instances || [];

  return (
    <>
      <Breadcrumb
        data={
          scheduler
            ? [
                { url: "/schedulers", label: t("Menu-schedulers") },
                {
                  url: resolvePath(ROUTE_SCHEDULER_DETAIL, {
                    schedulerName: scheduler.name,
                  }),
                  label: scheduler.name,
                },
              ]
            : []
        }
      />

      <Panel icon={"stopwatch"} title={t("Page-title-scheduler-detail")}>
        <Container title={t("Scheduler-field-main")}>
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
        </Container>
        <Container title={instances.length+" "+pluralizeIf(instances.length, t("Scheduler-field-instance"), t("Scheduler-field-instances")) || ""}>
          <div className="box" style={{ padding: "3rem" }}>
            <SchedulerInstanceTable schedulerInstances={instances} />
          </div>
        </Container>
      </Panel>
    </>
  );
};

export default SchedulerDetail;

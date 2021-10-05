import SchedulerInstanceTable from "business/scheduler/component/SchedulerInstanceTable";
import useSchedulers from "business/scheduler/hook/useSchedulers";
import { SchedulerInstance } from "business/scheduler/type";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import Breadcrumb from "_common/component/breadcrumb/Breadcrumb";
import Icon from "_common/component/element/icon/Icon";
import Container from "_common/component/layout/container/Container";
import Panel from "_common/component/layout/panel/Panel";
import Appear from "_common/component/transition/Appear";
import { pluralizeIf } from "_core/i18n";
import { resolvePath, ROUTE_SCHEDULERS, ROUTE_SCHEDULER_DETAIL } from "_core/router/routes";

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
                { linkTo: ROUTE_SCHEDULERS, label: t("Menu-schedulers") },
                {
                  linkTo: resolvePath(ROUTE_SCHEDULER_DETAIL, {
                    schedulerName: scheduler.name,
                  }),
                  label: scheduler.name,
                },
              ]
            : []
        }
      />

      <Panel icon={"stopwatch"} title={t("Page-title-scheduler-detail")}>
        <Appear visible={!!scheduler}>
          {(nodeRef) => (
            <div ref={nodeRef}>
              <Container
                title={
                  <>
                    <Icon name="cog" /> {t("Scheduler-field-main")}
                  </>
                }
              >
                <div>
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
            </div>
          )}
        </Appear>
         <hr style={{marginLeft:-20, width:"133%"}}/>
        
        <Appear visible={instances && instances.length > 0}>
          {(nodeRef) => (
            <div ref={nodeRef}>
              <Container
                title={
                  <>
                    <Icon name="copy" />
                    {instances.length +
                      " " +
                      pluralizeIf(instances.length, t("Scheduler-field-instance"), t("Scheduler-field-instances")) ||
                      ""}
                  </>
                }
              >
                <div>
                  <SchedulerInstanceTable schedulerInstances={instances} />
                </div>
              </Container>
            </div>
          )}
        </Appear>
      </Panel>
    </>
  );
};

export default SchedulerDetail;

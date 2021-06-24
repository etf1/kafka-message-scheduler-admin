import SchedulerTable from "business/scheduler/component/SchedulerTable";
import useSchedulers from "business/scheduler/hook/useSchedulers";
import React from "react";
import { useTranslation } from "react-i18next";
import Icon from "_common/component/element/icon/Icon";
import Panel from "_common/component/layout/panel/Panel";
import { ROUTE_SCHEDULER_DETAIL } from "_core/router/routes";

const Schedulers = () => {
  const { t } = useTranslation();
  const { schedulers, error } = useSchedulers();

  return (
    <Panel icon={"stopwatch"} title={t("Page-title-schedulers")}>
      {error && <div className="animate-opacity" style={{ fontWeight: 800, color: "red" }}><Icon name="exclamation-triangle"/> {t("LoadingError")}</div>}
      {!error && schedulers && <SchedulerTable schedulers={schedulers} detailUrl={ROUTE_SCHEDULER_DETAIL} />}
    </Panel>
  );
};

export default Schedulers;

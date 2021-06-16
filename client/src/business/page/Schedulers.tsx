import React from "react";
import { useTranslation } from "react-i18next";
import Panel from "_common/component/layout/panel/Panel";

const Schedulers = () => {

  const { t } = useTranslation();

  return <Panel icon={"stopwatch"} title={t("Schedulers")}>

  </Panel>
}


export default Schedulers;
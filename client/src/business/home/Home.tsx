import { useTranslation } from "react-i18next";

import {
  AppStat,
  getAppStats,
} from "business/scheduler/service/SchedulerService";
import { useEffect, useState } from "react";
import { clear } from "_common/service/SessionStorageService";
import AppStatCard from "./AppStartCard";
import Icon from "_common/component/element/icon/Icon";
import Loader from "_common/component/element/Loader";

import Style from "./Home.module.css";
import clsx from "clsx";

const Home = () => {
  const { t } = useTranslation();
  const [error, setError] = useState<Error>();
  const [stats, setStats] = useState<AppStat[]>();
  clear(() => false);
  useEffect(() => {
    (async () => {
      try {
        const stats = await getAppStats();

        setStats(stats);
        setError(undefined);
      } catch (err) {
        console.error(err);
        setError(err);
      }
    })();
  }, []);

  return (
    <div className={clsx("columns", Style.Home)}>
      {error && (
        <div className={clsx("animate-opacity", Style.LoadingError)}>
          <Icon name="exclamation-triangle" /> {t("LoadingError")}
        </div>
      )}
      {stats === undefined && <Loader />}
      {!error &&
        stats &&
        stats.map((st) => {
          return <AppStatCard key={st.scheduler} stat={st} />;
        })}
    </div>
  );
};

export default Home;

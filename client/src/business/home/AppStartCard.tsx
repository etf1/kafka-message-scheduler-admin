import { AppStat } from "business/scheduler/service/SchedulerService";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import Icon from "_common/component/element/icon/Icon";
import { ROUTE_ALL_SCHEDULES, ROUTE_LIVE_SCHEDULES } from "_core/router/routes";

type AppStatCardProps = {
  stat: AppStat;
};

const AppStatCard: React.FC<AppStatCardProps> = ({ stat }) => {
  const { t } = useTranslation();
  return (
    <div className="column is-4-tablet is-4-desktop">
      <div className="card">
        <div className="card-header" style={{ backgroundColor: "orange" }}>
          <h3 className="card-header-title is-inline" style={{ color: "white" }}>
              <Icon name={"stopwatch"} style={{marginRight:20}}  className={"has-tooltip-right"} />{" "}
               <span>{stat.scheduler}</span>
          </h3>
        </div>
        <div
          className="card-content"
          style={{
            maxHeight: 450,
            backgroundColor: "#f5f5f5",
            paddingLeft: 0,
            paddingRight: 0,
            paddingTop: "1rem",
          }}
        >
          <h3
            className="subtitle is-6"
            style={{
              margin: 0,
              padding: "1rem",
            }}
          >
            <Link to={ROUTE_LIVE_SCHEDULES + "?schedulerName=" + stat.scheduler}>
              <Icon name={"calendar"} style={{marginRight:20}}  className={"has-tooltip-right"} data-tooltip={t("SchedulesLive")} />{" "}
              {stat.total_live} {t("SchedulesLive")}
            </Link>
          </h3>
          <h3
            className="subtitle is-6"
            style={{
              margin: 0,
              padding: "1rem",
            }}
          >
            <Link to={ROUTE_ALL_SCHEDULES + "?schedulerName=" + stat.scheduler}>
              <Icon name={"calendar-alt"}  style={{marginRight:20}} className={"has-tooltip-right"} data-tooltip={t("Schedules")} />{" "}
              {stat.total} {t("Schedules")}
            </Link>
          </h3>
        </div>
      </div>
    </div>
  );
};

export default AppStatCard;

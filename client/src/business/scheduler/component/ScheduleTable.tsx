import { useTranslation } from "react-i18next";
import { ScheduleInfo } from "../type";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";
import React from "react";
import { Link } from "react-router-dom";
import Styles from "./ScheduleTable.module.css";
import clsx from "clsx";
import { resolvePath, ROUTE_SCHEDULE_DETAIL } from "_core/router/routes";
export type ScheduleTableProps = {
  data: ScheduleInfo[];
  onClick?: (schedule: ScheduleInfo) => void;
};

const ScheduleTable: React.FC<ScheduleTableProps> = ({ data, onClick }) => {
  const { t } = useTranslation();
  return (
    <table className="table is-striped is-narrow is-hoverable is-fullwidth">
      <thead>
        <tr>
          <th>{t("ScheduleTable-column-ID")}</th>
          <th>{t("ScheduleTable-column-Scheduler")}</th>
          <th>{t("ScheduleTable-column-CreationDate")}</th>
          <th>{t("ScheduleTable-column-TiggerDate")}</th>
          <th>{t("ScheduleTable-column-TargetTopic")}</th>
          <th>{t("ScheduleTable-column-TargetId")}</th>
        </tr>
      </thead>

      <tbody>
        {data.map((schedule) => {
          return (
            <tr key={`${schedule.scheduler}/${schedule.id}`} onClick={() => onClick && onClick(schedule)}>
              <td className={clsx(Styles.ColWithId, Styles.ColWithLink)}>
                <Link
                  to={resolvePath(ROUTE_SCHEDULE_DETAIL, {
                    schedulerName: schedule.scheduler,
                    scheduleId: schedule.id,
                  })}
                >
                  {schedule.id}
                </Link>
              </td>
              <td className={Styles.colWithId}>{schedule.scheduler}</td>
              <td>{format(fromUnixTime(schedule.timestamp), t("Calendar-date-format"))}</td>
              <td>{format(fromUnixTime(schedule.epoch), t("Calendar-date-format"))}</td>
              <td className={Styles.colWithId}>{schedule.targetTopic}</td>
              <td className={Styles.colWithId}>{schedule.targetId}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

export default ScheduleTable;

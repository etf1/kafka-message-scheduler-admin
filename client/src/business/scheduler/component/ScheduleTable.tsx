import { useTranslation } from "react-i18next";
import { ScheduleInfo } from "../type";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";
import React from "react";
import { Link } from "react-router-dom";
import Styles from "./ScheduleTable.module.css";
import clsx from "clsx";
import { resolvePath } from "_core/router/routes";

export type ScheduleTableProps = {
  data: ScheduleInfo[];
  onClick?: (schedule: ScheduleInfo) => void;
  detailUrl: string;
  showAsTable?: boolean;
};

const ScheduleTable: React.FC<ScheduleTableProps> = ({
  data,
  detailUrl,
  onClick,
  showAsTable,
}) => {
  const { t } = useTranslation();

  return showAsTable || showAsTable === undefined ? (
    <table
      key="table"
      className="table is-striped is-narrow is-hoverable is-fullwidth"
    >
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
            <tr
              key={`${schedule.scheduler}/${schedule.id}`}
              onClick={() => onClick && onClick(schedule)}
            >
              <td className={clsx(Styles.ColWithId, Styles.ColWithLink)}>
                <Link
                  to={resolvePath(detailUrl, {
                    schedulerName: schedule.scheduler,
                    scheduleId: schedule.id,
                  })}
                >
                  {schedule.id}
                </Link>
              </td>
              <td className={Styles.colWithId}>{schedule.scheduler}</td>
              <td>
                {format(
                  fromUnixTime(schedule.timestamp),
                  t("Calendar-date-format")
                )}
              </td>
              <td>
                {format(
                  fromUnixTime(schedule.epoch),
                  t("Calendar-date-format")
                )}
              </td>
              <td className={Styles.colWithId}>{schedule.targetTopic}</td>
              <td className={Styles.colWithId}>{schedule.targetId}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  ) : (
    <div>
      {data.map((schedule) => {
        return (
          <fieldset
            className="box "
            key={`cards${schedule.scheduler}/${schedule.id}`}
            disabled
            style={{ textAlign: "left", marginBottom: 20 }}
          >
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-id")}</strong>
              <Link
                to={resolvePath(detailUrl, {
                  schedulerName: schedule.scheduler,
                  scheduleId: schedule.id,
                })}
              >
                <span className={clsx(Styles.ValueField, Styles.ColWithLink)}>
                  {schedule.id}
                </span>
              </Link>
            </div>
            <div className="space-right">
              <strong className="space-right">
                {t("Schedule-field-scheduler")}
              </strong>
              <span className={Styles.ValueField}>{schedule.scheduler}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">
                {t("Schedule-field-creation-date")}
              </strong>
              <span className={clsx("space-right", Styles.ValueField)}>
                {format(
                  fromUnixTime(schedule.timestamp),
                  t("Calendar-date-format")
                )}
                ,{" "}
              </span>
              <strong className={clsx("space-right", Styles.ValueField)}>
                {t("Schedule-field-trigger-date")}
              </strong>
              <span className={Styles.ValueField}>
                {format(
                  fromUnixTime(schedule.epoch),
                  t("Calendar-date-format")
                )}
              </span>
            </div>

            <div className="space-right">
              <strong className="space-right">
                {t("Schedule-field-target-topic")}
              </strong>
              <span className={Styles.ValueField}>{schedule.targetTopic}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">
                {t("Schedule-field-target-id")}
              </strong>
              <span className={Styles.ValueField}>{schedule.targetId}</span>
            </div>
          </fieldset>
        );
      })}
    </div>
  );
};

export default ScheduleTable;

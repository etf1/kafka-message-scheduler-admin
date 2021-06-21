import { useTranslation } from "react-i18next";
import { ScheduleInfo } from "../type";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Style from "./ScheduleTable.module.css";
import clsx from "clsx";
import { resolvePath } from "_core/router/routes";
import { SortOrder, SortType } from "../service/SchedulerService";
import Icon from "_common/component/element/icon/Icon";

const formatUnixTime = (time: number, fmt: string) => {
  if (time) {
    const dt = fromUnixTime(time);
    return format(dt, fmt);
  }
  return "";
};
export type SortModel = { type: SortType; order: SortOrder } | undefined;
export type ScheduleTableProps = {
  data: ScheduleInfo[];
  onClick?: (schedule: ScheduleInfo) => void;
  onSort: (column: SortType, sortOrder: SortOrder) => void;
  detailUrl: string;
  showAsTable?: boolean;
  ref?: React.LegacyRef<HTMLTableElement> | undefined;
};

const ScheduleTable: React.FC<ScheduleTableProps> = ({ ref, data, detailUrl, onClick, onSort, showAsTable }) => {
  const { t } = useTranslation();
  const [sort, setSort] = useState<SortModel>();

  const handleSort = (column: SortType) => {
    if (sort && column === sort.type) {
      setSort({type:sort.type, order: sort.order==="asc"? "desc":"asc"});
    } else {
      setSort({type:column, order: "asc"});
    }
  };

  useEffect(()=>{

    sort && onSort(sort.type, sort.order);

  }, [sort, onSort]);

  const renderSortIcon = (key:SortType, sort:SortModel)=> {
    if (sort?.type === key) {
        if (sort?.order === "asc") {
          return <Icon name="long-arrow-alt-down" />
        } else {
          return <Icon name="long-arrow-alt-up" />
        }
    }

    return null;
  }

  return showAsTable || showAsTable === undefined ? (
    <table ref={ref} key="table" className="table is-striped is-hoverable is-fullwidth">
      <thead>
        <tr>
          <th style={{cursor:"pointer"}} onClick={() => handleSort("id")}>{t("ScheduleTable-column-ID")} {renderSortIcon("id", sort)}</th>
          <th style={{cursor:"pointer"}} onClick={() => handleSort("timestamp")}>{t("ScheduleTable-column-CreationDate")} {renderSortIcon("timestamp", sort)}</th>
          <th style={{cursor:"pointer"}} onClick={() => handleSort("epoch")}>{t("ScheduleTable-column-TiggerDate")} {renderSortIcon("epoch", sort)}</th>
          <th>{t("ScheduleTable-column-TargetTopic")}</th>
          <th>{t("ScheduleTable-column-TargetId")}</th>
        </tr>
      </thead>

      <tbody>
        {data.map((schedule, index) => {
          return (
            <tr key={`${index} ${schedule.scheduler}/${schedule.id}`} onClick={() => onClick && onClick(schedule)}>
              <td className={clsx(Style.ColWithId, Style.ColWithLink)}>
                <Link
                  to={resolvePath(detailUrl, {
                    schedulerName: schedule.scheduler,
                    scheduleId: schedule.id,
                  })}
                >
                  {schedule.id}
                </Link>
              </td>
              <td>{formatUnixTime(schedule.timestamp, t("Calendar-date-hour-format"))}</td>
              <td>{formatUnixTime(schedule.epoch, t("Calendar-date-hour-format"))}</td>
              <td className={Style.colWithId}>{schedule.targetTopic}</td>
              <td className={Style.colWithId}>{schedule.targetId}</td>
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
                <span className={clsx(Style.ValueField, Style.ColWithLink)}>{schedule.id}</span>
              </Link>
            </div>
          
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-creation-date")}</strong>
              <span className={clsx("space-right", Style.ValueField)}>
                {formatUnixTime(schedule.timestamp, t("Calendar-date-hour-format"))},{" "}
              </span>
              <strong className={clsx("space-right", Style.ValueField)}>{t("Schedule-field-trigger-date")}</strong>
              <span className={Style.ValueField}>
                {formatUnixTime(schedule.epoch, t("Calendar-date-hour-format"))}
              </span>
            </div>

            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-target-topic")}</strong>
              <span className={Style.ValueField}>{schedule.targetTopic}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-target-id")}</strong>
              <span className={Style.ValueField}>{schedule.targetId}</span>
            </div>
          </fieldset>
        );
      })}
    </div>
  );
};

export default ScheduleTable;

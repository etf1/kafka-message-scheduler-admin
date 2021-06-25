import { useTranslation } from "react-i18next";
import { Schedule } from "../type";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";
import React from "react";
import Style from "./ScheduleVersionTable.module.css";
import clsx from "clsx";
import { base64DecToArr, truncate, UTF8ArrToStr } from "_common/service/FunUtil";
import Icon from "_common/component/element/icon/Icon";
import ModalService from "_common/component/modal/ModalService";

const formatUnixTime = (time: number, fmt: string) => {
  if (time) {
    const dt = fromUnixTime(time);
    return format(dt, fmt);
  }
  return "";
};

const getScheduleValue = (value: string) => {
  try {
    
    var sortieUT8 = base64DecToArr(value);

    return UTF8ArrToStr(sortieUT8);
  } catch (err) {
    console.error(err);
  }
  return value;
};
export type ScheduleVersionTableProps = {
  data: Schedule[];
  onClick?: (schedule: Schedule) => void;
  showAsTable?: boolean;
};

const ScheduleVersionTable: React.FC<ScheduleVersionTableProps> = ({ data, onClick, showAsTable }) => {
  const { t } = useTranslation();

  const showValueDetail = (schedule: Schedule) => {
    ModalService.message({ title: t("Schedule-field-target-value"), message: getScheduleValue(schedule.value) });
  };

  return showAsTable || showAsTable === undefined ? (
    <table key="table" className="table is-striped is-hoverable is-fullwidth">
      <thead>
        <tr>
          <th style={{ minWidth: 190 }}>{t("ScheduleVersionTable-column-CreationDate")}</th>
          <th style={{ minWidth: 190 }}>{t("ScheduleVersionTable-column-TiggerDate")}</th>
          <th style={{ minWidth: 180 }}>{t("ScheduleVersionTable-column-TargetTopic")}</th>
          <th style={{ minWidth: 180 }}>{t("ScheduleVersionTable-column-TargetId")}</th>
          <th>{t("ScheduleVersionTable-column-Value")}</th>
        </tr>
      </thead>

      <tbody>
        {data.map((schedule, index) => {
          const value = getScheduleValue(schedule.value);
          return (
            <tr key={`${index} ${schedule.scheduler}/${schedule.id}`} onClick={() => onClick && onClick(schedule)}>
              <td>{formatUnixTime(schedule.timestamp, t("Calendar-date-hour-format"))}</td>
              <td>{formatUnixTime(schedule.epoch, t("Calendar-date-hour-format"))}</td>
              <td className={Style.colWithId}>{schedule.targetTopic}</td>
              <td className={Style.colWithId}>{schedule.targetId}</td>
              {value ? <td onClick={() => showValueDetail(schedule)} className={clsx(Style.colWithId, Style.ColWithLink)}>
                
                {truncate(value, 250)}{" "}
                <span style={{ color: "gray", fontStyle: "italic" }}>
                  ({value && value.length}&nbsp;{t("Chars")})
                </span>{" "}
                <Icon name="eye" />
              </td> : <td />}
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
              <span className={clsx(Style.ValueField, Style.ColWithLink)}>{schedule.id}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-scheduler")}</strong>
              <span className={Style.ValueField}>{schedule.scheduler}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-creation-date")}</strong>
              <span className={clsx("space-right", Style.ValueField)}>
                {formatUnixTime(schedule.timestamp, t("Calendar-date-hour-format"))},{" "}
              </span>
              <strong className={clsx("space-right", Style.ValueField)}>{t("Schedule-field-trigger-date")}</strong>
              <span className={Style.ValueField}>{formatUnixTime(schedule.epoch, t("Calendar-date-hour-format"))}</span>
            </div>

            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-target-topic")}</strong>
              <span className={Style.ValueField}>{schedule.targetTopic}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-target-id")}</strong>
              <span className={Style.ValueField}>{schedule.targetId}</span>
            </div>
            <div className="space-right">
              <strong className="space-right">{t("Schedule-field-target-value")}</strong>
              <span className={Style.ValueField}>{truncate(getScheduleValue(schedule.value), 80)}</span>
            </div>
          </fieldset>
        );
      })}
    </div>
  );
};

export default ScheduleVersionTable;

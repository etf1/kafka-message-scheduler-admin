import { useTranslation } from "react-i18next";
import { ScheduleInfo } from "../type";
import fromUnixTime from "date-fns/fromUnixTime";
import format from "date-fns/format";
import { CSSProperties } from "react";
export type ScheduleTableProps = {
  data: ScheduleInfo[];
  onClick: (schedule: ScheduleInfo) => void;
};
const styles: { colWithId: CSSProperties } = {
  colWithId: {
    textAlign: "left",
    maxWidth: 210,
    wordBreak: "break-all",
  },
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
            <tr key={`${schedule.scheduler}/${schedule.id}`} className="pointer">
              <td style={styles.colWithId}>{schedule.id}</td>
              <td style={styles.colWithId}>{schedule.scheduler}</td>
              <td>{format(fromUnixTime(schedule.timestamp), t("Calendar-date-format"))}</td>
              <td>{format(fromUnixTime(schedule.epoch), t("Calendar-date-format"))}</td>
              <td style={styles.colWithId}>{schedule.targetTopic}</td>
              <td style={styles.colWithId}>{schedule.targetId}</td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

export default ScheduleTable;

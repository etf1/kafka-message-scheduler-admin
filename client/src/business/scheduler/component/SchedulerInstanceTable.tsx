import clsx from "clsx";
import { useTranslation } from "react-i18next";
import { SchedulerInstance } from "../type";
import Style from "./SchedulerInstanceTable.module.css";

export type SchedulerInstanceTableProps = {
    schedulerInstances : SchedulerInstance[];
    onClick?: (scheduler:SchedulerInstance) => void;
}
/*
bootstrap_servers: "localhost:9092"
hostname: ["localhost"]
ip: "127.0.0.1"
topics: ["schedules"]
0: "schedules"
*/
const SchedulerInstanceTable:React.FC<SchedulerInstanceTableProps> = ({schedulerInstances, onClick})=>{
    const { t } = useTranslation();
    return (
        <table key="table" className="table is-striped is-hoverable is-fullwidth">
        <thead>
          <tr>
            <th style={{cursor:"pointer", minWidth:150}}>{t("SchedulerInstanceTable-column-Ip")}</th>
            <th style={{cursor:"pointer", minWidth:150}}>{t("SchedulerInstanceTable-column-Hostname")}</th>
            <th style={{cursor:"pointer"}}>{t("SchedulerInstanceTable-column-BootstrapServers")}</th>
            <th style={{cursor:"pointer"}}>{t("SchedulerInstanceTable-column-Topics")}</th>
     
          </tr>
        </thead>
  
        <tbody>
          {schedulerInstances.map((inst) => {
            return (
              <tr key={`${inst.ip}`} onClick={() => onClick && onClick(inst)}>
                <td className={clsx(Style.ColWithId, Style.ColWithLink)}>
                {inst.ip}
                </td>
                <td> 
                {inst.hostname.join( ', ')}
                </td>
                <td> 
                {inst.bootstrap_servers}
                </td>
                <td> 
                {inst.topics.join( ', ')}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    )
}


export default SchedulerInstanceTable;
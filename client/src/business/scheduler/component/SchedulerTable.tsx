import clsx from "clsx";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { resolvePath } from "_core/router/routes";
import { Scheduler } from "../type";
import Styles from "./SchedulerTable.module.css";

export type SchedulerTableProps = {
    schedulers : Scheduler[];
    onClick?: (scheduler:Scheduler) => void;
    detailUrl: string;
}

const SchedulerTable:React.FC<SchedulerTableProps> = ({schedulers, detailUrl, onClick})=>{
    const { t } = useTranslation();
    return (
        <table key="table" className="table is-striped is-hoverable is-fullwidth">
        <thead>
          <tr>
            <th style={{cursor:"pointer"}}>{t("SchedulerTable-column-Name")}</th>
            <th style={{cursor:"pointer"}}>{t("SchedulerTable-column-Port")}</th>
            <th style={{cursor:"pointer"}}>{t("SchedulerTable-column-Nb-Instances")}</th>
     
          </tr>
        </thead>
  
        <tbody>
          {schedulers.map((scheduler) => {
            return (
              <tr key={`${scheduler.name}`} onClick={() => onClick && onClick(scheduler)}>
                <td className={clsx(Styles.ColWithId, Styles.ColWithLink)}>
                  <Link
                    to={resolvePath(detailUrl, {
                      schedulerName: scheduler.name
                    })}
                  >
                    {scheduler.name}
                  </Link>
                </td>
                <td> 
                {scheduler.http_port}
                </td>
                <td> 
                {scheduler.instances.length}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    )
}


export default SchedulerTable;
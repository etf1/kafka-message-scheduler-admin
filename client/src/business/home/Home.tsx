import { AppStat, getAppStats } from "business/scheduler/service/SchedulerService";
import { useEffect, useState } from "react";
import AppStatCard from "./AppStartCard";

const Home = () => {
  const [stats, setStats] = useState<AppStat[]>([]);

  useEffect( () => {
    (async ()=>{
      console.log("Home");
      const stats = await getAppStats();

      setStats(stats);
    })();
  }, [])

  return <div className="columns" style={{margin:"3rem"}}>
      {stats.map (st => {
        return <AppStatCard key={st.scheduler} stat={st} />
      })}
  </div>;
};

export default Home;

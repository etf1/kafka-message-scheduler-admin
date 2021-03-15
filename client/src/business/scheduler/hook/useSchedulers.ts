import { useEffect, useState } from "react";
import useRefresh from "_common/hook/useRefresh";
import { listAllSchedulers } from "../service/SchedulerService";
import { Scheduler } from "./../type/index";
const useSchedulers = () => {
  const [refresh, count] = useRefresh();
  const [schedulers, setSchedulers] = useState<Scheduler[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  useEffect(() => {
    setIsLoading(true);
    (async () => {
      const res = await listAllSchedulers();
      setSchedulers(res);
      setIsLoading(false);
    })();
  }, [count]);

  return { schedulers, isLoading, refresh };
};

export default useSchedulers;

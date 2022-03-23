import { useEffect, useState } from "react";
import useRefresh from "_common/hook/useRefresh";
import { listAllSchedulers } from "../service/SchedulerService";
import { Scheduler } from "./../type/index";
const useSchedulers = () => {
  const [refresh, count] = useRefresh();
  const [schedulers, setSchedulers] = useState<Scheduler[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error>();
  useEffect(() => {
    setIsLoading(true);
    (async () => {
      try {
        const res = await listAllSchedulers();
        setSchedulers(res);
        setIsLoading(false);
        setError(undefined);
      } catch (err) {
        console.error(err);
        setError(err);
      }
    })();
  }, [count]);

  return { schedulers, isLoading, refresh, error };
};

export default useSchedulers;

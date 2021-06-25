export type Scheduler = {
  name: string;
  instances: SchedulerInstance[];
  http_port: string;
};

export type SchedulerInstance = {
  bootstrap_servers: string;
  hostname: string[];
  ip: string;
  topics: string[];
};

export type Header = { name: string; value: string };

export type ScheduleInfo = {
  id: string;
  scheduler: string;
  timestamp: number; // date de création de la planif
  epoch: number; // date a laquelle va etre déclenché la planif
  targetTopic: string;
  targetId: string;
};

export type Schedule = ScheduleInfo & {
  topic: string;
  value: string;
};


export type ScheduleType = "live" | "all" | "history"
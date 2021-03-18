export type Scheduler = {
  name: string;
  instances: SchedulerInstance[];
};

export type SchedulerInstance = {
  id: string;
  names: string[];
  topics: string[];
  partitions: number[];
  bootstrapServers: string[];
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
  headers: Header[];
  body: string;
};

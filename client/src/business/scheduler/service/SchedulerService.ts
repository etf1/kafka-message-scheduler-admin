import { Schedule, ScheduleInfo, Scheduler } from "../type";

export type SortType = "id" | "epoch" | "timestamp";
export type SortOrder = "asc" | "desc";

export type SearchParams = {
  scheduleId?: string;
  epochFrom?: number;
  epochTo?: number;
  sort?: SortType;
  sortOrder?: SortOrder;
  schedulerName: string;
  max: number; // -1  for all
};

export const makeSearchArgs = (p: SearchParams): string => {
  let res = `?scheduler-name=${p.schedulerName}`;
  if (p.scheduleId) {
    res += `&schedule-id=${p.scheduleId}`;
  }
  if (p.max) {
    res += `&max=${p.max}`;
  }
  if (p.sort) {
    res += `&sort=${p.sort} ${p.sortOrder || "asc"}`;
  }
  if (p.epochFrom) {
    res += `&epoch-from=${p.epochFrom}`;
  }
  if (p.epochTo) {
    res += `&epoch-to=${p.epochTo}`;
  }
  return res;
};

export const listAllSchedulers = async (): Promise<Scheduler[]> => {
  const response = await fetch("/api/schedulers");
  const result: { schedulers: Scheduler[] } = await response.json();
  return result.schedulers;
};
export const makeScheduleModel = (schedules: any[]): ScheduleInfo[] => {
  if (schedules) {
    return schedules.map((o) => {
      return {
        id: o.id,
        scheduler: o.scheduler,
        timestamp: o.timestamp,
        epoch: o.epoch,
        targetTopic: o["target-topic"],
        targetId: o["target-id"],
      };
    });
  }
  return schedules;
};
export const searchLiveSchedules = async (
  p: SearchParams
): Promise<ScheduleInfo[]> => {
  const response = await fetch("/api/live/schedules" + makeSearchArgs(p));
  const result: { total: number; schedules: any[] } = await response.json();
  return makeScheduleModel(result.schedules);
};
export const searchSchedules = async (
  p: SearchParams
): Promise<ScheduleInfo[]> => {
  const response = await fetch("/api/schedules" + makeSearchArgs(p));
  const result: { total: number; schedules: any[] } = await response.json();
  return makeScheduleModel(result.schedules);
};
export const getScheduleDetail = async (
  scheduleName: string,
  id: string
): Promise<Schedule> => {
  const response = await fetch(`/api/scheduler/${scheduleName}/schedule/${id}`);
  const result: Schedule = await response.json();
  return result;
};
/*
/schedulers : return the list of addresses for each scheduler, containing one or many instance names
{
	schedulers: [{
		name: "kafka-message-scheduler.platform.svc.cluster.local",
		instances: [{
			ip: "10.1.25.23",
			names: ["10-1-25-23.scheduler-1.platform.svc.cluster.local."]
			topics: ["topic-1", "topic-2"]
			partitions: [0, 1]
			bootstrap-servers: ["kafka-1", "kafka-2", "kafka-3"]
		},
		{
			ip: "10.1.40.164"
			names: ["10.1.40.164.scheduler-2.platform.svc.cluster.local."]
			topics: ["topic-1", "topic-2"]
			partitions: [0, 1]
			bootstrap-servers: ["kafka-1", "kafka-2", "kafka-3"]
		}]
	},
	{
		name: "video-retrier-scheduler.platform.svc.cluster.local"
		instances: [{
			ip: "10.1.25.24"
			names: ["10-1-25-24.scheduler-vid.platform.svc.cluster.local."]
			topics: ["topic-2"]
			partitions: [0, 1]
			bootstrap-servers: ["kafka-1", "kafka-2", "kafka-3"]
		}]
	},
	{
		name: "10.1.25.24"
		instances: [{
			ip: "10.1.25.24"
			names: ["10-1-25-24.scheduler-vid.platform.svc.cluster.local."]
			topics: ["topic-3"]
			partitions: [0, 1]
			bootstrap-servers: ["kafka-4", "kafka-5", "kafka-6"]
		}]
	}]
}
/scheduler/{scheduler-name}/schedule/{id}: get a unique schedule
HTTP 302 Found
{
	id: "schedule-1",
	scheduler: "kafka-message-scheduler.platform.svc.cluster.local",
	timestamp: 589656366,
	epoch: 12563656,
	topic: "topic-1"
	target-topic: "xxxx",
	target-id: "xxx"
	headers: [
		{name: "", value: ""}
	],
	body: "xxxx",
}
HTTP 404 Not Found
<EMPTY>
/schedules: search for schedules based on search parameters
?schedule-id=xxx&epoch-from=0&epoch-to=1000&sort=id|epoch|timestamp asc|desc&scheduler-name=kafka-message-scheduler.platform.svc.cluster.local&max=-1
{
	found: 1500,
	schedules: [
		{
			id: "schedule-1",
			scheduler: "kafka-message-scheduler.platform.svc.cluster.local"
			epoch: 12563656,
			timestamp: 4589666
		}
	]
}
HTTP 404 Not Found
<EMPTY>
/live/schedules : returns the list of schedules in memory and planned for a scheduler (sorted by epoch)
?schedule-id=xxx&epoch-from=0&epoch-to=1000&sort=id|epoch|timestamp asc|desc&scheduler-name=kafka-message-scheduler.platform.svc.cluster.local&max=-1
HTTP 302 Found
{
	total: 150,
	schedules: [
		{
			id: "schedule-1",
			scheduler: "kafka-message-scheduler.platform.svc.cluster.local"
			epoch: 12563656,
			timestamp: 4589666
		}
	]
}
HTTP 404 Not Found
<EMPTY>

*/

import { get } from "_common/service/ApiUtil";
import { getLiveScheduleDetailUrl, getLiveSchedulesUrl, getScheduleDetailUrl, getSchedulersUrl, getSchedulesUrl } from "_core/service/config";
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
  return await get(getSchedulersUrl());
};
export const makeScheduleInfoModel = (schedules: any[]): ScheduleInfo[] => {
  if (schedules) {
    return schedules.map((o) => {
      return {
        id: o.schedule.id,
        scheduler: o.scheduler,
        timestamp: o.schedule.timestamp,
        epoch: o.schedule.epoch,
        targetTopic: o.schedule["target-topic"],
        targetId: o.schedule["target-key"],
		value: o.schedule.value
      };
    });
  }
  return schedules;
};

export const makeScheduleModel = (schedule: any, schedulerName:string): Schedule => {
	return {
		id: schedule.id,
		scheduler: schedulerName,
		timestamp: schedule.timestamp,
		epoch: schedule.epoch,
		targetTopic: schedule["target-topic"],
		targetId: schedule["target-key"],
		value: schedule.value,
		topic: schedule.topic
	  };
  };
export const searchLiveSchedules = async (p: SearchParams): Promise<ScheduleInfo[]> => {
  const result: { found: number; schedules: any[] } = await get(getLiveSchedulesUrl() + makeSearchArgs(p));
 
  const res = makeScheduleInfoModel(result.schedules);
  console.log(res);
  return res;
};
export const searchSchedules = async (p: SearchParams): Promise<ScheduleInfo[]> => {
  const result: { found: number; schedules: any[] } = await get(getSchedulesUrl() + makeSearchArgs(p));
  return makeScheduleInfoModel(result.schedules);
};
export const getScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule> => {
  const result: Schedule[]  = await get(getScheduleDetailUrl(schedulerName, id));
  
  if (result.length>0) {
	return makeScheduleModel(result[0], schedulerName);
  }
  throw new Error("Not found");
  
  
};

export const getLiveScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule> => {
	const result: Schedule[]  = await get(getLiveScheduleDetailUrl(schedulerName, id));
	
	if (result.length>0) {
	  return makeScheduleModel(result[0], schedulerName);
	}
	throw new Error("Not found");
	
	
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

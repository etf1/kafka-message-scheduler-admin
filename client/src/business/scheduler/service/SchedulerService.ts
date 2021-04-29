import { get } from "_common/service/ApiUtil";
import {
  getLiveScheduleDetailUrl,
  getLiveSchedulesUrl,
  getScheduleDetailUrl,
  getSchedulersUrl,
  getSchedulesUrl,
} from "_core/service/config";
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
        value: o.schedule.value,
      };
    });
  }
  return schedules;
};

export const makeScheduleModel = (schedule: any, schedulerName: string): Schedule => {
  return {
    id: schedule.id,
    scheduler: schedulerName,
    timestamp: schedule.timestamp,
    epoch: schedule.epoch,
    targetTopic: schedule["target-topic"],
    targetId: schedule["target-key"],
    value: schedule.value,
    topic: schedule.topic,
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
  const result: Schedule[] = await get(getScheduleDetailUrl(schedulerName, id));

  if (result.length > 0) {
    return makeScheduleModel(result[0], schedulerName);
  }
  throw new Error("Not found");
};

export const getLiveScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule> => {
  const result: Schedule[] = await get(getLiveScheduleDetailUrl(schedulerName, id));

  if (result.length > 0) {
    return makeScheduleModel(result[0], schedulerName);
  }
  throw new Error("Not found");
};

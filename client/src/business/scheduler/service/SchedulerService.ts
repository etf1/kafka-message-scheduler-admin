import { ScheduleType } from './../type/index';
import { get } from "_common/service/ApiUtil";
import {
  getAppStatsUrl,
  getLiveScheduleDetailUrl,
  getLiveSchedulesUrl,
  getHistoryScheduleDetailUrl,
  getHistorySchedulesUrl,
  getScheduleDetailUrl,
  getSchedulersUrl,
  getSchedulesUrl,
} from "_core/service/config";
import { Schedule, ScheduleInfo, Scheduler } from "../type";

export type SortType = "id" | "epoch" | "timestamp";
export type SortOrder = "asc" | "desc";
export type AppStat = {
  scheduler: string;
  total_live:number;
  history:number;
  total:number;
}

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
  let res = '?';//`?scheduler-name=${p.schedulerName}`;
  if (p.scheduleId) {
    res += `&schedule-id=${encodeURIComponent(p.scheduleId)}`;
  }
 /* if (p.max) {
    res += `&max=${p.max}`;
  }*/
  if (p.sort) {
    res += `&sort-by=${p.sort} ${p.sortOrder || "asc"}`;
  }
  if (p.epochFrom) {
    res += `&epoch-from=${encodeURIComponent(p.epochFrom)}`;
  }
  if (p.epochTo) {
    res += `&epoch-to=${encodeURIComponent(p.epochTo)}`;
  }
  return res;
};

export const getAppStats = async (): Promise<AppStat[]> => {
  return await get(getAppStatsUrl());
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

export const makeScheduleModel = ({ schedule }: any, schedulerName: string): Schedule => {
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
export const searchLiveSchedules = async (p: SearchParams): Promise<{found: number; schedules:ScheduleInfo[]}> => {
  const result: { found: number; schedules: any[] } = await get(
    getLiveSchedulesUrl(p.schedulerName) + makeSearchArgs(p)
  );

  const res = {found: result.found, schedules:makeScheduleInfoModel(result.schedules)};
  return res;
};
export const searchSchedules = async (p: SearchParams): Promise<{found: number; schedules:ScheduleInfo[]}> => {
  const result: { found: number; schedules: any[] } = await get(getSchedulesUrl(p.schedulerName) + makeSearchArgs(p));
  return {found: result.found, schedules:makeScheduleInfoModel(result.schedules)};
};
export const getScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule[]> => {
  const result: Schedule[] = await get(getScheduleDetailUrl(schedulerName, id));

  if (result.length > 0) {
    return result.map((sch) => makeScheduleModel(sch, schedulerName));
  }
  throw new Error("Not found");
};

export const getLiveScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule[]> => {
  
  const result: Schedule[] = await get(getLiveScheduleDetailUrl(schedulerName, id));

  if (result.length > 0) {
    return result.map((sch) => makeScheduleModel(sch, schedulerName));
  }
  throw new Error("Not found");
};

export const searchHistorySchedules = async (p: SearchParams): Promise<{found: number; schedules:ScheduleInfo[]}> => {
  const result: { found: number; schedules: any[] } = await get(
    getHistorySchedulesUrl(p.schedulerName) + makeSearchArgs(p)
  );

  const res = {found: result.found, schedules:makeScheduleInfoModel(result.schedules)};
  return res;
};
export const getHistoryScheduleDetail = async (schedulerName: string, id: string): Promise<Schedule[]> => {
  
  const result: Schedule[] = await get(getHistoryScheduleDetailUrl(schedulerName, id));

  if (result.length > 0) {
    return result.map((sch) => makeScheduleModel(sch, schedulerName));
  }
  throw new Error("Not found");
};

export function getScheduleDetailByType (type:ScheduleType) {
    switch (type) {
      case "history":
        return getHistoryScheduleDetail;
      case "live":
        return getLiveScheduleDetail;
      default:
        return getScheduleDetail;
    }
}


export function getSearchScheduleDetailByType (type:ScheduleType) {
  switch (type) {
    case "history":
      return searchHistorySchedules;
    case "live":
      return searchLiveSchedules;
    default:
      return searchSchedules;
  }
}
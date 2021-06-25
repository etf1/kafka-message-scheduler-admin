import { get } from "_common/service/ApiUtil";

let apiRoot = "";
let schedulersUrl = "";
let schedulesUrl = "";
let scheduleDetailUrl = "";
let liveSchedulesUrl = "";
let liveScheduleDetailUrl = "";
let historySchedulesUrl = "";
let historyScheduleDetailUrl = "";
let appStatUrl = "";

async function init() {
  const response = await get("/configuration.json");
  apiRoot = response["api-root"];
  appStatUrl= response["stats"];
  schedulersUrl = response["schedulers"];
  schedulesUrl = response["schedules"];
  scheduleDetailUrl = response["schedule-detail"];
  liveSchedulesUrl = response["live-schedules"];
  liveScheduleDetailUrl = response["live-schedule-detail"];
  historySchedulesUrl = response["history-schedules"];
  historyScheduleDetailUrl = response["history-schedule-detail"];
}
export function getApiRoot() {
  if (apiRoot === null || apiRoot === undefined) {
    return "";
    //throw new Error ("Error, configuration is not initialized, 'init()' function should be executed and terminated before any other calls.");
  }
  return apiRoot;
}
export function getAppStatsUrl() {
  return getApiRoot()+appStatUrl;
}
export function getSchedulersUrl() {
  return getApiRoot()+schedulersUrl;
}
export function getSchedulesUrl(schedulerName: string) {
  return getApiRoot()+schedulesUrl.replace("{name}", schedulerName);
}
export function getScheduleDetailUrl(schedulerName: string, id: string) {
  return getApiRoot()+scheduleDetailUrl.replace("{name}", schedulerName).replace("{id}", id);
}
export function getLiveSchedulesUrl(schedulerName: string) {
  return getApiRoot()+liveSchedulesUrl.replace("{name}", schedulerName);
}
export function getLiveScheduleDetailUrl(schedulerName: string, id: string) {
  return getApiRoot()+liveScheduleDetailUrl.replace("{name}", schedulerName).replace("{id}", id);
}
export function getHistorySchedulesUrl(schedulerName: string) {
  return getApiRoot()+historySchedulesUrl.replace("{name}", schedulerName);
}
export function getHistoryScheduleDetailUrl(schedulerName: string, id: string) {
  return getApiRoot()+historyScheduleDetailUrl.replace("{name}", schedulerName).replace("{id}", id);
}

export default init;

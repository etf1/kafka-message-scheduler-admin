import { get } from "_common/service/ApiUtil";

let apiRoot = "";
let schedulersUrl = "";
let schedulesUrl = "";
let scheduleDetailUrl = "";
let liveSchedulesUrl = "";
let liveScheduleDetailUrl = "";

async function init() {
  const response = await get("/configuration.json");
  apiRoot = response["api-root"];
  schedulersUrl = response["schedulers"];
  schedulesUrl = response["schedules"];
  scheduleDetailUrl = response["schedule-detail"];
  liveSchedulesUrl = response["live-schedules"];
  liveScheduleDetailUrl = response["live-schedule-detail"];
}
export function getApiRoot() {
  if (!apiRoot) {
    throw new Error ("Error, configuration is not initialized, 'init()' function should be executed and terminated before any other calls.");
  }
  return apiRoot;
}
export function getSchedulersUrl() {
  return getApiRoot()+schedulersUrl;
}
export function getSchedulesUrl() {
  return getApiRoot()+schedulesUrl;
}
export function getScheduleDetailUrl(schedulerName: string, id: string) {
  return getApiRoot()+scheduleDetailUrl.replace("{name}", schedulerName).replace("{id}", id);
}
export function getLiveSchedulesUrl() {
  return getApiRoot()+liveSchedulesUrl;
}
export function getLiveScheduleDetailUrl(schedulerName: string, id: string) {
  return getApiRoot()+liveScheduleDetailUrl.replace("{name}", schedulerName).replace("{id}", id);
}

export default init;

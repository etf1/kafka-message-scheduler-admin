import { Dictionary } from "./../../_common/type/utils";
import { replaceAll } from "_common/service/FunUtil";
import { lazy } from "react";

const Home = lazy(() => import("business/page/Home"));
const About = lazy(() => import("business/about/About"));

const SchedulesLive = lazy(() => import("business/page/SchedulesLive"));
const SchedulesAll = lazy(() => import("business/page/SchedulesAll"));
const ScheduleDetail = lazy(() => import("business/page/ScheduleDetail"));

export type RouteConfig = {
  path: string;
  key: string;
  component: React.LazyExoticComponent<() => JSX.Element>;
  exact: boolean;
};

export const ROUTE_HOME = "/";
export const ROUTE_ABOUT = "/about";
export const ROUTE_LIVE_SCHEDULES = "/live";
export const ROUTE_ALL_SCHEDULES = "/all";
export const ROUTE_SCHEDULE_DETAIL = "/detail/:schedulerName/:scheduleId";
export const ROUTE_SCHEDULE_LIVE_DETAIL = "/live/detail/:schedulerName/:scheduleId";
export const ROUTE_SCHEDULE_ALL_DETAIL = "/all/detail/:schedulerName/:scheduleId";

export const resolvePath = (path: string, variables: Dictionary) => {
  if (path.indexOf(":") > -1) {
    Object.keys(variables).forEach((key) => {
      path = replaceAll(path, ":" + key, variables[key]);
    });
  }
  return path;
};

const routes: RouteConfig[] = [
  {
    path: ROUTE_ABOUT,
    key: "about",
    component: About,
    exact: true,
  },
  {
    path: ROUTE_LIVE_SCHEDULES,
    key: "live",
    component: SchedulesLive,
    exact: true,
  },
  {
    path: ROUTE_ALL_SCHEDULES,
    key: "all",
    component: SchedulesAll,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULE_DETAIL,
    key: "schedule",
    component: ScheduleDetail,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULE_LIVE_DETAIL,
    key: "schedule",
    component: ScheduleDetail,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULE_ALL_DETAIL,
    key: "schedule",
    component: ScheduleDetail,
    exact: true,
  },
  {
    path: ROUTE_HOME,
    key: "home",
    component: Home,
    exact: false,
  },
];

export default routes;

import { ScheduleType } from "./../../business/scheduler/type/index";
import { Dictionary } from "./../../_common/type/utils";
import { replaceAll } from "_common/service/FunUtil";
import { lazy } from "react";
const Home = lazy(() => import("business/home/Home"));
const About = lazy(() => import("business/about/About"));
const Schedulers = lazy(() => import("business/scheduler/page/Schedulers"));
const SchedulerDetail = lazy(
  () => import("business/scheduler/page/SchedulerDetail")
);
const SchedulesLive = lazy(
  () => import("business/scheduler/page/SchedulesLive")
);
const SchedulesAll = lazy(() => import("business/scheduler/page/SchedulesAll"));
const ScheduleDetail = lazy(
  () => import("business/scheduler/page/ScheduleDetail")
);
const ScheduleDetailLive = lazy(
  () => import("business/scheduler/page/ScheduleDetailLive")
);
const SchedulesHistory = lazy(
  () => import("business/scheduler/page/SchedulesHistory")
);
const ScheduleDetailHistory = lazy(
  () => import("business/scheduler/page/ScheduleDetailHistory")
);
export type RouteConfig = {
  path: string;
  key: string;
  component: React.LazyExoticComponent<() => JSX.Element>;
  exact: boolean;
  menu?: { label: string; icon: string; position: number };
};

export const ROUTE_HOME = "/";
export const ROUTE_ABOUT = "/about";
export const ROUTE_SCHEDULERS = "/scheduler";
export const ROUTE_SCHEDULER_DETAIL = "/scheduler/detail/:schedulerName";
export const ROUTE_ALL_SCHEDULES = "/all";
export const ROUTE_SCHEDULE_ALL_DETAIL =
  "/all/detail/:schedulerName/:scheduleId";
export const ROUTE_LIVE_SCHEDULES = "/live";
export const ROUTE_SCHEDULE_LIVE_DETAIL =
  "/live/detail/:schedulerName/:scheduleId";
export const ROUTE_HISTORY_SCHEDULES = "/history";
export const ROUTE_SCHEDULE_HISTORY_DETAIL =
  "/history/detail/:schedulerName/:scheduleId";

export const getRouteScheduleDetailByType = (scheduleType: ScheduleType) => {
  switch (scheduleType) {
    case "live":
      return ROUTE_SCHEDULE_LIVE_DETAIL;
    case "history":
      return ROUTE_SCHEDULE_HISTORY_DETAIL;
    default:
      return ROUTE_SCHEDULE_ALL_DETAIL;
  }
};

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
    menu: {
      label: "Menu-schedules-live",
      icon: "bolt",
      position: 2,
    },
  },
  {
    path: ROUTE_HISTORY_SCHEDULES,
    key: "history",
    component: SchedulesHistory,
    exact: true,
    menu: {
      label: "Menu-schedules-history",
      icon: "history",
      position: 4,
    },
  },
  {
    path: ROUTE_ALL_SCHEDULES,
    key: "all",
    component: SchedulesAll,
    exact: true,
    menu: {
      label: "Menu-schedules-all",
      icon: "calendar-alt",
      position: 3,
    },
  },
  {
    path: ROUTE_SCHEDULE_LIVE_DETAIL,
    key: "schedule",
    component: ScheduleDetailLive,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULE_ALL_DETAIL,
    key: "schedule",
    component: ScheduleDetail,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULE_HISTORY_DETAIL,
    key: "schedule",
    component: ScheduleDetailHistory,
    exact: true,
  },
  {
    path: ROUTE_SCHEDULERS,
    key: "schedulers",
    component: Schedulers,
    exact: true,
    menu: {
      label: "Menu-schedulers",
      icon: "stopwatch",
      position: 5,
    },
  },
  {
    path: ROUTE_SCHEDULER_DETAIL,
    key: "scheduler-detail",
    component: SchedulerDetail,
    exact: true,
  },
  {
    path: ROUTE_HOME,
    key: "home",
    component: Home,
    exact: false,
    menu: {
      label: "Menu-home",
      icon: "home",
      position: 1,
    },
  },
];

export const routesWithMenu = [...routes.filter((r) => r.menu)];
routesWithMenu.sort(
  (a, b) => (a.menu && b.menu && a.menu?.position - b.menu?.position) || 0
);

export default routes;

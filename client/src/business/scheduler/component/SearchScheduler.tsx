import React, { useCallback, useEffect, useReducer, useState } from "react";
import { TFunction, useTranslation } from "react-i18next";
import format from "date-fns/format";
import add from "date-fns/add";
import Calendar from "_common/component/calendar/Calendar";
import Dropdown from "_common/component/element/Dropdown";
import SearchInput from "_common/component/element/SearchInput";
import { load, save } from "_common/service/StorageService";
import useSchedulers from "../hook/useSchedulers";
import { searchLiveSchedules, SearchParams, searchSchedules, SortOrder, SortType } from "../service/SchedulerService";
import { ScheduleInfo, Scheduler } from "../type";
import ScheduleTable from "./ScheduleTable";
import { ROUTE_SCHEDULE_LIVE_DETAIL, ROUTE_SCHEDULE_ALL_DETAIL } from "_core/router/routes";
import startOfDay from "date-fns/startOfDay";


type SearchParamsModel = {
  scheduler?: Scheduler;
  scheduleId?: string;
  epochFrom?: Date;
  epochTo?: Date;
  sort?: SortType;
  sortOrder?: SortOrder;
  max?: number; // -1  for all
};

export type SearchParamsReducerAction =
  | { type: "init"; payload: SearchParamsModel }
  | { type: "scheduler-changed"; payload: Scheduler }
  | { type: "scheduleId-changed"; payload: string }
  | { type: "epochFrom-changed"; payload: Date }
  | { type: "epochTo-changed"; payload: Date }
  | { type: "sort-changed"; payload: SortType }
  | { type: "sortOrder-changed"; payload: SortOrder }
  | { type: "max-changed"; payload: number };

export type SearchParamsReducer = (state: SearchParamsModel, action: SearchParamsReducerAction) => SearchParamsModel;

const searchParamsReducer: SearchParamsReducer = (state: SearchParamsModel, action) => {
  switch (action.type) {
    case "init":
      return { ...state, ...action.payload };
    case "scheduler-changed":
      return { ...state, scheduler: action.payload };
    case "scheduleId-changed":
      return { ...state, scheduleId: action.payload };
    case "epochFrom-changed":
      return { ...state, epochFrom: action.payload };
    case "epochTo-changed":
      return { ...state, epochTo: action.payload };
    case "sort-changed":
      return { ...state, sort: action.payload };
    case "sortOrder-changed":
      return { ...state, sortOrder: action.payload };
    case "max-changed":
      return { ...state, max: action.payload };
    default:
      throw new Error();
  }
};

const makeParams = (model: SearchParamsModel): SearchParams | undefined => {
  if (model.scheduler?.name) {
    return {
      scheduleId: model.scheduleId,
      epochFrom: model.epochFrom && parseInt((model.epochFrom.getTime() / 1000).toFixed(0)),
      epochTo: model.epochTo && parseInt((model.epochTo.getTime() / 1000).toFixed(0)),
      sort: model.sort,
      sortOrder: model.sortOrder,
      schedulerName: model.scheduler.name,
      max: model.max || 150,
    };
  } else {
    return undefined;
  }
};

const buildSearchModelLabel = (model: SearchParamsModel, t: TFunction<string>): React.ReactNode => {
  const result: React.ReactNode[] = [];
  const addSeparator = () => {
    if (result.length > 0) {
      result.push(
        <span key={result.length} className="space-right">
          ,
        </span>
      );
    }
  };
  const addLabel = (key: string, label: string, value: string) => {
    result.push(
      <span key={key} style={{ fontStyle: "italic" }}>
        <label style={{ fontStyle: "normal", fontWeight: 600 }}>{label}</label>: "{value}"
      </span>
    );
  };
  if (model) {
    if (model.scheduler) {
      addLabel("scheduler", t("Scheduler"), model.scheduler.name);
    }
    if (model.scheduleId) {
      addSeparator();
      addLabel("schedule-id", t("Scheduler-search-field-schedule-id"), model.scheduleId);
    }
    if (model.epochFrom) {
      addSeparator();
      addLabel("start-at", t("Scheduler-search-field-start-at"), format(model.epochFrom, t("Calendar-date-format")));
    }
    if (model.epochTo) {
      addSeparator();
      addLabel("end-at", t("Scheduler-search-field-end-at"), format(model.epochTo, t("Calendar-date-format")));
    }

    result.unshift(t("Scheduler-search-summary") + ": ");
  }

  return result;
};
export type SearchSchedulerProps = {
  live: boolean;
};
const SearchScheduler: React.FC<SearchSchedulerProps> = ({ live }) => {
  const { t } = useTranslation();
  const { schedulers } = useSchedulers();
  const [result, setResult] = useState<ScheduleInfo[]>([]);

  const [searchModel, dispatch] = useReducer<SearchParamsReducer>(searchParamsReducer, {
    scheduler: load<Scheduler>("SearchParamsModel-Scheduler", undefined),
    scheduleId: load<string>("SearchParamsModel-Scheduler-id", ""),
    epochFrom: startOfDay(new Date()),
    epochTo: startOfDay(add(new Date(), {
      days: 1,
    })),
  });

  useEffect(() => {
    if (searchModel) {
      save("SearchParamsModel-Scheduler", searchModel.scheduler);
      save("SearchParamsModel-Scheduler-id", searchModel.scheduleId);
    }
  }, [searchModel]);

  useEffect(() => {
    if (schedulers && schedulers.length > 0) {
      dispatch({ type: "scheduler-changed", payload: schedulers[0] });
    }
  }, [schedulers]);

  useEffect(() => {
    const searchMethod = live ? searchLiveSchedules : searchSchedules;
    const searchParams: SearchParams | undefined = makeParams(searchModel);
    if (searchParams) {
      searchMethod(searchParams).then((result) => {
        setResult(result);
      });
    }
  }, [searchModel, live]);

  const renderOption = (option: Scheduler) => {
    return <span key={option.name}>{option.name}</span>;
  };
  const handleSearchInputChanged = useCallback((value) => {
    dispatch({ type: "scheduleId-changed", payload: value || "" });
  }, []);
  return (
    <React.Fragment key="SearchScheduler">
      <h2 className="subtitle" style={{fontSize:"1rem"}}>{buildSearchModelLabel(searchModel, t)}</h2>
      <div className="app-box">
        <div className="container">
          <div className="panel">
            <div className="panel-heading">{t("Schedules")}</div>
            <div className="panel-block space-top more-space-bottom">
              <div className="field is-horizontal">
                <div className="field-body">
                  <div className="field">
                    <label className="label">{t("Scheduler")}</label>
                    <div className="control has-icons-left">
                      <Dropdown
                        placeholder={t("Please choose some scheduler")}
                        options={schedulers}
                        getKey={(scheduler) => scheduler.name}
                        renderOption={renderOption}
                        onChange={(s) => dispatch({ type: "scheduler-changed", payload: s })}
                        value={searchModel.scheduler}
                      />
                    </div>
                  </div>
                  <div className="field">
                    <label className="label">{t("Scheduler-search-field-schedule-id")}</label>
                    <div className="control">
                      <SearchInput
                        className="input"
                        onChange={handleSearchInputChanged}
                        value={searchModel.scheduleId}
                      />
                    </div>
                  </div>
                  <div className="field">
                    <label className="label">{t("Scheduler-search-field-start-at")}</label>
                    <div className="control">
                      <Calendar
                        uid="epochFrom"
                        className=""
                        onChange={(d) => dispatch({ type: "epochFrom-changed", payload: d })}
                        value={searchModel.epochFrom}
                      />
                    </div>
                  </div>
                  <div className="field">
                    <label className="label">{t("Scheduler-search-field-end-at")}</label>
                    <div className="control">
                      <Calendar
                        uid="epochTo"
                        className=""
                        onChange={(d) => dispatch({ type: "epochTo-changed", payload: d })}
                        value={searchModel.epochTo}
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div className="panel-block">
              <div className="container">
                {(!result || result.length === 0) && <strong>Pas de résultat...</strong>}
                {result && result.length > 0 && (
                  <ScheduleTable
                    key="table"
                    data={result}
                    detailUrl={live ? ROUTE_SCHEDULE_LIVE_DETAIL : ROUTE_SCHEDULE_ALL_DETAIL}
                  />
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </React.Fragment>
  );
};

export default SearchScheduler;

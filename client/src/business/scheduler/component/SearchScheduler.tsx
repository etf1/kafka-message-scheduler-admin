import React, { useCallback, useEffect, useReducer, useState } from "react";
import { useTranslation } from "react-i18next";
import Calendar from "_common/component/calendar/Calendar";
import Dropdown from "_common/component/element/Dropdown";
import SearchInput from "_common/component/element/SearchInput";
import { load, save } from "_common/service/StorageService";
import useSchedulers from "../hook/useSchedulers";
import { searchLiveSchedules, SearchParams, searchSchedules, SortOrder, SortType } from "../service/SchedulerService";
import { ScheduleInfo, Scheduler } from "../type";
import ScheduleTable from "./ScheduleTable";

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
      save<Scheduler>("selected-scheduler", action.payload)
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

export type SearchSchedulerProps = {
  live: boolean;
};
const SearchScheduler: React.FC<SearchSchedulerProps> = ({ live }) => {
  const { t } = useTranslation();
  const { schedulers, isLoading, refresh } = useSchedulers();
  const [result, setResult] = useState<ScheduleInfo[]>([]);

  const [searchModel, dispatch] = useReducer<SearchParamsReducer>(searchParamsReducer, { scheduler: load<Scheduler>("selected-scheduler", undefined), epochFrom: new Date() });

  useEffect (()=>{
    if (schedulers && schedulers.length>0) {
      dispatch({ type: "scheduler-changed", payload: schedulers[0]})
    }
  },[schedulers])

  useEffect(() => {
    const searchMethod = live ? searchLiveSchedules : searchSchedules;
    const searchParams: SearchParams | undefined = makeParams(searchModel);
    if (searchParams) {
      searchMethod(searchParams).then((result) => {
        setResult(result);
      });
    }
  }, [searchModel]);

  const handleResultRowClick = (schedule: ScheduleInfo) => {
    console.log("clicked");
  };

  const renderOption = (option: Scheduler) => {
    return <span>{option.name}</span>;
  };
  return (
    <div className="container">
      <div className="panel">
        <div className="panel-heading">{t("Schedules")}</div>
        <div className="panel-block more-space-bottom">
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
                <label className="label">{t("ID Planif.")}</label>
                <div className="control">
                  <SearchInput
                    className="input"
                    onChange={(value) => dispatch({ type: "scheduleId-changed", payload: value || "" })}
                    value={searchModel.scheduleId}
                  />
                </div>
              </div>
              <div className="field">
                <label className="label">{t("Start at")}</label>
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
                <label className="label">{t("End at")}</label>
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
        <div className="panel-block" >
          <div className="container">
            {(!result || result.length===0) && <strong>Pas de résultat...</strong>}
            {result && <ScheduleTable data={result} onClick={handleResultRowClick} />}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchScheduler;

import add from "date-fns/add";
import startOfDay from "date-fns/startOfDay";
import React, { useCallback, useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
import DatePicker from "_common/component/calendar/DatePicker";
import Dropdown from "_common/component/element/dropdown/Dropdown";
import SearchInput from "_common/component/element/search-input/SearchInput";
import { load, save } from "_common/service/StorageService";
import { getDateLocale } from "_core/i18n";
import useSchedulers from "../hook/useSchedulers";
import { SortOrder, SortType } from "../service/SchedulerService";
import { Scheduler } from "../type";

export type SearchParamsModel = {
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
  | { type: "epochFrom-changed"; payload: Date | undefined }
  | { type: "epochTo-changed"; payload: Date | undefined }
  | { type: "sort-changed"; payload: SortType }
  | { type: "sortOrder-changed"; payload: SortOrder }
  | { type: "max-changed"; payload: number };

export type SearchParamsReducer = (
  state: SearchParamsModel,
  action: SearchParamsReducerAction
) => SearchParamsModel;

const searchParamsReducer: SearchParamsReducer = (
  state: SearchParamsModel,
  action
) => {
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

type SearchSchedulerFormType = {
  onChange: (model: SearchParamsModel) => void;
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: Date;
  epochTo?: Date;
};
const SearchSchedulerForm: React.FC<SearchSchedulerFormType> = ({
  onChange,
  schedulerName,
  scheduleId,
  epochFrom,
  epochTo,
}) => {
  const { t } = useTranslation();
  const { schedulers } = useSchedulers();
  const [model, dispatch] = useReducer<SearchParamsReducer>(
    searchParamsReducer,
    {
      scheduler: load<Scheduler>(
        "SearchParamsModel-Scheduler",
        schedulers.find((s) => s.name === schedulerName) || undefined
      ),
      scheduleId: scheduleId || "",
      epochFrom: epochFrom || startOfDay(new Date()),
      epochTo:
        epochTo ||
        startOfDay(
          add(new Date(), {
            days: 1,
          })
        ),
    }
  );

  useEffect(() => {
    if (model) {
      save("SearchParamsModel-Scheduler", model.scheduler);
    }
  }, [model]);

  useEffect(() => {
    if (schedulers && schedulers.length > 0) {
      dispatch({ type: "scheduler-changed", payload: schedulers[0] });
    }
  }, [schedulers]);

  useEffect(() => {
    onChange(model);
  }, [model, onChange]);

  const renderOption = (option: Scheduler) => {
    return <span key={option.name}>{option.name}</span>;
  };
  const handleSearchInputChanged = useCallback((value) => {
    dispatch({ type: "scheduleId-changed", payload: value || "" });
  }, []);

  return (
    <div
      className="field is-horizontal"
      style={{ textAlign: "left", width: "100%", margin: "1rem" }}
    >
      <div className="field-label is-normal">
        <label className="label">Critères</label>
      </div>
      <div className="field-body columns is-mobile is-multiline">
        <div className="column">
          <Dropdown
            placeholder={t("Please choose some scheduler")}
            options={schedulers}
            getKey={(scheduler) => scheduler.name}
            renderOption={renderOption}
            onChange={(s) =>
              dispatch({ type: "scheduler-changed", payload: s })
            }
            value={model.scheduler}
          />
        </div>
        <div className="column" style={{ width: 150 }}>
          <SearchInput
            onChange={handleSearchInputChanged}
            placeholder={t("Scheduler-search-field-schedule-id")}
            value={model.scheduleId}
          />
        </div>
        <div className="column" style={{ flexGrow: 0 }}>
          <DatePicker
            placeholder={t("Scheduler-search-field-start-at")}
            value={model.epochFrom}
            onChange={(d) =>
              dispatch({ type: "epochFrom-changed", payload: d })
            }
            locale={getDateLocale()}
            dateFormat={t("Calendar-date-format")}
            todayLabel={t("Calendar-btn-label-Today")}
          />
        </div>
        <div className="column">
          <DatePicker
            placeholder={t("Scheduler-search-field-end-at")}
            value={model.epochTo}
            onChange={(d) => dispatch({ type: "epochTo-changed", payload: d })}
            locale={getDateLocale()}
            dateFormat={t("Calendar-date-format")}
            todayLabel={t("Calendar-btn-label-Today")}
          />
        </div>
      </div>
    </div>
  );
};

export default SearchSchedulerForm;

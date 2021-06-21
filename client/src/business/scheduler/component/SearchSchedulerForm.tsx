import endOfDay from "date-fns/endOfDay";
import React, { useCallback, useEffect, useReducer } from "react";
import { useTranslation } from "react-i18next";
import DatePicker from "_common/component/calendar/DatePicker";
import Icon from "_common/component/element/icon/Icon";
import SearchInput from "_common/component/element/search-input/SearchInput";
import Select from "_common/component/element/select/Select";
import { load, save } from "_common/service/LocalStorageService";
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
      return { ...state, epochTo: action.payload && endOfDay(action.payload) };
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
  onRefresh: () => void;
};
const SearchSchedulerForm: React.FC<SearchSchedulerFormType> = ({
  onChange,
  schedulerName,
  scheduleId,
  epochFrom,
  epochTo,
  onRefresh,
}) => {
  const { t } = useTranslation();
  const { schedulers } = useSchedulers();
  const [model, dispatch] = useReducer<SearchParamsReducer>(searchParamsReducer, {
    scheduler: load<Scheduler>(
      "SearchParamsModel-Scheduler",
      (schedulers && schedulers.find((s) => s.name === schedulerName)) || undefined
    ),
    scheduleId: scheduleId || "",
    epochFrom: epochFrom, //|| startOfDay(new Date()),
    epochTo: epochTo /* ||
      endOfDay(
        add(new Date(), {
          days: 1,
        })
      ),*/,
  });

  useEffect(() => {
    if (model) {
      save("SearchParamsModel-Scheduler", model.scheduler);
    }
  }, [model]);

  useEffect(() => {
    if (schedulers && schedulers.length > 0) {
      dispatch({ type: "scheduler-changed", payload: schedulers.find (s => s.name === schedulerName) || schedulers[0] });
    }
  }, [schedulers, schedulerName]);

  useEffect(() => {
    onChange(model);
  }, [model, onChange]);

  const handleSearchInputChanged = useCallback((value) => {
    dispatch({ type: "scheduleId-changed", payload: value || "" });
  }, []);

  return (
    <div className="field " style={{ textAlign: "left", width: "100%", margin: "0" }}>
      <div className=" columns is-mobile is-multiline">
        <div className="column is-3">
          <div className={"field fieldWithNoLabel"}>
            <label className="label">{t("Scheduler")}</label>
            <div className={"control"}>
              <Select
                value={model.scheduler}
                onChange={(s) => s && dispatch({ type: "scheduler-changed", payload: s })}
                className="column is-3"
                labelField={"name"}
                keyField={"name"}
                options={schedulers}
              />
            </div>
          </div>
        </div>
        <div className="column is-4">
          <label className="label">ID Planif.</label>
          <SearchInput
            onChange={handleSearchInputChanged}
            placeholder={t("Scheduler-search-field-schedule-id")}
            value={model.scheduleId}
          />
        </div>
        <div className="column" style={{ flexGrow: 0 }}>
          <label className="label">Début</label>
          <DatePicker
            placeholder={t("Scheduler-search-field-start-at")}
            value={model.epochFrom}
            onChange={(d) => dispatch({ type: "epochFrom-changed", payload: d })}
            locale={getDateLocale()}
            dateFormat={t("Calendar-date-format")}
            todayLabel={t("Calendar-btn-label-Today")}
          />
        </div>
        <div className="column">
          <label className="label">Fin</label>
          <DatePicker
            placeholder={t("Scheduler-search-field-end-at")}
            value={model.epochTo}
            onChange={(d) => dispatch({ type: "epochTo-changed", payload: d })}
            locale={getDateLocale()}
            dateFormat={t("Calendar-date-format")}
            todayLabel={t("Calendar-btn-label-Today")}
          />
        </div>
        <div className="column">
          <label className="label">&nbsp;</label>
          <button onClick={onRefresh} className="button is-primary">
            <Icon name="sync-alt" marginRight={10} /> {t("Refresh")}
          </button>
        </div>
      </div>
    </div>
  );
};

export default SearchSchedulerForm;

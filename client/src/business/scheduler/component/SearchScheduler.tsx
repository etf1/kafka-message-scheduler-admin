import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import format from "date-fns/format";
import { searchLiveSchedules, SearchParams, searchSchedules, SortOrder, SortType } from "../service/SchedulerService";
import { ScheduleInfo } from "../type";
import ScheduleTable from "./ScheduleTable";
import { ROUTE_SCHEDULE_LIVE_DETAIL, ROUTE_SCHEDULE_ALL_DETAIL } from "_core/router/routes";
import useMedia from "_common/hook/useMedia";
import SearchSchedulerForm, { SearchParamsModel } from "./SearchSchedulerForm";
import { useHistory } from "react-router";
import { pluralizeIf } from "_core/i18n";
import Container from "_common/component/layout/container/Container";
import startOfDay from "date-fns/startOfDay";
import endOfDay from "date-fns/endOfDay";

const makeParams = (model: SearchParamsModel | undefined): SearchParams | undefined => {
  if (model && model.scheduler?.name) {
    return {
      scheduleId: model.scheduleId,
      epochFrom: model.epochFrom && parseInt((model.epochFrom.getTime() / 1000).toFixed(0)),
      epochTo: model.epochTo && parseInt((model.epochTo.getTime() / 1000).toFixed(0)),
      sort: model.sort,
      sortOrder: model.sortOrder,
      schedulerName: model.scheduler.name,
      max: model.max || 300,
    };
  } else {
    return undefined;
  }
};
export type SearchSchedulerProps = {
  live: boolean;
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: Date;
  epochTo?: Date;
};

const SearchScheduler: React.FC<SearchSchedulerProps> = ({ live, schedulerName, scheduleId, epochFrom, epochTo }) => {
  const { t } = useTranslation();
  const history = useHistory();
  const [searchModel, setSearchModel] = useState<SearchParamsModel | undefined>(); //;load<SearchParamsModel>("SearchParamsModel"+live?"live":"all", undefined));
  const [result, setResult] = useState<{ found: number; schedules: ScheduleInfo[] }>();
  const smallScreen = useMedia(["(max-width: 1250px)", "(min-width: 1250px)"], [true, false], true);
  const schedules = result?.schedules || [];

  const buildResultLabel = () => {
    if (result && result.found > 0) {
      const limitedResultLabel =
        schedules.length < (result?.found || 0)
          ? `(${t("Schedule-Search-limited-result-label")} ${schedules.length})`
          : "";
      return `${result.found} ${pluralizeIf(
        schedules.length,
        t("Schedule-Search-result"),
        t("Schedule-Search-results")
      )} ${limitedResultLabel}`;
    }
    return "";
  };

  useEffect(() => {
    const searchMethod = live ? searchLiveSchedules : searchSchedules;
    //  save("SearchParamsModel"+live?"live":"all", searchModel);
    //  console.log(searchModel);
    const searchParams: SearchParams | undefined = makeParams(searchModel);
    if (searchParams) {
      searchMethod(searchParams).then((result) => {
        setResult(result);
      });
    }
  }, [searchModel, live]);

  const handleSearchChange = useCallback(
    (searchModel: SearchParamsModel) => {
      const newPath = [];
      if (searchModel.scheduler) {
        newPath.push(`schedulerName=${searchModel.scheduler.name}`);
      }
      if (searchModel.scheduleId) {
        newPath.push(`scheduleId=${searchModel.scheduleId}`);
      }
      if (searchModel.epochFrom) {
        newPath.push(`epochFrom=${format(searchModel.epochFrom, t("Calendar-date-format"))}`);
      }
      if (searchModel.epochTo) {
        newPath.push(`epochTo=${format(searchModel.epochTo, t("Calendar-date-format"))}`);
      }
      history.replace(window.location.pathname + "?" + newPath.join("&"));
      setSearchModel(searchModel);
    },
    [history, t]
  );

  const handleSort = useCallback(
    (type: SortType, order: SortOrder) => {
      if (searchModel && (searchModel.sort !== type || searchModel.sortOrder !== order)) {
        searchModel.sort = type;
        searchModel.sortOrder = order;
        setSearchModel({ ...searchModel });
      }
    },
    [searchModel]
  );

  const handleRefresh = useCallback(() => {
    setSearchModel((m) => ({ ...m }));
  }, []);
  return (
    <React.Fragment key="SearchScheduler">
      <div className="app-box">
        <div className="container">
          <div className="more-space-top more-space-bottom">
            <SearchSchedulerForm
              onChange={handleSearchChange}
              schedulerName={schedulerName}
              scheduleId={scheduleId}
              epochFrom={epochFrom || (live ? undefined : startOfDay(new Date()))}
              epochTo={epochTo || (live ? undefined : endOfDay(new Date()))}
              onRefresh={handleRefresh}
            />
          </div>
          <Container title={buildResultLabel()}>
            {(!schedules || schedules.length === 0) && <strong>Pas de résultat...</strong>}
            {schedules && schedules.length > 0 && (
              <ScheduleTable
                key="table"
                data={schedules}
                showAsTable={!smallScreen}
                onSort={handleSort}
                detailUrl={live ? ROUTE_SCHEDULE_LIVE_DETAIL : ROUTE_SCHEDULE_ALL_DETAIL}
              />
            )}
          </Container>
        </div>
      </div>
    </React.Fragment>
  );
};

export default SearchScheduler;

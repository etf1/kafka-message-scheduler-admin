import React, { useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import format from "date-fns/format";
import { getSearchScheduleDetailByType, SearchParams, SortOrder, SortType } from "../service/SchedulerService";
import { ScheduleInfo, ScheduleType } from "../type";
import ScheduleTable from "./ScheduleTable";
import { getRouteScheduleDetailByType } from "_core/router/routes";
import useMedia from "_common/hook/useMedia";
import SearchSchedulerForm, { SearchParamsModel } from "./SearchSchedulerForm";
import { useHistory } from "react-router";
import { pluralizeIf } from "_core/i18n";
import Container from "_common/component/layout/container/Container";
import Appear from "_common/component/transition/Appear";
import { save } from "_common/service/SessionStorageService";
import clsx from "clsx";
import Icon from "_common/component/element/icon/Icon";
import useRefresh from "_common/hook/useRefresh";

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
  scheduleType: ScheduleType;
  schedulerName?: string;
  scheduleId?: string;
  epochFrom?: Date;
  epochTo?: Date;
};

const SearchScheduler: React.FC<SearchSchedulerProps> = ({
  scheduleType,
  schedulerName,
  scheduleId,
  epochFrom,
  epochTo,
}) => {
  const { t } = useTranslation();
  const history = useHistory();
  const [searchModel, setSearchModel] = useState<SearchParamsModel | undefined>(); //;load<SearchParamsModel>("SearchParamsModel"+live?"live":"all", undefined));
  const [result, setResult] = useState<{ found: number; schedules: ScheduleInfo[] }>();
  const smallScreen = useMedia(["(max-width: 1250px)", "(min-width: 1250px)"], [true, false], true);
  const schedules = result?.schedules || [];
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error>();
  const [refresh, count] = useRefresh();

  const prevSearhParamStr = useRef<string>();
  const prevCount = useRef<number>();
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
    const searchParams: SearchParams | undefined = makeParams(searchModel);
    const searchParamStr = JSON.stringify(searchParams);
    if (
      (searchParams && count !== prevCount.current) ||
      (searchParams && searchParamStr !== prevSearhParamStr.current)
    ) {
      setIsLoading(true);
      prevSearhParamStr.current = searchParamStr;
      getSearchScheduleDetailByType(scheduleType)(searchParams)
        .then((result) => {
          setResult(result);
          setIsLoading(false);
          setError(undefined);
        })
        .catch((err: Error) => {
          console.error(err);
          setError(err);
        });
    }
  }, [searchModel, scheduleType, count]);

  const handleSearchChange = useCallback(
    (searchModel: SearchParamsModel) => {
      const newPath = [];
      if (searchModel.scheduler) {
        newPath.push(`schedulerName=${searchModel.scheduler.name}`);
        save(scheduleType + "SchedulerName", searchModel.scheduler.name);
      }
      if (searchModel.scheduleId) {
        newPath.push(`scheduleId=${searchModel.scheduleId}`);
      }
      save(scheduleType + "ScheduleId", searchModel.scheduleId);

      const epochFrom = searchModel.epochFrom && format(searchModel.epochFrom, t("Calendar-date-format"));
      save(scheduleType + "EpochFrom", epochFrom);
      if (epochFrom) {
        newPath.push(`epochFrom=${epochFrom}`);
      }
      const epochTo = searchModel.epochTo && format(searchModel.epochTo, t("Calendar-date-format"));
      save(scheduleType + "EpochTo", epochTo);
      if (epochTo) {
        newPath.push(`epochTo=${epochTo}`);
      }

      history.replace(window.location.pathname + "?" + newPath.join("&"));
      setSearchModel(searchModel);
    },
    [history, scheduleType, t]
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

  return (
    <React.Fragment key="SearchScheduler">
      <div className="app-box">
        <div className="container">
          <div style={{  paddingBottom: 0 }}>
            <div className="space-top space-bottom">
              <SearchSchedulerForm
                onChange={handleSearchChange}
                schedulerName={schedulerName}
                scheduleId={scheduleId}
                epochFrom={epochFrom || undefined}
                epochTo={epochTo || undefined}
                onRefresh={refresh}
              />
            </div>
          </div>
          <hr style={{marginLeft:-20, width:"133%"}}/>
          <Container title={buildResultLabel()}>
            {(!schedules || schedules.length === 0) && (
              <strong style={{ color: "gray", fontStyle: "italic" }}>
                {isLoading ? t("Loading") : t("NoResults")}
              </strong>
            )}
            <div>
              {error && (
                <div className="animate-opacity" style={{ fontWeight: 800, color: "red" }}>
                  <Icon name="exclamation-triangle" /> {t("LoadingError")}
                </div>
              )}
              {!error && (
                <Appear visible={schedules && schedules.length > 0}>
                  {(nodeRef) => (
                    <div ref={nodeRef} className={clsx(isLoading && "animate-opacity")}>
                      <ScheduleTable
                        key="table"
                        data={schedules}
                        showAsTable={!smallScreen}
                        onSort={handleSort}
                        detailUrl={getRouteScheduleDetailByType(scheduleType)}
                      />
                    </div>
                  )}
                </Appear>
              )}
            </div>
          </Container>
        </div>
      </div>
    </React.Fragment>
  );
};

export default SearchScheduler;

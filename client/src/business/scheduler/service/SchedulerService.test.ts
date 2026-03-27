import { makeScheduleInfoModel, makeScheduleModel, makeSearchArgs, type SearchParams } from './SchedulerService';

// ---------------------------------------------------------------------------
// makeSearchArgs
// ---------------------------------------------------------------------------

describe('makeSearchArgs', () => {
  it('returns bare "?" when no optional params are provided', () => {
    const p: SearchParams = { schedulerName: 'sched-1', max: 300 };
    expect(makeSearchArgs(p)).toBe('?');
  });

  it('encodes scheduleId into the query string', () => {
    const p: SearchParams = { schedulerName: 'sched-1', max: 300, scheduleId: 'invoice/FR-001' };
    const result = makeSearchArgs(p);
    expect(result).toContain('schedule-id=invoice%2FFR-001');
  });

  it('includes sort-by when sort and sortOrder are provided', () => {
    const p: SearchParams = {
      schedulerName: 'sched-1',
      max: 300,
      sort: 'timestamp',
      sortOrder: 'desc',
    };
    const result = makeSearchArgs(p);
    expect(result).toContain('sort-by=timestamp desc');
  });

  it('defaults sort order to "asc" when sortOrder is omitted', () => {
    const p: SearchParams = { schedulerName: 'sched-1', max: 300, sort: 'id' };
    const result = makeSearchArgs(p);
    expect(result).toContain('sort-by=id asc');
  });

  it('encodes epochFrom and epochTo', () => {
    const p: SearchParams = {
      schedulerName: 'sched-1',
      max: 300,
      epochFrom: 1_700_000_000,
      epochTo: 1_800_000_000,
    };
    const result = makeSearchArgs(p);
    expect(result).toContain('epoch-from=1700000000');
    expect(result).toContain('epoch-to=1800000000');
  });

  it('combines all optional params', () => {
    const p: SearchParams = {
      schedulerName: 'sched-1',
      max: 50,
      scheduleId: 'sched-42',
      sort: 'epoch',
      sortOrder: 'asc',
      epochFrom: 1_700_000_000,
      epochTo: 1_800_000_000,
    };
    const result = makeSearchArgs(p);
    expect(result).toContain('schedule-id=sched-42');
    expect(result).toContain('sort-by=epoch asc');
    expect(result).toContain('epoch-from=1700000000');
    expect(result).toContain('epoch-to=1800000000');
  });
});

// ---------------------------------------------------------------------------
// makeScheduleInfoModel
// ---------------------------------------------------------------------------

// Raw API shape returned by the Go server (JSON field names use kebab-case for some keys)
const rawSchedule = (id: string, schedulerName: string) => ({
  scheduler: schedulerName,
  schedule: {
    id,
    timestamp: 1_700_000_001,
    epoch: 1_700_001_000,
    'target-topic': 'invoice.output',
    'target-key': `key-${id}`,
    value: 'base64payload==',
  },
});

describe('makeScheduleInfoModel', () => {
  it('maps raw API objects to ScheduleInfo shape', () => {
    const raw = [rawSchedule('schedule-1', 'sched-1'), rawSchedule('schedule-2', 'sched-1')];
    const result = makeScheduleInfoModel(raw);

    expect(result).toHaveLength(2);

    expect(result[0]).toMatchObject({
      id: 'schedule-1',
      scheduler: 'sched-1',
      timestamp: 1_700_000_001,
      epoch: 1_700_001_000,
      targetTopic: 'invoice.output',
      targetId: 'key-schedule-1',
    });
  });

  it('returns an empty array for an empty input', () => {
    expect(makeScheduleInfoModel([])).toEqual([]);
  });

  it('passes through null/undefined input unchanged (API consistency guard)', () => {
    // The Go API should never return null, but the function is permissive.
    const result = makeScheduleInfoModel(null as any);
    expect(result).toBeNull();
  });
});

// ---------------------------------------------------------------------------
// makeScheduleModel
// ---------------------------------------------------------------------------

describe('makeScheduleModel', () => {
  it('maps a raw schedule detail response to the Schedule model', () => {
    const raw = rawSchedule('schedule-99', 'sched-1');
    const result = makeScheduleModel(raw, 'sched-1');

    expect(result).toMatchObject({
      id: 'schedule-99',
      scheduler: 'sched-1',
      timestamp: 1_700_000_001,
      epoch: 1_700_001_000,
      targetTopic: 'invoice.output',
      targetId: 'key-schedule-99',
      value: 'base64payload==',
    });
  });

  it('uses the provided schedulerName (not the nested scheduler field)', () => {
    const raw = rawSchedule('schedule-1', 'original-scheduler');
    // The Go detail endpoint wraps result in {schedule: {...}} — scheduler comes from caller.
    const result = makeScheduleModel(raw, 'override-scheduler');
    expect(result.scheduler).toBe('override-scheduler');
  });
});

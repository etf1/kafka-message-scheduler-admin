import { HttpResponse, http } from 'msw';

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

// sample from https://mswjs.io/docs/getting-started/mocks/rest-api
export const handlers = [
  http.get('/api/schedulers', () => {
    return HttpResponse.json({
      schedulers: [
        {
          name: 'kafka-message-scheduler.platform.svc.cluster.local',
          instances: [
            {
              ip: '42.42.25.23',
              names: ['42-42-25-23.scheduler-1.platform.svc.cluster.local.'],
              topics: ['topic-1', 'topic-2'],
              partitions: [0, 1],
              bootstrapServers: ['kafka-1', 'kafka-2', 'kafka-3'],
            },
            {
              ip: '42.42.40.164',
              names: ['42.42.40.164.scheduler-2.platform.svc.cluster.local.'],
              topics: ['topic-1', 'topic-2'],
              partitions: [0, 1],
              bootstrapServers: ['kafka-1', 'kafka-2', 'kafka-3'],
            },
          ],
        },
        {
          name: 'video-retrier-scheduler.platform.svc.cluster.local',
          instances: [
            {
              ip: '42.42.25.24',
              names: ['42-42-25-24.scheduler-vid.platform.svc.cluster.local.'],
              topics: ['topic-2'],
              partitions: [0, 1],
              bootstrapServers: ['kafka-1', 'kafka-2', 'kafka-3'],
            },
          ],
        },
        {
          name: '42.42.25.24',
          instances: [
            {
              ip: '42.42.25.24',
              names: ['42-42-25-24.scheduler-vid.platform.svc.cluster.local.'],
              topics: ['topic-3'],
              partitions: [0, 1],
              bootstrapServers: ['kafka-4', 'kafka-5', 'kafka-6'],
            },
          ],
        },
      ],
    });
  }),
  http.get('/api/live/schedules', ({ request }) => {
    const url = new URL(request.url);
    const schedulerName = url.searchParams.get('scheduler-name');
    const scheduleId = url.searchParams.get('schedule-id');
    const max = scheduleId ? 8 - scheduleId.length : Number(url.searchParams.get('max') ?? 300);
    const schedules = [];
    if (!scheduleId || scheduleId.length < 8) {
      for (let i = 0; i < max; i++) {
        schedules.push({
          id: `video:ce67bd${getRandomInt(9)}${getRandomInt(9)}-4997-441a-b762-a29e1cc8c44${i}:apps:offline`,
          scheduler: schedulerName,
          epoch: 1615532584,
          timestamp: 1605542514,
          'target-topic': 'backend.queueing.catalog.video.v1',
          'target-id': 'ce67bd87-4997-441a-b762-a29e1cc8c446',
        });
      }
    }
    return HttpResponse.json({ total: max, schedules });
  }),
  http.get('/api/schedules', ({ request }) => {
    const url = new URL(request.url);
    const schedulerName = url.searchParams.get('scheduler-name');
    const scheduleId = url.searchParams.get('schedule-id');
    const max = scheduleId ? 8 - scheduleId.length : Number(url.searchParams.get('max') ?? 300);
    const schedules = [];
    if (!scheduleId || scheduleId.length < 8) {
      for (let i = 0; i < max; i++) {
        schedules.push({
          id: `video:ce67bd${getRandomInt(9)}${getRandomInt(9)}-4997-441a-b762-a29e1cc8c44${i}:apps:offline`,
          scheduler: schedulerName,
          epoch: 1615532584,
          timestamp: 1605542514,
          'target-topic': 'backend.queueing.catalog.video.v1',
          'target-id': 'ce67bd87-4997-441a-b762-a29e1cc8c446',
        });
      }
    }
    return HttpResponse.json({ total: max, schedules });
  }),
  http.get('/api/scheduler/:schedulerName/schedule/:id', ({ params }) => {
    const { schedulerName, id } = params;
    return HttpResponse.json({
      id,
      scheduler: schedulerName,
      epoch: 1615532584,
      timestamp: 1605542514,
      topic: 'topic-1',
      'target-topic': 'backend.queueing.catalog.video.v1',
      'target-id': `ce67bd8${getRandomInt(9)}-4997-441a-b762-a29e1cc8c446`,
      headers: [
        { name: 'header-name', value: 'a-header-value' },
        { name: 'another-header-name', value: 'xxx-value' },
      ],
      body: 'xxxx',
    });
  }),
];

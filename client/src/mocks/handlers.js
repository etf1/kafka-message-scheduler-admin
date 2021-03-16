import { rest } from "msw";

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

// sample from https://mswjs.io/docs/getting-started/mocks/rest-api
export const handlers = [
  rest.get("/api/schedulers", (req, res, ctx) => {
    return res(
      // Respond with a 200 status code
      // eslint-disable-next-line
      ctx.status(200),
      ctx.json({
        schedulers: [
          {
            name: "kafka-message-scheduler.platform.svc.cluster.local",
            instances: [
              {
                ip: "42.42.25.23",
                names: ["42-42-25-23.scheduler-1.platform.svc.cluster.local."],
                topics: ["topic-1", "topic-2"],
                partitions: [0, 1],
                bootstrapServers: ["kafka-1", "kafka-2", "kafka-3"],
              },
              {
                ip: "42.42.40.164",
                names: ["42.42.40.164.scheduler-2.platform.svc.cluster.local."],
                topics: ["topic-1", "topic-2"],
                partitions: [0, 1],
                bootstrapServers: ["kafka-1", "kafka-2", "kafka-3"],
              },
            ],
          },
          {
            name: "video-retrier-scheduler.platform.svc.cluster.local",
            instances: [
              {
                ip: "42.42.25.24",
                names: ["42-42-25-24.scheduler-vid.platform.svc.cluster.local."],
                topics: ["topic-2"],
                partitions: [0, 1],
                bootstrapServers: ["kafka-1", "kafka-2", "kafka-3"],
              },
            ],
          },
          {
            name: "42.42.25.24",
            instances: [
              {
                ip: "42.42.25.24",
                names: ["42-42-25-24.scheduler-vid.platform.svc.cluster.local."],
                topics: ["topic-3"],
                partitions: [0, 1],
                bootstrapServers: ["kafka-4", "kafka-5", "kafka-6"],
              },
            ],
          },
        ],
      })
    );
  }),
  rest.get("/api/live/schedules", (req, res, ctx) => {
    const schedulerName = req.url.searchParams.get("scheduler-name");
    const scheduleId = req.url.searchParams.get("schedule-id");
    const max = scheduleId ? (8-scheduleId.length) : req.url.searchParams.get("max") || 150;
    /*const sort = req.url.searchParams.get("sort");
    const epochFrom = req.url.searchParams.get("epoch-from");
    const epochTo = req.url.searchParams.get("epoch-to");*/
    const schedules = [];
    if (!scheduleId || scheduleId.length < 8) {
      for (let i = 0; i < max; i++) {
        schedules.push({
          id: `video:ce67bd${getRandomInt(9)}${getRandomInt(9)}-4997-441a-b762-a29e1cc8c44${i}:apps:offline`,
          scheduler: schedulerName,
          epoch: 1615532584,
          timestamp: 1605542514,

          "target-topic": "backend.queueing.catalog.video.v1",
          "target-id": "ce67bd87-4997-441a-b762-a29e1cc8c446",
        });
      }
    }
    return res(
      // Respond with a 200 status code
      ctx.status(200),
      ctx.json({
        total: max,
        schedules: schedules,
      })
    );
  }),
  rest.get("/api/schedules", (req, res, ctx) => {
    const schedulerName = req.url.searchParams.get("scheduler-name");
    const scheduleId = req.url.searchParams.get("schedule-id");
    const max = scheduleId ? (8-scheduleId.length) : req.url.searchParams.get("max") || 150;
    //const sort = req.url.searchParams.get("sort");
    //const epochFrom = req.url.searchParams.get("epoch-from");
    //const epochTo = req.url.searchParams.get("epoch-to");
    const schedules = [];
    if (!scheduleId || scheduleId.length < 8) {
      for (let i = 0; i < max; i++) {
        schedules.push({
          id: `video:ce67bd${getRandomInt(9)}${getRandomInt(9)}-4997-441a-b762-a29e1cc8c44${i}:apps:offline`,
          scheduler: schedulerName,
          epoch: 1615532584,
          timestamp: 1605542514,
          "target-topic": "backend.queueing.catalog.video.v1",
          "target-id": "ce67bd87-4997-441a-b762-a29e1cc8c446",
        });
      }
    }

    return res(
      // Respond with a 200 status code
      ctx.status(200),
      ctx.json({
        total: max,
        schedules: schedules,
      })
    );
  }),
  rest.get("/api/scheduler/:schedulerName/schedule/:id", (req, res, ctx) => {
    const { schedulerName, id } = req.params;
    return res(
      // Respond with a 200 status code
      ctx.status(200),
      ctx.json({
        id: id,
        scheduler: schedulerName,
        epoch: 1615532584,
        timestamp: 1605542514,
        topic: "topic-1",
        "target-topic": "backend.queueing.catalog.video.v1",
        "target-id": `ce67bd8${getRandomInt(9)}-4997-441a-b762-a29e1cc8c446`,
        headers: [{ name: "header-name", value: "a-header-value" },{ name: "another-header-name", value: "xxx-value" }],
        body: "xxxx",
      })
    );
  }),
];

# API Reference — `kafka-message-scheduler-admin`

> **Base URL:** `http://localhost:9000`  
> **API prefix:** `/api` (when SPA mode, i.e. `API_SERVER_ONLY=false`)  
> **Content-Type:** `application/json; charset=UTF-8`  
> **Auth:** None (internal admin tool, protected by network/VPN)

---

## Overview

The Kafka Message Scheduler Admin exposes a read-only REST API to query scheduled messages across multiple Kafka scheduler instances. Data is sourced from three stores:

| Store | Prefix | Description |
|-------|--------|-------------|
| **Cold** | `/scheduler/{name}/...` | Persistent (BoltDB + Bleve index), all known schedules |
| **Live** | `/live/scheduler/{name}/...` | Real-time from scheduler instances via REST |
| **History** | `/history/scheduler/{name}/...` | Persistent (BoltDB + Bleve index), past/triggered schedules |

---

## Endpoints

### `GET /api/stats`

Returns aggregate statistics for all known schedulers.

**Response:** `200 OK`

```json
[
  {
    "scheduler": "scheduler-1",
    "total_live": 42,
    "total_history": 150,
    "total": 300
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| `scheduler` | string | Scheduler instance name |
| `total_live` | int | Schedules currently live in the scheduler instance |
| `total_history` | int | Historical/triggered schedules |
| `total` | int | Total schedules in cold storage |

---

### `GET /api/schedulers`

Lists all discovered scheduler instances with their Kafka configuration.

**Response:** `200 OK`

```json
[
  {
    "name": "scheduler-1",
    "http_port": "8000",
    "instances": [
      {
        "ip": "10.0.0.1",
        "hostname": ["scheduler-1.local"],
        "topics": ["schedules"],
        "history_topic": "history",
        "bootstrap_servers": "kafka:9092"
      }
    ]
  }
]
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Scheduler hostname |
| `http_port` | string | HTTP port for scheduler REST API |
| `instances[].ip` | string | Instance IP address |
| `instances[].hostname` | string[] | Reverse DNS hostnames |
| `instances[].topics` | string[] | Kafka topics consumed by this scheduler |
| `instances[].history_topic` | string | Kafka topic for historical schedules |
| `instances[].bootstrap_servers` | string | Kafka bootstrap servers |

---

### `GET /api/scheduler/{name}/schedules`

Search schedules in **cold storage** (persistent, indexed).

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | string | Scheduler name |

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `schedule-id` | string | — | Filter by schedule ID (supports wildcards `*` in Bleve) |
| `epoch-from` | int64 | — | Minimum epoch (Unix timestamp, inclusive) |
| `epoch-to` | int64 | — | Maximum epoch (Unix timestamp, inclusive) |
| `sort-by` | string | `timestamp desc` | Sort field + order. Format: `{field} {order}`. Fields: `timestamp`, `id`, `epoch`. Orders: `asc`, `desc`. |
| `max` | int | 300 | Maximum results to return. Valid range: 1–999. Values ≥ 1000 or ≤ 0 silently revert to the default (300). |

**Response:** `200 OK`

```json
{
  "found": 42,
  "schedules": [
    {
      "scheduler": "scheduler-1",
      "schedule": {
        "id": "schedule-123",
        "epoch": 1711324800,
        "timestamp": 1711234567,
        "topic": "schedules",
        "target-topic": "my-target-topic",
        "target-key": "my-key",
        "value": "base64-or-raw-bytes"
      }
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `found` | int | Total matching results (may exceed returned count if `max` applies) |
| `schedules[].scheduler` | string | Scheduler name |
| `schedules[].schedule.id` | string | Schedule unique ID |
| `schedules[].schedule.epoch` | int64 | When the schedule fires (Unix timestamp) |
| `schedules[].schedule.timestamp` | int64 | When the schedule was created (Unix timestamp) |
| `schedules[].schedule.topic` | string | Source Kafka topic |
| `schedules[].schedule.target-topic` | string | Destination Kafka topic |
| `schedules[].schedule.target-key` | string | Message key for the target topic |
| `schedules[].schedule.value` | bytes | Message body (may be decoded via HTTP decoder) |

**Note:** Response is streamed as JSON — `found` and opening `[` are written first, then each schedule is individually encoded and comma-separated.

---

### `GET /api/scheduler/{name}/schedule/{id}`

Get a specific schedule from **cold storage**, including all versions.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `name` | string | Scheduler name |
| `id` | string | Schedule ID |

**Response:** `200 OK` — Array of schedule versions (most recent first)

```json
[
  {
    "scheduler": "scheduler-1",
    "schedule": {
      "id": "schedule-123",
      "epoch": 1711324800,
      "timestamp": 1711234567,
      "topic": "schedules",
      "target-topic": "my-target-topic",
      "target-key": "my-key",
      "value": "..."
    }
  }
]
```

**Response:** `404 Not Found` — if schedule does not exist

```json
null
```

---

### `GET /api/live/scheduler/{name}/schedules`

Search schedules currently **live** in the scheduler instance (real-time via REST).

Same query parameters and response format as `/api/scheduler/{name}/schedules`.

**Note:** Live data is fetched in real-time from the scheduler's `/schedules` REST endpoint. No persistent index — filtering and sorting are done in-memory.

---

### `GET /api/live/scheduler/{name}/schedule/{id}`

Get a specific schedule from the **live** scheduler instance.

Same path parameters and response format as `/api/scheduler/{name}/schedule/{id}`.

---

### `GET /api/history/scheduler/{name}/schedules`

Search **historical** (triggered/expired) schedules.

Same query parameters and response format as `/api/scheduler/{name}/schedules`.

---

### `GET /api/history/scheduler/{name}/schedule/{id}`

Get a specific **historical** schedule.

Same path parameters and response format as `/api/scheduler/{name}/schedule/{id}`.

---

## Internal Endpoints (not proxied to `/api`)

These are bound to separate ports and are not part of the main API:

### `GET /metrics` (port 9001)

Prometheus metrics endpoint. Exposes default Go runtime metrics (via `promhttp.Handler()` and `go-runtime-metrics`). No custom application counters are registered.

### `GET /debug/pprof/` (port 6060, localhost only)

Go pprof profiling endpoints. **Disabled by default** — enabled via `PPROF_ENABLED=true` env var or `SIGUSR1` signal.

---

## Configuration

The client loads API URLs from `/configuration.json`:

```json
{
  "stats": "/api/stats",
  "schedulers": "/api/schedulers",
  "schedules": "/api/scheduler/{name}/schedules",
  "schedule-detail": "/api/scheduler/{name}/schedule/{id}",
  "live-schedules": "/api/live/scheduler/{name}/schedules",
  "live-schedule-detail": "/api/live/scheduler/{name}/schedule/{id}",
  "history-schedules": "/api/history/scheduler/{name}/schedules",
  "history-schedule-detail": "/api/history/scheduler/{name}/schedule/{id}"
}
```

---

## Error Response

All error responses use the same format:

```json
{
  "error": "description of the error"
}
```

HTTP status code: `500 Internal Server Error`

---

## Route Map

```
GET /health                                             → health()
GET /api/stats                                          → stats()
GET /api/schedulers                                     → listSchedulers()
GET /api/scheduler/{name}/schedules?...                 → searchSchedules(coldDB)
GET /api/scheduler/{name}/schedule/{id}                 → getSchedule(coldDB)
GET /api/live/scheduler/{name}/schedules?...            → searchSchedules(liveDB)
GET /api/live/scheduler/{name}/schedule/{id}            → getSchedule(liveDB)
GET /api/history/scheduler/{name}/schedules?...         → searchSchedules(historyDB)
GET /api/history/scheduler/{name}/schedule/{id}         → getSchedule(historyDB)
```

SPA mode (`API_SERVER_ONLY=false`, default):
- `/api/*` → Go API router (above)
- `/*` → Static file server (`STATIC_FILES_DIR`, default: `../client/build`)
- SPA fallback: any non-file path → `index.html`

API-only mode (`API_SERVER_ONLY=true`):
- All routes above are served without `/api` prefix (e.g., `/stats`, `/schedulers`)


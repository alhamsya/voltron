# <img align="right" src="https://avatars.githubusercontent.com/u/56905970?s=60&v=4" alt="voltron" title="voltron" /> voltron

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alhamsya/voltron)
[![Sourcegraph](https://sourcegraph.com/github.com/alhamsya/voltron/-/badge.svg)](https://sourcegraph.com/github.com/alhamsya/voltron?badge)
[![Documentation](https://godoc.org/github.com/alhamsya/voltron?status.svg)](https://godoc.org/github.com/alhamsya/voltron)
[![codecov](https://codecov.io/gh/alhamsya/voltron/graph/badge.svg?token=PIN65DKRGQ)](https://codecov.io/gh/alhamsya/voltron)
[![Go Report Card](https://goreportcard.com/badge/github.com/alhamsya/voltron)](https://goreportcard.com/report/github.com/alhamsya/voltron)
[![License](https://img.shields.io/github/license/alhamsya/voltron?color=blue)](https://raw.githubusercontent.com/alhamsya/voltron/master/LICENSE)

## ⚡️ Dashboard
- https://power-meter-dashboard.up.railway.app/

## ⚡️ Meter Value Simulator
- https://v0-meter-value-simulator.vercel.app/

## 👀 Architecture
- architecture repository: Hexagonal Architecture
- web framework: fiber (https://gofiber.io/)

## 🎯 Structure
- `cmd`: directory for main entry points or commands of the application
- `internal`: directory for containing application code that should not exposed to external packages
- `core`: directory that contains the central business logic of the application
    - `domain`: directory that contains domain models/entities representing core business concepts
    - `port`: directory that contains defined interfaces or contracts that adapters must follow
    - `service`: directory that contains the business logic or services of the application
- `pkg`: shared managing to support service and utilities
```
cmd/
└── rest
internal/
├── adapter/
│   ├── postgresql
│   ├── handler/
│   │   └── rest
│   └── redis
├── core/
│   ├── domain/
│   │   ├── config
│   │   ├── constant
│   │   ├── request
│   │   └── response
│   ├── port/
│   │   ├── repository
│   │   └── service
│   └── service/
│       └── meter
└── mock/
    ├── repository
    └── service
pkg/
├── manager/
│   ├── config
│   ├── graceful
│   ├── protocol
│   ├── response
│   └── xhttp
└── util
```

## ⚡️ Prerequisites
1. **Go**: version 1.25.5 or higher is required
2. **Make**: running Makefile commands
3. **Docker**: running infrastructure

## ⚙️ Installation and Setup
1. clone this repository:
    ```bash
    git clone https://github.com/alhamsya/voltron.git
    cd voltron
    ```
2. download dependencies
    ```bash
    go mod download
    ```
3. first setup in local environment
    ```bash
    make start-local
    ```

## ⚙️ Running Tests
1. start service in local environment
    ```bash
    go run ./cmd/main.go run rest
    ```

## ⚡️ Mock Documentation
by default host mock: `localhost:3000`

## ⚡️ API Documentation
### Dashboard Endpoint
**Endpoint**:
- `GET /v1/api/power/latest`
- `GET /v1/api/power/time-series`
- `GET /v1/api/power/daily-usage`

### Meter Value Simulator
**Endpoint**: 
- `POST /v1/api/power/meter`

**Request Body**:
```json
{
    "PM": [
        {
            "date": "28 10:54:15/12/2025",
            "data": "[233.664283]",
            "name": "Volts"
        },
        {
            "date": "28 10:54:15/12/2025",
            "data": "[5.977841]",
            "name": "Current"
        },
        {
            "date": "28 10:54:15/12/2025",
            "data": "[1396.807870]",
            "name": "Active_Power"
        },
        {
            "date": "28 10:54:15/12/2025",
            "data": "[1.007820]",
            "name": "Total_Import_kWh"
        }
    ]
}
```

**Response** `OK (200)`:
```json
{
    "data": {
        "data": null
    },
    "message": "success meter reading successfully"
}
```

**Response** `Bad request (400)`:
```json
{
    "data": null,
    "message": "please check request date",
    "error": {
        "Layout": "02/01/2006 15:04:05",
        "Value": "",
        "LayoutElem": "02",
        "ValueElem": "",
        "Message": ""
    }
}
```

**Response** `Internal Server Error (500)`:
```json
{
    "data": null,
    "message": "please try again"
}
```
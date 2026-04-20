# 🐹 Gopher-Watch: Cloud-Native Monitoring Engine --

![CI Status](https://github.com/aditya2319/gopher-watch-fork/actions/workflows/go-ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)

Gopher-Watch is a high-performance, concurrent service monitoring engine designed for reliability and zero-noise alerting. Built in Go, it utilizes a JSON-based configuration schema (targets.json) to monitor distributed targets (APIs, Web Servers, Microservices) and manages the full incident lifecycle.



##  Core Features

* **Concurrency-First Architecture**: Utilizes Go routines and `sync.WaitGroup` to probe multiple targets in parallel without performance degradation.
* **Intelligent Alert Suppression**: Implements a failure-streak algorithm (3-cycle threshold) to prevent alert fatigue from transient network blips.
* **Incident Lifecycle Management**: Automatically detects service recovery and notifies Slack with total downtime duration.
* **Observability Ready**: Exposes a `/metrics` endpoint in Prometheus text format for integration with Grafana.
* **Hot-Reloadable**: Reloads `targets.json` configuration on every cycle, allowing for real-time monitoring updates without service downtime.

##  Roadmap & Upcoming Features
* **Full CD Pipeline**: Automated deployment to AWS/DigitalOcean via GitHub Actions.
* **Dynamic Configuration**: Support for remote config fetching (S3 or Consul).
* **Distributed Probing**: Running multiple "Gopher-Agents" across different geographic regions to detect regional outages.
* **Grafana Dashboards**: Pre-configured JSON dashboards for instant observability.

##  Technical Architecture

* **Registry**: A thread-safe singleton using `sync.RWMutex` to manage global state across concurrent probes.
* **Prober**: A resilient HTTP engine with configurable timeouts, retries, and status/body validation.
* **Notifier**: An asynchronous Slack integration using webhooks and environment-based secret management.



##  Getting Started

### Prerequisites
* Go 1.22+
* Slack Webhook URL (optional)

### Installation & Run
1. **Clone the repository**
   ```bash
   git clone [https://github.com/aditya2319/gopher-watch-fork](https://github.com/aditya2319/gopher-watch-fork.git)
   cd gopher-watch-
   ```

2. **Create a .env file in the root:**
  ```bash
SLACK_WEBHOOK_URL=[https://hooks.slack.com/services/your/webhook/url](https://hooks.slack.com/services/your/webhook/url)
```

3. **Launch the Engine**
   ```bash
   make # Runs pre-flight checks, builds, and starts monitoring
   ```
## Testing & Quality Assurance
This project maintains a high bar for reliability, utilizing the Go race detector and automated CI

    Run Unit Tests: make test

    Race Detection: All tests are executed with the -race flag to ensure thread safety.

    CI: GitHub Actions automatically verifies every push to main and all Pull Requests.
```bash
# Example Output
PASS
coverage: 75.0% of statements
ok      [github.com/aditya2319/gopher-watch-/internal/monitor](https://github.com/aditya2319/gopher-watch-/internal/monitor)  1.015s
```

## Monitoring Logs & Stats

Utilize the built-in Makefile tools to analyze engine health in real-time:

    make stats: View a summary of total passed vs. failed probes.

    make errors: Instantly filter and view JSON-formatted error logs.

    make summary: View the last 20 probe cycles with timestamps.

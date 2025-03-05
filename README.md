# Prometheus Metrics Example

这是一个使用 Prometheus 客户端库（`client_golang`）实现自定义指标收集和暴露的示例项目。项目展示了如何定义自定义指标、收集系统资源使用情况（如内存、CPU 负载、Goroutine 数量等），并通过 HTTP 端点暴露这些指标。

## 功能

1. **自定义指标收集器**：
    - 收集系统内存使用情况（堆内存、空闲内存等）。
    - 收集系统 CPU 负载（基于 Goroutine 数量与 CPU 核心数的比例）。
    - 收集当前运行的 Goroutine 数量。
    - 记录最后处理请求的时间戳。

2. **预定义指标**：
    - 请求总数计数器（`app_requests_total`）。
    - 当前并发请求数仪表盘（`app_concurrent_requests`）。
    - 请求延迟直方图（`app_request_duration_seconds`）。

3. **模拟请求处理**：
    - 模拟多个并发请求，记录请求的延迟和状态。

4. **Prometheus 指标暴露**：
    - 通过 `/metrics` 端点暴露所有指标，供 Prometheus 抓取。

## 指标说明

以下是项目暴露的主要指标：

| 指标名称                          | 类型      | 描述                                       |
|-----------------------------------|-----------|--------------------------------------------|
| `system_load_factor`              | Gauge     | 系统 CPU 负载因子（Goroutine 数量/CPU 核心数） |
| `system_memory_usage_bytes`       | Gauge     | 系统内存使用情况（堆内存、空闲内存等）       |
| `last_processed_timestamp`        | Gauge     | 最后处理请求的时间戳                        |
| `system_goroutines_count`         | Gauge     | 当前运行的 Goroutine 数量                   |
| `app_requests_total`              | Counter   | 请求总数，按方法、端点和状态分类             |
| `app_concurrent_requests`         | Gauge     | 当前并发请求数                              |
| `app_request_duration_seconds`    | Histogram | 请求延迟分布                                |


## 许可证

本项目基于 [MIT 许可证](LICENSE) 开源
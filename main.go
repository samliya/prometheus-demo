package main

import (
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServiceMetricsCollector 自定义的收集器
// 需要实现Collector 的两个接口 Describe和Collect
// 向Collect的ch传递数据时，需要指明使用哪一个desc、metrics的类型是什么以及具体的值，还可以携带的label对应的值
// 具体的采集频率由prometheus.yaml 的scrape_interval 决定
type ServiceMetricsCollector struct {
	systemLoadDesc    *prometheus.Desc
	memoryUsageDesc   *prometheus.Desc
	lastProcessedDesc *prometheus.Desc
	goroutinesDesc    *prometheus.Desc // 新增：goroutine数量指标
}

func NewServiceMetricsCollector() *ServiceMetricsCollector {
	return &ServiceMetricsCollector{
		systemLoadDesc: prometheus.NewDesc(
			"system_load_factor",
			"Current system load factor",
			[]string{"component"},
			nil,
		),
		memoryUsageDesc: prometheus.NewDesc(
			"system_memory_usage_bytes",
			"Current memory usage in bytes",
			[]string{"type"},
			nil,
		),
		lastProcessedDesc: prometheus.NewDesc(
			"last_processed_timestamp",
			"Timestamp of the last processed request",
			nil,
			nil,
		),
		goroutinesDesc: prometheus.NewDesc(
			"system_goroutines_count",
			"Number of goroutines currently running",
			nil,
			nil,
		),
	}
}

func (c *ServiceMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.systemLoadDesc
	ch <- c.memoryUsageDesc
	ch <- c.lastProcessedDesc
	ch <- c.goroutinesDesc
}

func (c *ServiceMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	// 获取内存统计信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 收集堆内存使用情况
	ch <- prometheus.MustNewConstMetric(
		c.memoryUsageDesc,
		prometheus.GaugeValue,
		float64(memStats.HeapAlloc), // 当前堆内存使用量
		"heap_alloc",
	)
	ch <- prometheus.MustNewConstMetric(
		c.memoryUsageDesc,
		prometheus.GaugeValue,
		float64(memStats.HeapSys), // 从系统获取的堆内存
		"heap_sys",
	)
	ch <- prometheus.MustNewConstMetric(
		c.memoryUsageDesc,
		prometheus.GaugeValue,
		float64(memStats.HeapIdle), // 空闲堆内存
		"heap_idle",
	)

	// 收集系统CPU使用率（这里使用runtime.NumCPU()作为基准负载）
	ch <- prometheus.MustNewConstMetric(
		c.systemLoadDesc,
		prometheus.GaugeValue,
		float64(runtime.NumGoroutine())/float64(runtime.NumCPU()), // 每个CPU的平均goroutine数
		"cpu_load",
	)

	// 收集goroutine数量
	ch <- prometheus.MustNewConstMetric(
		c.goroutinesDesc,
		prometheus.GaugeValue,
		float64(runtime.NumGoroutine()),
	)

	// 记录最后处理时间
	ch <- prometheus.MustNewConstMetric(
		c.lastProcessedDesc,
		prometheus.GaugeValue,
		float64(time.Now().Unix()),
	)
}

// 预定义的指标
var (
	// 计数器：记录请求总数
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of requests processed",
		},
		[]string{"method", "endpoint", "status"},
	)

	// 仪表盘：记录当前并发请求数
	concurrentRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_concurrent_requests",
			Help: "Number of concurrent requests being processed",
		},
	)

	// 直方图：请求延迟
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_request_duration_seconds",
			Help:    "Request latency in seconds",
			Buckets: prometheus.DefBuckets, // 默认延迟桶
		},
		[]string{"method", "endpoint"},
	)
)

// 模拟服务处理请求
func handleRequest(method, endpoint string) {
	// 增加并发请求计数
	concurrentRequests.Inc()
	defer concurrentRequests.Dec()

	// 记录请求开始时间
	start := time.Now()

	// 模拟业务处理
	time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(900)))

	// 计算请求耗时
	duration := time.Since(start)

	// 记录请求指标
	requestCounter.WithLabelValues(method, endpoint, "200").Inc()
	requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func main() {
	customCollector := NewServiceMetricsCollector()

	reg := prometheus.NewRegistry()

	// 注册默认指标和自定义收集器
	reg.MustRegister(customCollector)
	reg.MustRegister(requestCounter)
	reg.MustRegister(concurrentRequests)
	reg.MustRegister(requestDuration)

	// 模拟请求处理的并发goroutine
	go func() {
		for {
			go handleRequest("GET", "/users")
			go handleRequest("POST", "/login")
			go handleRequest("GET", "/metrics")
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Println("Starting metrics server on :2115")
	log.Fatal(http.ListenAndServe(":2115", nil))
}

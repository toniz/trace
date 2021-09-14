/*
 * Create By Xinwenjia 2021-09-06
 */

package trace

import (
    "net/http"
    "time"
    "context"
    "errors"
    "io/ioutil"
    "bytes"
    "strconv"

    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/encoding/gzip"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    //"go.opentelemetry.io/otel/baggage"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
    "go.opentelemetry.io/otel/semconv/v1.4.0"

    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
    slsProjectHeader         = "x-sls-otel-project"
    slsInstanceIDHeader      = "x-sls-otel-instance-id"
    slsAccessKeyIDHeader     = "x-sls-otel-ak-id"
    slsAccessKeySecretHeader = "x-sls-otel-ak-secret"
    slsSecurityTokenHeader   = "x-sls-otel-token"
)

var (
    OtelSpanKindUnspecified = strconv.Itoa(int(trace.SpanKindUnspecified))
    OtelSpanKindInternal = strconv.Itoa(int(trace.SpanKindInternal))
    OtelSpanKindServer = strconv.Itoa(int(trace.SpanKindServer))
    OtelSpanKindClient = strconv.Itoa(int(trace.SpanKindClient))
    OtelSpanKindProducer = strconv.Itoa(int(trace.SpanKindProducer))
    OtelSpanKindConsumer = strconv.Itoa(int(trace.SpanKindConsumer))
)

var (
    TraceExportNotSupport error = errors.New("TraceExportNotSupport")
    TraceExportAlreadySet error = errors.New("TraceExportAlreadySet")
    TraceConfigCheckFailed error = errors.New("TraceConfigCheckFailed")
    TraceNotWork error = errors.New("TraceNotWork")
    TraceEndSpanFailed error = errors.New("TraceEndSpanFailed")
    TraceHttpCodeError error = errors.New("TraceHttpCodeError")
)


type Otel struct {
    grpcexp *otlptrace.Exporter
    stdexp *stdouttrace.Exporter
    tp *sdktrace.TracerProvider
}

func init() {
    s := new(Otel)
    SetTracer(s)
}

// Set Http Export
func (c *Otel) SetHttpExport(ctx context.Context, filename, serviceName, version string) error {
    return TraceExportNotSupport
}

// Set Grpc Export
func (c *Otel) SetGrpcExport(ctx context.Context, filename, serviceName, version string) error {
    if c.IsWork() == nil {
        return TraceExportAlreadySet
    }

    params := make(map[string]string)
    err := LoadFile(filename, &params)
    if err != nil {
        return err
    }

    traceEndpoint := params["trace_endpoint"]
    projectId := params["project_id"]
    serviceId := params["service_id"]
    accessKeyId := params["access_key_id"]
    accessKeySecret := params["access_key_secret"]
    if len(traceEndpoint) == 0 || len(projectId) == 0 ||len(serviceId) == 0 || len(accessKeyId) == 0 || len(accessKeySecret) == 0 {
        return TraceConfigCheckFailed
    }

    headers := map[string]string {
        slsProjectHeader:         projectId,
        slsInstanceIDHeader:      serviceId,
        slsAccessKeyIDHeader:     accessKeyId,
        slsAccessKeySecretHeader: accessKeySecret,
    }

    secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
    if params["insecure"] == "1" {
        secureOption = otlptracegrpc.WithInsecure()
    }

    opts := []otlptracegrpc.Option{
        secureOption,
        otlptracegrpc.WithEndpoint(traceEndpoint),
        otlptracegrpc.WithHeaders(headers),
        otlptracegrpc.WithCompressor(gzip.Name),
        otlptracegrpc.WithReconnectionPeriod(50 * time.Millisecond),
    }

    var otlpErr error
    client := otlptracegrpc.NewClient(opts...)
    c.grpcexp, otlpErr = otlptrace.New(ctx, client)
    if otlpErr != nil {
        return otlpErr
    }

    err = c.SetGlobalProvider(ctx, serviceName, version)
    if err != nil {
        return err
    }

    return nil
}

// Set Default Exprot: stdout
func (c *Otel) SetDefaultExport(ctx context.Context, serviceName, version string) error {
    if c.IsWork() == nil {
        return TraceExportAlreadySet
    }

    var expErr error
    c.stdexp, expErr = stdouttrace.New(stdouttrace.WithPrettyPrint())
    if expErr != nil {
        return expErr
    }

    err := c.SetGlobalProvider(ctx, serviceName, version)
    if err != nil {
        return err
    }
    return nil
}

// Set Glocal Trace Provider
func (c *Otel) SetGlobalProvider(ctx context.Context, serviceName, version string) error {
    myResource, rcErr := resource.New(ctx,
        resource.WithProcessPID(),
        resource.WithTelemetrySDK(),
        resource.WithHost(),
        resource.WithOSType(),
        resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
        resource.WithAttributes(semconv.ServiceVersionKey.String(version)),
    )

    if rcErr != nil {
        return rcErr
    }

    var bsp sdktrace.SpanProcessor
    if c.grpcexp != nil {
        bsp = sdktrace.NewBatchSpanProcessor(c.grpcexp)
    } else {
        bsp = sdktrace.NewBatchSpanProcessor(c.stdexp)
    }

    c.tp = sdktrace.NewTracerProvider (
        sdktrace.WithSpanProcessor(bsp),
        sdktrace.WithResource(myResource),
    )

    otel.SetTracerProvider(c.tp)
    propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
    otel.SetTextMapPropagator(propagator)

    return nil
}

func (c *Otel) Close(ctx context.Context) error {
    if c.tp != nil {
        c.tp.Shutdown(ctx)
    }

    if c.grpcexp != nil {
        c.grpcexp.Shutdown(ctx)
    }

    if c.stdexp != nil {
        c.stdexp.Shutdown(ctx)
    }
    return nil
}

func (c *Otel) IsWork() error {
    if c.tp == nil {
        return TraceNotWork
    }

    return nil
}

// Start A child Span
// 0: SpanKindUnspecified
// 1: SpanKindInternal
// 2: SpanKindServer
// 3: SpanKindClient
// 4: SpanKindProducer
// 5: SpanKindConsumer
func (c *Otel) NewSpan(ctx context.Context, name string, kind string) (context.Context, error) {
    if c.IsWork() != nil {
        return ctx, TraceNotWork
    }

    // Label Span Level
    level, ok := ctx.Value(ContextSpanLevelKey).(int)
    if !ok {
        level = 0
    }
    ctxl := context.WithValue(ctx, ContextSpanLevelKey, level+1)

    iKind, _ := strconv.Atoi(kind)
    tracer := c.tp.Tracer(semconv.SchemaURL, trace.WithInstrumentationVersion(otel.Version()))
    ctxc, _ := tracer.Start(ctxl, name, trace.WithSpanKind(trace.SpanKind(iKind)))
    return ctxc, nil
}

// Call Span End
func (c *Otel) EndSpan(ctx context.Context) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    span := trace.SpanFromContext(ctx)
    if span == nil {
        return TraceEndSpanFailed
    }
    span.End()
    return nil
}

// Adding An Attribute To A Span
func (c *Otel) AddSpanAttribute(ctx context.Context, params map[string]string) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    var attr []attribute.KeyValue
    for k, v := range params {
        attr = append(attr, attribute.String(k, v))
    }

    span := trace.SpanFromContext(ctx)
    span.SetAttributes(attr...)
    return nil
}

// Adding An Event To A Span
func (c *Otel) AddSpanEvent(ctx context.Context, event string, params map[string]string) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    var attr []attribute.KeyValue
    for k, v := range params {
        attr = append(attr, attribute.String(k, v))
    }

    span := trace.SpanFromContext(ctx)
    span.AddEvent(event, trace.WithAttributes(attr...))
    return nil
}

// Set Span Status OK
func (c *Otel) SetSpanOK(ctx context.Context, message string) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    if len(message) == 0 {
        message = "Success"
    }

    span := trace.SpanFromContext(ctx)
    span.SetStatus(codes.Ok, message)
    return nil
}

// Set Span Status Error
func (c *Otel) SetSpanError(ctx context.Context, err error) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    span := trace.SpanFromContext(ctx)
    span.SetStatus(codes.Error, err.Error())
    errTime := time.Now()
    span.RecordError(err, trace.WithTimestamp(errTime))

    return nil
}

// Http Handle Hook
func (c *Otel) NewHandler(handler http.Handler, name string) http.Handler {
    if c.IsWork() != nil {
        return handler
    }

    opt := otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents)
    return otelhttp.NewHandler(handler, name, opt)
}

// Http Do Function
func (c *Otel) HttpDo(ctx context.Context, req *http.Request, name string, options ...interface{}) (*http.Response, error) {
    client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Call Http Get
func (c *Otel) HttpGet(ctx context.Context, url string, name string) (string, error) {
    if c.IsWork() != nil {
        return "", TraceNotWork
    }

    resp, errOH := otelhttp.Get(ctx, url)
    if errOH != nil {
        return "", errOH
    }

    defer resp.Body.Close()
    resBody, errIO := ioutil.ReadAll(resp.Body)
    if errIO != nil {
        return "", errIO
    }

    if resp.StatusCode > 399 {
        return "", TraceHttpCodeError
    }

    return string(resBody), nil
}

// Call Http Post
func (c *Otel) HttpPost(ctx context.Context, url, contentType, body, name string) (string, error) {
    if c.IsWork() != nil {
        return "", TraceNotWork
    }

    reqBody := bytes.NewBuffer([]byte(body))
    resp, errOH := otelhttp.Post(ctx, url, contentType, reqBody)
    if errOH != nil {
        return "", errOH
    }

    defer resp.Body.Close()
    resBody, errIO := ioutil.ReadAll(resp.Body)
    if errIO != nil {
        return "", errIO
    }

    if resp.StatusCode > 399 {
        return "", TraceHttpCodeError
    }

    return string(resBody), nil 
}


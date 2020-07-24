// 2019. Mahesh Voleti (mvoleti@vmware.com)

package internal

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	otrext "github.com/opentracing/opentracing-go/ext"
        wfreporter "github.com/wavefronthq/wavefront-opentracing-sdk-go/reporter"
        wftracer "github.com/wavefronthq/wavefront-opentracing-sdk-go/tracer"
        application "github.com/wavefronthq/wavefront-sdk-go/application"
        wavefront "github.com/wavefronthq/wavefront-sdk-go/senders"
)

/*func NewGlobalTracer(serviceName string) io.Closer {

	config, enverr := jaegercfg.FromEnv()
	if enverr != nil {
		log.Println("Couldn't parse Jaeger env vars", enverr.Error())
		os.Exit(1)
	}

	config.ServiceName = serviceName
	config.Sampler.Type = jaeger.SamplerTypeConst
	config.Sampler.Param = 1
	config.Reporter.LogSpans = true

	if GlobalConfig.JaegerHostPort != "" {
		config.Reporter.LocalAgentHostPort = GlobalConfig.JaegerHostPort
	}

	log.Printf("Connecting to Jaeger @ %s\n", config.Reporter.LocalAgentHostPort)

	closer, err := config.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(jmetrics.NullFactory),
	)
	if err != nil {
		log.Println("Couldn't initialize tracer", err.Error())
		os.Exit(1)
	}

	return closer
}*/

func NewDirectGlobalTracer(serviceName string, cluster string, token string, batchSize int, maxBufferSize int, flushInterval int, applicationName string) io.Closer {

	config := &wavefront.DirectConfiguration{
		Server: cluster,
		Token: token,
		BatchSize: batchSize,
		MaxBufferSize: maxBufferSize,
		FlushIntervalSeconds: flushInterval,
	}
	sender, err := wavefront.NewDirectSender(config)
	if err != nil {
		log.Printf("Failed to create Wavefront Sender: %s\n", err.Error())
		os.Exit(1)
	}

	appTags := application.New(applicationName, serviceName)

	directrep := wfreporter.New(sender, appTags)
	consolerep := wfreporter.NewConsoleSpanReporter(serviceName)

	reporter := wfreporter.NewCompositeSpanReporter(directrep, consolerep)

	tracer := wftracer.New(reporter)

	opentracing.SetGlobalTracer(tracer)

	return ioutil.NopCloser(nil)

}

func NewProxyGlobalTracer(serviceName string, proxyIp string, tracingPort int, metricsPort int, distributionPort int, flushInterval int, applicationName string) io.Closer {

        config := &wavefront.ProxyConfiguration{
                Host:        proxyIp,
                TracingPort: tracingPort,
		MetricsPort: metricsPort,
		DistributionPort: distributionPort,
		FlushIntervalSeconds: flushInterval,
        }

        sender, err := wavefront.NewProxySender(config)
        if err != nil {
                log.Printf("Failed to create Wavefront Sender: %s\n", err.Error())
                os.Exit(1)
        }

	appTags := application.New(applicationName, serviceName)

        directrep := wfreporter.New(sender, appTags)
        consolerep := wfreporter.NewConsoleSpanReporter(serviceName)

        reporter := wfreporter.NewCompositeSpanReporter(directrep, consolerep)

        tracer := wftracer.New(reporter)

        opentracing.SetGlobalTracer(tracer)

        return ioutil.NopCloser(nil)

}


func NewServerSpan(req *http.Request, spanName string) opentracing.Span {
	tracer := opentracing.GlobalTracer()
	parentCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	var span opentracing.Span
	if err == nil { // has parent context
		span = tracer.StartSpan(spanName, opentracing.ChildOf(parentCtx))
	} else if err == opentracing.ErrSpanContextNotFound { // no parent
		span = tracer.StartSpan(spanName)
	} else {
		log.Printf("Error in extracting tracer context: %s", err.Error())
	}

	otrext.SpanKindRPCServer.Set(span)

	return span
}


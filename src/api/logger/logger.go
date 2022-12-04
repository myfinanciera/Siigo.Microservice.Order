package logger

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/Siigo.Core.Logs.Golang.git/easy"
	"github.com/ic2hrmk/promtail"
	"github.com/sirupsen/logrus"
	"os"
	"siigo.com/order/src/api/config"
)

const BusinessFieldKey = "business"
const BusinessFieldValue = "yes"

type BusinessHook struct {
	client promtail.Client
}

// NewLogrus Create a new instance of logrus with custom hooks
func NewLogrus(l *logrus.Logger, config *config.Configuration) *logrus.Logger {

	l.SetReportCaller(config.Log.ReportCaller)
	l.SetFormatter(&easy.Formatter{})
	l.SetOutput(os.Stdout)

	l.AddHook(NewBusinessHook(config))

	levelLog, er := logrus.ParseLevel(config.Log.Level)
	if er != nil {
		levelLog = logrus.InfoLevel
	}
	l.SetLevel(levelLog)

	return l
}

func NewBusinessHook(config *config.Configuration) *BusinessHook {

	promtailClient, err := promtail.NewJSONv1Client(config.BusinessLogger.LokiUrl, config.BusinessLogger.DefaultLabels)

	if err != nil {
		panic(err)
	}

	return &BusinessHook{
		client: promtailClient,
	}

}

func (hook *BusinessHook) Fire(entry *logrus.Entry) error {

	if _, exist := entry.Data[BusinessFieldKey]; !exist {
		return nil
	}

	go hook.client.LogfWithLabels(promtail.Info, FieldsToMap(entry.Data), entry.Message)
	return nil
}

func (hook *BusinessHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func FieldsToMap(fields logrus.Fields) map[string]string {
	mapFields := make(map[string]string)
	for key, element := range fields {
		mapFields[key] = element.(string)
	}
	return mapFields
}

func WithBusinessFields(fields logrus.Fields) logrus.Fields {
	fields[BusinessFieldKey] = BusinessFieldValue
	return fields
}

type LoggerBuilder struct {
	Fields map[string]interface{}
}

func NewBusinessLogger() LoggerBuilder {
	c := LoggerBuilder{Fields: map[string]interface{}{}}
	c.Fields[BusinessFieldKey] = BusinessFieldValue
	return c
}

func (receiver LoggerBuilder) WithMetricType(metric string) LoggerBuilder {
	receiver.Fields["MetricName"] = metric
	return receiver
}

func (receiver LoggerBuilder) WithApiLayer() LoggerBuilder {
	receiver.Fields["Layer"] = "Api"
	return receiver
}

func (receiver LoggerBuilder) WithDomainLayer() LoggerBuilder {
	receiver.Fields["Layer"] = "Domain"
	return receiver
}

func (receiver LoggerBuilder) WithInfraLayer() LoggerBuilder {
	receiver.Fields["Layer"] = "Infra"
	return receiver
}

func (receiver LoggerBuilder) WithApplicationLayer() LoggerBuilder {
	receiver.Fields["Layer"] = "Application"
	return receiver
}

func (receiver LoggerBuilder) ToFields() map[string]interface{} {
	return receiver.Fields
}

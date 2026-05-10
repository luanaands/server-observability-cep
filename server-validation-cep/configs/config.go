package configs

import (
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

type Conf struct {
	MyCoreHost               string `mapstructure:"MY_CORE_HOST"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME"`
	Title                    string `mapstructure:"TITLE"`
	BackgroundColor          string `mapstructure:"BACKGROUND_COLOR"`
	ExternalCallURL          string `mapstructure:"EXTERNAL_CALL_URL"`
	ExternalCallMethod       string `mapstructure:"EXTERNAL_CALL_METHOD"`
	RequestNameOTEL          string `mapstructure:"REQUEST_NAME_OTEL"`
	OTELTracer               trace.Tracer
}

func LoadConfig() (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

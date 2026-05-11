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
	OTELTracer               trace.Tracer
}

func LoadConfig() (*Conf, error) {
	var cfg *Conf

	// Configurar viper para ler variáveis de ambiente
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Tentar ler arquivo de config se existir (opcional)
	_ = viper.ReadInConfig()

	// Bind explícito das variáveis de ambiente
	viper.BindEnv("MY_CORE_HOST", "MY_CORE_HOST")
	viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_ENDPOINT")
	viper.BindEnv("OTEL_SERVICE_NAME", "OTEL_SERVICE_NAME")
	viper.BindEnv("TITLE", "TITLE")

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

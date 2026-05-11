package configs

import (
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

type Conf struct {
	ViaCepApiHost            string `mapstructure:"VIA_CEP_API_HOST"`
	ApiWeatherKey            string `mapstructure:"API_WEATHER_KEY"`
	ApiWeatherHost           string `mapstructure:"API_WEATHER_HOST"`
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

	// Configurar viper para ler variáveis de ambiente
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Tentar ler arquivo de config se existir (opcional)
	_ = viper.ReadInConfig()

	// Bind explícito das variáveis de ambiente
	viper.BindEnv("VIA_CEP_API_HOST", "VIA_CEP_API_HOST")
	viper.BindEnv("API_WEATHER_KEY", "API_WEATHER_KEY")
	viper.BindEnv("API_WEATHER_HOST", "API_WEATHER_HOST")
	viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_ENDPOINT")
	viper.BindEnv("OTEL_SERVICE_NAME", "OTEL_SERVICE_NAME")
	viper.BindEnv("TITLE", "TITLE")
	viper.BindEnv("BACKGROUND_COLOR", "BACKGROUND_COLOR")
	viper.BindEnv("EXTERNAL_CALL_URL", "EXTERNAL_CALL_URL")
	viper.BindEnv("EXTERNAL_CALL_METHOD", "EXTERNAL_CALL_METHOD")
	viper.BindEnv("REQUEST_NAME_OTEL", "REQUEST_NAME_OTEL")

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

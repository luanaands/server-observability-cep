# Server Observability - CEP e Clima com Tracing

Projeto em Go com duas APIs que consultam CEP e clima em tempo real, com observabilidade distribuida usando OpenTelemetry + Zipkin.

## Descricao

Arquitetura do fluxo:

1. `server-validation-cep` (Service A) recebe o CEP em `POST /cep`.
2. Service A chama o `server-core-cep` (Service B).
3. Service B consulta ViaCEP e WeatherAPI e retorna cidade e temperaturas.
4. A entrada HTTP dos dois servicos esta instrumentada com `otelhttp.NewHandler`, preservando o contexto de trace entre Service A e Service B.
5. Os spans sao enviados para o OpenTelemetry Collector via OTLP gRPC (`4317`).
6. O Collector exporta os traces para o Zipkin.

## Componentes existentes

- `server-validation-cep`:
  - Porta `8082`
  - Endpoint principal: `POST /cep`
  - Valida CEP e orquestra chamada para o Service B
  - Entrada HTTP instrumentada com `otelhttp.NewHandler` (span `service-a.request /cep`)
- `server-core-cep`:
  - Porta `8081`
  - Endpoint principal: `POST /cep`
  - Consulta ViaCEP e WeatherAPI
  - Entrada HTTP instrumentada com `otelhttp.NewHandler` (span `service-b.request /cep`)
  - Spans de negocio para chamadas externas: `viacep.lookup` e `weatherapi.lookup`
- `otel-collector`:
  - Porta `4317` (OTLP gRPC)
  - Config em `./.docker/otel-collector-config.yaml`
  - Receiver configurado para `0.0.0.0:4317`
- `zipkin`:
  - UI em `http://localhost:9411`
  - Recebe spans do Collector

## Pre-requisitos

- Docker + Docker Compose
- Go `1.25+` (apenas se for executar local sem Docker)

## Como executar (recomendado: Docker Compose)

No diretorio raiz do projeto:

```bash
docker compose up --build -d
```

Verificar status dos containers:

```bash
docker compose ps
```

Parar o ambiente:

```bash
docker compose down
```

## Como executar sem Docker (opcional)

### 1) Service B (`server-core-cep`)

```bash
cd server-core-cep
go mod tidy
```

Variaveis minimas:

```env
VIA_CEP_API_HOST=http://viacep.com.br/ws/
API_WEATHER_HOST=http://api.weatherapi.com/v1/current.json
API_WEATHER_KEY=<SUA_CHAVE>
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=server-core-cep
TITLE=Service B
```

No `docker-compose.yaml`, esse mesmo servico usa `TITLE=service-b` por padrao.

No `docker-compose.yaml`, o valor padrao de `TITLE` para esse servico esta como `service-b`.

Executar:

```bash
go run cmd/server/main.go
```

### 2) Service A (`server-validation-cep`)

```bash
cd server-validation-cep
go mod tidy
```

Variaveis minimas:

```env
MY_CORE_HOST=http://localhost:8081/cep
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=server-validation-cep
TITLE=Service A
```

No `docker-compose.yaml`, esse mesmo servico usa por padrao:

```env
MY_CORE_HOST=http://server-core-cep:8081/cep
TITLE=service-a
```

No `docker-compose.yaml`, os valores padrao desse servico sao:

```env
MY_CORE_HOST=http://server-core-cep:8081/cep
TITLE=service-a
```

Executar:

```bash
go run cmd/server/main.go
```

## Como testar os endpoints

Fluxo completo (recomendado - entra pelo Service A):

```bash
curl -X POST http://localhost:8082/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"41720000"}'
```

Teste direto no Service B:

```bash
curl -X POST http://localhost:8081/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01001201"}'
```

Arquivos HTTP de apoio:

- `server-validation-cep/test/cep.http`
- `server-core-cep/test/cep.http`

## Swagger

- Service A: `http://localhost:8082/docs/index.html`
- Service B: `http://localhost:8081/docs/index.html`

## Como verificar os traces

1. Suba os servicos com `docker compose up --build -d`.
2. Gere trafego com uma requisicao `POST /cep` no Service A.
3. Abra o Zipkin em `http://localhost:9411`.
4. Filtre por servico:
   - `server-validation-cep`
   - `server-core-cep`
5. Abra um trace e valide:
   - Span de entrada do Service A: `service-a.request /cep`
   - Span de entrada do Service B: `service-b.request /cep`
   - Spans filhos de negocio no Service B:
      - `viacep.lookup`
      - `weatherapi.lookup`
   - Spans HTTP automaticos (`otelhttp.NewTransport`) para:
     - chamada Service A -> Service B
     - chamada Service B -> ViaCEP
     - chamada Service B -> WeatherAPI

## Troubleshooting de traces

Se aparecer erro como `connection refused` ao exportar traces:

- Verifique se o `otel-collector` esta rodando (`docker compose ps`).
- Verifique logs do collector (`docker compose logs otel-collector`).
- Confirme em `./.docker/otel-collector-config.yaml` que o receiver gRPC esta em `0.0.0.0:4317`.
- Se o trace estiver quebrado entre Service A e B, confirme que as requisicoes entram por `POST /cep` (rota instrumentada com `otelhttp.NewHandler`) nos dois servicos.
- Recrie o collector se necessario:

```bash
docker compose up -d --force-recreate otel-collector
```

## Testes

Rodar testes por modulo:

```bash
cd server-core-cep && go test ./...
cd ../server-validation-cep && go test ./...
```

## Contato

Desenvolvido por Luana Andrade - luanaands@gmail.com

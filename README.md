# Clima por CEP — Go + Cloud Run

Recebe um CEP brasileiro, descobre a cidade via [ViaCEP](https://viacep.com.br/) e retorna a temperatura atual via [WeatherAPI](https://www.weatherapi.com/) em Celsius, Fahrenheit e Kelvin.

## URL no Cloud Run

**https://cep-weather-458435769570.us-central1.run.app**

Exemplo: `https://cep-weather-458435769570.us-central1.run.app/weather/01001000`

## API

`GET /weather/{cep}` — CEP com 8 dígitos, sem hífen.

| Cenário | Status | Corpo |
|---|---|---|
| Sucesso | 200 | `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}` |
| CEP com formato inválido | 422 | `invalid zipcode` |
| CEP não encontrado | 404 | `can not find zipcode` |

Conversões: `F = C × 1.8 + 32` e `K = C + 273` (valores arredondados em 2 casas).

## Configuração

Crie uma chave gratuita em https://www.weatherapi.com/signup.aspx e configure o `.env`:

```bash
cp .env.example .env   # preencha WEATHER_API_KEY
```

| Variável | Obrigatória | Padrão |
|---|---|---|
| `WEATHER_API_KEY` | sim | — |
| `PORT` | não | `8080` |

Variáveis já definidas no ambiente têm prioridade sobre o `.env`.

## Rodando localmente com Docker

```bash
docker compose up --build
curl localhost:8080/weather/01001000
```

## Rodando localmente sem Docker

Requer Go 1.27+. Execute na raiz do projeto, onde fica o `.env`:

```bash
go run ./cmd
curl localhost:8080/weather/01001000
```

Ou gere o binário:

```bash
go build -o bin/server ./cmd
./bin/server
```

Sem `.env`, passe as variáveis direto: `WEATHER_API_KEY=... PORT=9090 go run ./cmd`.

O servidor encerra de forma graciosa ao receber `Ctrl+C` (SIGINT) ou SIGTERM.

## Testes

Via Docker (o build do estágio `test` falha se algum teste falhar):

```bash
docker build --target test .
```

Ou direto com Go:

```bash
go test ./... -v
```

Os testes cobrem:

- conversões de temperatura;
- clientes HTTP do ViaCEP e da WeatherAPI (com servidores `httptest`);
- todos os cenários do handler (200, 404, 422 e 500);
- carregamento de configuração (porta padrão e chave obrigatória);
- roteamento e ciclo de vida do servidor (graceful shutdown e porta ocupada).

Nenhum teste acessa a internet.

## Deploy no Cloud Run

```bash
gcloud run deploy cep-weather \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=...
```

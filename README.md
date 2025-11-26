# Weather API by CEP

Sistema em Go que recebe um CEP brasileiro, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## Funcionalidades

- Recebe CEP válido de 8 dígitos
- Consulta localização via API ViaCEP
- Consulta temperatura via WeatherAPI
- Converte temperaturas para Celsius, Fahrenheit e Kelvin
- Tratamento de erros adequado

## API Endpoints

### GET /weather/{cep}

Retorna as temperaturas para o CEP informado.

#### Respostas

**Sucesso (200 OK):**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

**CEP inválido (422 Unprocessable Entity):**
```json
{
  "message": "invalid zipcode"
}
```

**CEP não encontrado (404 Not Found):**
```json
{
  "message": "can not find zipcode"
}
```

## Como executar

### Usando Docker Compose (Recomendado)

1. Clone o repositório
2. (Opcional) Configure a API key do WeatherAPI:
   ```bash
   cp .env.example .env
   # Edite o arquivo .env com sua API key
   ```
3. Execute com Docker Compose:
   ```bash
   docker-compose up --build
   ```

### Executando localmente

1. Instale as dependências:
   ```bash
   go mod tidy
   ```

2. Execute os testes:
   ```bash
   go test -v
   ```

3. Execute a aplicação:
   ```bash
   go run main.go
   ```

## Testando a API

### Exemplos de CEPs válidos:
- `01310100` - Av. Paulista, São Paulo/SP
- `20040020` - Centro, Rio de Janeiro/RJ
- `30112000` - Centro, Belo Horizonte/MG

### Testando com curl:

```bash
# CEP válido
curl http://localhost:8080/weather/01310100

# CEP inválido (formato)
curl http://localhost:8080/weather/123

# CEP não encontrado
curl http://localhost:8080/weather/99999999
```

## Deploy no Google Cloud Run

### Pré-requisitos
- Google Cloud CLI instalado e configurado
- Projeto no Google Cloud Platform
- Billing habilitado no projeto

### Passos para deploy:

1. Configure o projeto:
   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

2. Habilite as APIs necessárias:
   ```bash
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable run.googleapis.com
   ```

3. Faça o build e deploy:
   ```bash
   gcloud run deploy weather-api \
     --source . \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated \
     --set-env-vars="WEATHER_API_KEY=your_api_key_here"
   ```

4. A URL do serviço será exibida após o deploy.

## Tecnologias utilizadas

- **Go 1.21** - Linguagem de programação
- **Gorilla Mux** - Router HTTP
- **ViaCEP API** - Consulta de CEPs brasileiros
- **WeatherAPI** - Dados meteorológicos
- **Docker** - Containerização
- **Google Cloud Run** - Plataforma de deploy

## Estrutura do projeto

```
.
├── main.go              # Código principal da aplicação
├── main_test.go         # Testes automatizados
├── go.mod              # Dependências do Go
├── Dockerfile          # Configuração do Docker
├── docker-compose.yml  # Orquestração local
├── .env.example        # Exemplo de variáveis de ambiente
└── README.md           # Este arquivo
```

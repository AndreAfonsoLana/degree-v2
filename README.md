# Clima por CEP com Observabilidade (OTEL + Zipkin)

Sistema distribuído construído em Go para consultar a temperatura de uma cidade baseada no CEP, com foco em resiliência e na implementação de rastreamento distribuído (Distributed Tracing) utilizando OpenTelemetry e Zipkin.

## 🏗️ Arquitetura do Projeto

O sistema é composto por dois microsserviços principais que se comunicam entre si:
- **Serviço A (Validador de CEP):** Recebe a requisição do usuário, valida se o CEP possui exatamente 8 dígitos (sem hífens) e encaminha a solicitação para o Serviço B propagando o contexto de rastreamento.
- **Serviço B (Orquestrador de Clima):** Recebe o CEP válido, busca o nome da localidade (via API ViaCEP) e, em seguida, busca a temperatura atual dessa cidade (via WeatherAPI).

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** Go (Golang) 1.24/1.25
- **Observabilidade:** OpenTelemetry (OTEL)
- **Tracing / UI:** Zipkin
- **Infraestrutura:** Docker & Docker Compose

## 🚀 Como executar o projeto

### Pré-requisitos
- Docker e Docker Compose instalados na máquina.
- Uma chave de API válida da [WeatherAPI](https://www.weatherapi.com/).

### Passo a Passo

1. Clone o repositório para sua máquina.
2. Certifique-se de inserir sua chave da API de clima na configuração do projeto (arquivo `docker-compose.yaml` na variável de ambiente correspondente ou via arquivo `.env`).
3. Na raiz do projeto, construa as imagens e suba a infraestrutura completa utilizando o Docker:
   ```bash
   docker compose up -d --build
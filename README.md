# Magalu Challenge

## Descrição
Solução para o desafio [Magalu](./CHALLENGE.md), desenvolvido em Go com PostgreSQL. O projeto conta com testes unitarios e de integração, alem do uso de Docker.

O projeto é uma API REST que permite o agendamento de notificaçãos. Ele possui endpoints para agendar, consultar status e cancelar/excluir agendamento.
O banco escolhifo foi o [PostgreSQL](https://www.postgresql.org/) por ser um banco de dados relacional robusto, amplamente utilizado e extremamente versátil.

## Requisitos
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Git](https://git-scm.com/)

## Setup
1. Clone o repositório
2. Entre na pasta do projeto
3. Copie o arquivo `.env.example` para `.env`
    ```
    cp .env.example .env
    ```
4. Execute o projeto
   ```bash
   docker-compose up -d
   ```
   > Se você estiver executando pela primeira vez o container pode demorar um pouco para subir. Execute `docker-compose logs -f` para acompanhar o processo.
   
   > As migrações do banco de dados serão executadas automaticamente na inicialização do container de banco de dados.

## Testes
Para rodar os testes, execute o comando:
```bash
docker-compose run --rm api go test ./...
```

Para rodar os testes com cobertura de código, execute o comando:
```bash
docker-compose run --rm api go test -coverprofile=coverage.out ./...
```

Para visualizar o relatório de cobertura, execute o comando:
```bash
go tool cover -html=coverage.out
```
> Certifique-se de ter o Go instalado na sua máquina para executar o comando acima.

## Endpoints
### `GET /`
Endpoint de healthcheck
```bash
curl "http://localhost:8080/
```

### `GET /notifications/{id}/status`
Consulta o status de um agendamento

```bash
curl "http://localhost:8080/notifications/{id}/status"
```

### `POST /`
Cria um novo agendamento

```bash
curl -X POST -d '{"type": "sms", "recipient": "test"}' "http://localhost:8080/notifications"
```
> Valores possiveis para o campo `type`: `email`, `sms`, `push` e `whatsapp`.

### `DELETE /notifications/{id}`
Cancelar/Excluir um agendamento

```bash
curl -X DELETE "localhost:8080/notifications/{id}"
```

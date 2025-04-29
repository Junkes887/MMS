# MMS API - Variações de Médias Móveis Simples

Serviço que entrega as médias móveis simples (MMS) de 20, 50 e 200 dias para os pares BRLBTC e BRLETH listados no Mercado Bitcoin.

---

## 🚀 Como Rodar o Projeto

### 1. Pré-requisitos
- Go 1.23+
- Docker e Docker Compose (opcional para ambiente)
- PostgreSQL rodando

### 2. Variáveis de Ambiente

Crie um arquivo `.env` com:

```env
PORT=9000
DB_USER=adm
DB_PASSWORD=123
DB_NAME=mms
DB_HOST=localhost
DB_PORT=5432
DELAY_FETCHER_SECONDS=120
EMAIL_SENHA=senha-do-app-gmail
EMAIL_ALERT=seu@email.com
```

### 3. Build e execução local

```bash
# Instalar dependências
go mod tidy

# Rodar a aplicação localmente
go run cmd/mms/main.go
```

### 4. Build da imagem Docker

```bash
docker build -t mms-api .
```

Rodar com Docker:

```bash
docker run -p 9000:9000 --env-file .env mms-api
```

### 5. Docker Compose (opcional)

Para rodar API + Banco juntos rode o comando `docker-compose.yml`.

---

## 🗺️ Rotas Disponíveis

### Buscar MMS de um par

**GET** `/:pair/mms`

**Parâmetros:**

| Parâmetro | Tipo    | Obrigatório | Descrição                           |
|-----------|---------|-------------|-------------------------------------|
| `pair`    | Path    | Sim         | BRLBTC ou BRLETH                    |
| `from`    | Query   | Sim         | Timestamp Unix de início            |
| `to`      | Query   | Não          | Timestamp Unix de fim (padrão: hoje)|
| `range`   | Query   | Sim         | 20, 50 ou 200 dias                  |

**Exemplo de chamada:**

```bash
curl "http://localhost:9000/BRLBTC/mms?from=1713744000&to=1716336000&range=20"
```

**Resposta:**

```json
[
  { "timestamp": 1713744000, "mms": 123456.78 },
  { "timestamp": 1713830400, "mms": 124789.45 }
]
```

---

## 🛡️ Estratégia de Resiliência

O job diário é resiliente a falhas da API externa (Mercado Bitcoin):

- Tenta novamente até 5 vezes, com atraso fixo entre as tentativas.
- O valor do atraso é configurado pela variável de ambiente DELAY_FETCHER_SECONDS; se não informado, o valor padrão é 120 segundos.
- Se falhar todas as tentativas, o ciclo é abortado com log de erro.
- O job retoma normalmente na próxima execução (24h depois), continuando a partir do último dia válido.

---

## 🔎 Estratégia de Monitoramento de Dados

O sistema valida diariamente a integridade dos registros dos últimos 365 dias:

- Para cada par (BRLBTC, BRLETH), checa se há um registro por dia.
- Se encontrar dias ausentes, dispara um **e-mail de alerta**.
- O e-mail detalha a lista dos dias faltando.

**Variáveis necessárias:**

- `EMAIL_SENHA`: Senha de App do Gmail para envio.
- `EMAIL_ALERT`: Endereço de e-mail que receberá os alertas.

---

## 🧪 Testes

Para testar as funções principais:

```bash
go test ./...
```

Testes cobrem:
- Serviço de busca de MMS
- Repositório de acesso ao banco
- Job de incremento diário

---

## 📚 Documentação API do Mercado Bitcoin
- [Docs Oficiais](https://api.mercadobitcoin.net/api/v4/docs)


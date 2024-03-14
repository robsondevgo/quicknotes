# Desenvolvimento Web com Go - Do Zero ao Deploy

## Quicknotes

Aplicação web desenvolvida durante o curso. Trata-se de uma aplicação de gerenciamento de anotações e lembretes, similar ao Google Keep.

## Configuração

A aplicação pode ser configurada através de variáveis de ambiente. As variáveis disponíveis podem ser vistas no arquivo .env ou na struct Config do pacote main.

A struct Config possui uma tag (env) que define o nome da variável de ambiente e um valor default ou a palavra required, para os casos em que o valor precisa obrigatoriamente vir de uma variável de ambiente (por exemplo, valores confidenciais são exemplos de valores que devem ser configurados em variáveis de ambiente).

Abaixo podemos ver um exemplo de como configurar uma propriedade na struct Config.

```go
type Config struct {
    NomePropriedade string `env:"NOME_ENV_VAR,valor_default"`
    SecretValue string `env:"SECRET_VALUE,required"`
}
```

- NomePropriedade: nome da propriedade de configuração
- NOME_ENV_VAR: nome da variável de ambiente de onde o valor será lido

## Rotas da aplicação

| Método | Rota            | Handler    | Descrição                         |
|:-------|:----------------|:-----------|:----------------------------------|
| GET    | /               | NoteList   | Home Page                         |
| GET    | /note/{id}      | NoteView   | Visualiza uma anotação            |
| GET    | /note/new       | NoteNew    | Form de Criação de uma anotação   |
| POST   | /note/save      | NoteSave   | Cria uma anotação                 |
| DELETE | /note/{id}      | NoteEdit   | Remove uma anotação               |
| GET    | /note/{id}/edit | NoteEdit   | Form de alteração de uma anotação |
| GET    | /user/signup    | SignupForm | Form de registro de usuários      |
| POST   | /user/signup    | Signup     | Adiciona o usuário no banco       |

## Modelo do Banco de Dados

### NOTES

| CAMPO      | TIPO      | CONSTRAINT   |
|:-----------|:----------|:-------------|
| ID         | BIGSERIAL | PK, NOT NULL |
| TITLE      | TEXT      | NOT NULL     |
| CONTENT    | TEXT      |              |
| COLOR      | TEXT      | NOT NULL     |
| CREATED_AT | TIMESTAMP |              |
| UPDATED_AT | TIMESTAMP |              |

### USERS

| CAMPO      | TIPO      | CONSTRAINT             |
|:-----------|:----------|:-----------------------|
| ID         | BIGSERIAL | PK, NOT NULL           |
| EMAIL      | TEXT      | NOT NULL UNIQUE        |
| PASSWORD   | TEXT      | NOT NULL               |
| ACTIVE     | TEXT      | NOT NULL DEFAULT false |
| CREATED_AT | TIMESTAMP |                        |
| UPDATED_AT | TIMESTAMP |                        |
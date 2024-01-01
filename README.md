# Desenvolvimento Web com Go - Do Zero do Deploy

## Quicknotes

Aplicação web desenvolvida durante o curso. Trata-se de uma aplicação de gerenciamento de anotações e lembretes, similar ao Google Keep.

## Rotas da aplicação

| Método | Rota         | Handler         | Descrição                              |
|:-------|:-------------|:----------------|:---------------------------------------|
| ALL    | /            | noteList        | Home Page                              |
| ALL    | /note/view   | noteView        | Visualiza uma anotação                 |
| ALL    | /note/new    | noteNew         | Form de Criação de uma anotação        |
| POST   | /note/create | noteCreate      | Cria uma anotação                      |
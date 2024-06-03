# Стена Комментариев на GraphQL

Система для добавления и чтения постов и комментариев с использованием GraphQL, аналогичная комментариям к постам на платформах, таких как Хабр или Reddit.

## Характеристики системы постов:
- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до 2000 символов.
- Система пагинации для получения списка комментариев.

## Запуск приложения

Склонируйте репозиторий и перейдите в корневую папку проекта:

```bash
git clone <URL репозитория>
cd <корневая папка проекта>
```

### 🗄️ Запуск приложения с базой данных Redis

```bash
make docker_build DB_TYPE=redis
make docker DB_TYPE=redis
```

### 🗃️ Запуск приложения с базой данных PostgreSQL

```bash
make docker_build DB_TYPE=postgres
make docker DB_TYPE=postgres
```

## Пример использования

### 📌 Создание поста:

```graphql
mutation {
  createPost(title: "Заголовок", content: "Содержание", commentsDisabled: false) {
    id
    title
    content
    commentsActive
    createdAt
    updatedAt
  }
}
```

### 💬 Создание комментария к посту:

```graphql
mutation {
  createComment(postId: "post_id", content: "Это комментарий") {
    id
    postId
    content
    createdAt
    updatedAt
  }
}
```

### 💬📥 Создание ответа на комментарий:

```graphql
mutation {
  createComment(postId: "post_id", parentId: "parent_comment_id", content: "Это ответ на комментарий") {
    id
    postId
    parentId
    content
    createdAt
    updatedAt
  }
}
```

### 📄 Получение данных о постах и комментариях:

```graphql
query {
  posts {
    id
    title
    content
    comments {
      id
      content
      replies {
        id
        content
      }
    }
  }
}
```

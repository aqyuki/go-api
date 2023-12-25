# 実装されている API について

## API について

実装されているすべての API は`/api/v1`prefix がついています。

## 認証系 API

現在実装されている認証系の API は以下のとおりです。

### `POST /signup`

アカウントを新規作成します。

**リクエストボディー**

```json
{
  "user_id": "foo",
  "password": "bar",
  "name": "bar",
  "bio": "baz"
}
```

`Content-Type`に`application/json`を指定してリクエストを送信する必要があります。

**レスポンス**

```json
{
  "token": "TOKEN"
}
```

JWT トークンが返されます。一部の API にアクセスするにはこのトークンを指定する必要があります。

### `POST /signin`

アカウントにサインインします。

**リクエストボディー**

```json
{
  "user_id": "foo",
  "password": "bar"
}
```

`Content-Type`に`application/json`を指定してリクエストを送信する必要があります。

**レスポンス**

```json
{
  "token": "TOKEN"
}
```

JWT トークンが返されます。一部の API にアクセスするにはこのトークンを指定する必要があります。

## 情報を取得する API

現在実装されている API は以下のとおりです。

### `GET /users/:user_id`

**この API は JWT トークンを指定する必要があります。**

`Authorization`ヘッダーに JWT トークンを指定してリクエストを送信してください。

```txt
Authorization: Bearer <YOUR JWT TOKEN>
```

**レスポンス**

```json
{
  "id": "user_id",
  "name": "name",
  "bio": "your bio"
}
```

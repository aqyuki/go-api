# README

このアプリケーションは JWT を使った認証についてのデモプログラムが含まれています。

## 起動方法

### API サーバー

`make server`を実行すると、ポート番号**8080**でサーバーが起動します。

### リクエストの送信方法

1. Sign Up

```bash
curl -X POST -H "Content-Type: application/json" -d '{"user_id":"foo","password":"bar","name":"bar","bio":"baz"}' http://localhost:8080/api/v1/signup
```

以下のフォーマットに沿ったデータをリクエストボディーに付与して送信すると登録される。

```json
{
  "user_id": "foo",
  "password": "bar",
  "name": "bar",
  "bio": "baz"
}
```

2. Sign In

```bash
curl -X POST -H "Content-Type: application/json" -d '{"user_id":"foo","password":"bar"}' http://localhost:8080/api/v1/signin
```

以下のフォーマットに沿ったデータをリクエストボディーに付与して送信するとログインされる。

```json
{
  "user_id": "foo",
  "password": "bar"
}
```

3. User Information

```bash
curl -X GET -H "Authorization: Bearer <YOUR_TOKEN>" http://localhost:8080/api/v1/users/<user_id>
```

- `YOUR_TOKEN`

  登録・ログイン時に返されるレスポンス内のトークンを指定

- `user_id`

  情報を表示したいユーザーの ID

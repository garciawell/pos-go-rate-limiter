# RATE LIMITER



## 1 - passo

Rodar o comando do docker para subir o container do app e DB REDIS.

```
    docker-compose  up -d
```

<br/>

## 2 - Middleware

Acessar o arquivo e executar a request `api.http` para enviar uma request, o middleware irá interceptar a request e adicionar as informação de ip e/ou API_KEY no REDIS.
Caso o valor de API_KEY não seja enviado, somente irá salvar os dados relacionados ao IP.

A API_KEY sempre irá ter preferencia relacionado ao IP, portanto, se o limite por IP é de 10 req/s e um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

<br/>


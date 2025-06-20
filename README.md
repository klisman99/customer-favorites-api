# Customer Favorites API

Essa é uma API Restful que permite gerenciar produtos favoritos de clientes.

**Importante:** Esta API possui integração com a API https://fakestoreapi.com/products para obter os produtos. Portanto, o ID dos produtos devem ser utilizados para adicionar e remover produtos favoritos.

### Funcionalidades

- Autenticação de usuário
- Listagem de clientes
- Listagem de produtos favoritos de um cliente
- Adição de um produto favorito a um cliente
- Remoção de um produto favorito de um cliente
- Atualização de um cliente
- Remoção de um cliente

### Tecnologias

- Go com Gin
- Docker
- PostgreSQL
- Swagger
- JWT
- Bcrypt

### Instalação

Esta aplicação utiliza docker para portabilizar a aplicação e banco de dados.

Para iniciar a aplicação, siga os passos abaixo:

1. Clone o repositório:
```
git clone https://github.com/klisman99/customer-favorites-api.git
```

2. Navegue até o diretório do projeto:
```
cd customer-favorites-api
```

3. Construa e inicie o container:
```
docker-compose up --build
```

4. Acesse a API em:
```
http://localhost:3002
```

### Documentação

Com o sistema iniciado, você pode acessar a documentação da API em:
```
http://localhost:3002/swagger/index.html
```

Para acessar as rotas autenticadas, você deve adicionar o token no header da requisição:
```
Authorization: Bearer <token>-
```

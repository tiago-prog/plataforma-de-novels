# Estrutura API de Login - Resumo

## Endpoints Principais

### 1. Login
**POST /api/auth/login**

**Request:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "minhasenha123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Login realizado com sucesso",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "refresh_token_here_longer_lived",
  "expires_in": 86400,
  "user": {
    "id": "12345",
    "username": "usuario_exemplo",
    "email": "usuario@exemplo.com",
    "role": "creator",
    "avatar": "https://example.com/avatars/user123.jpg",
    "verified": true
  }
}
```

**Response (Error - 401):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email ou senha inválidos"
  }
}
```

### 2. Registro
**POST /api/auth/register**

**Request:**
```json
{
  "username": "novo_usuario",
  "email": "novo@exemplo.com",
  "password": "senhaSegura123!",
  "confirm_password": "senhaSegura123!",
  "role": "reader",
  "terms_accepted": true
}
```

### 3. Refresh Token
**POST /api/auth/refresh-token**

**Request:**
```json
{
  "refresh_token": "refresh_token_here_longer_lived"
}
```

### 4. Logout
**POST /api/auth/logout**

**Request:**
```json
{
  "refresh_token": "refresh_token_to_invalidate"
}
```

## Códigos de Erro

| Código | HTTP Status | Descrição |
|--------|-------------|-----------|
| `INVALID_CREDENTIALS` | 401 | Email ou senha inválidos |
| `TOKEN_EXPIRED` | 401 | Token de acesso expirado |
| `EMAIL_ALREADY_EXISTS` | 409 | Email já cadastrado |
| `WEAK_PASSWORD` | 422 | Senha muito fraca |
| `ACCOUNT_LOCKED` | 423 | Conta bloqueada temporariamente |

## Configurações JWT

- **Algoritmo:** HS256
- **Access Token TTL:** 24 horas (86400 segundos)
- **Refresh Token TTL:** 30 dias (2592000 segundos)
- **Issuer:** "plataforma-conteudo-api"

## Headers Obrigatórios

Para requisições autenticadas:
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

## Rate Limiting

- **Login:** 5 tentativas por IP a cada 15 minutos
- **Registro:** 3 registros por IP a cada hora
- **Refresh Token:** 10 tentativas por minuto

## Papéis de Usuário

- **Reader:** Ler conteúdo, comentar, curtir, seguir autores
- **Creator:** Todas as permissões de reader + criar/editar conteúdo próprio
- **Admin:** Todas as permissões + moderação e administração

## Endpoints Complementares

| Endpoint | Método | Descrição | Auth Required |
|----------|--------|-----------|---------------|
| `/api/auth/verify-email` | POST | Verifica email via código | Não |
| `/api/auth/forgot-password` | POST | Inicia recuperação de senha | Não |
| `/api/auth/reset-password` | POST | Redefine senha com token | Não |
| `/api/auth/me` | GET | Dados do usuário autenticado | Sim |
| `/api/auth/change-password` | PUT | Altera senha do usuário | Sim |

## Implementação Frontend (isso aqui é sobre implementação do JWT)

### Armazenamento
```javascript
localStorage.setItem('access_token', response.token);
// Duvido que vamos implementar o refresh Token, mas se conseguir top xD
localStorage.setItem('refresh_token', response.refresh_token);
```

### Configuração Axios
```javascript
axios.defaults.headers.common['Authorization'] = 
  `Bearer ${localStorage.getItem('access_token')}`;
```

### Interceptor para Renovação Automática
```javascript
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const newToken = await refreshAccessToken();
      error.config.headers['Authorization'] = `Bearer ${newToken}`;
      return axios.request(error.config);
    }
    return Promise.reject(error);
  }
);
```
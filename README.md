# Plano de Projeto: Plataforma de Conteúdo (1h/dia)



### Organização em Sprints Semanais

- **Duração:** 1 semana
- **Capacidade Semanal:** 5h (1h/dia, 5 dias por semana)
- **Meta de Pontos/Semana:** 3 a 5 PF
- **Foco por Sprint:** Um conjunto coeso de funcionalidades

O projeto será dividido em **Fases**, cada uma com múltiplos Sprints, visando a entrega de versões testáveis e funcionais ao final de cada fase.

---


### 🧱 Módulo 1: Fundação e Autenticação

| ID    | Tarefa                              | Descrição                                                                 | PF  |
|:-----:|:------------------------------------|:--------------------------------------------------------------------------|:---:|
| A01   | Configuração do Ambiente (Backend)  | Projeto Golang/gin-gonic, estrutura, DB                               |  5  |
| A02   | Configuração do Ambiente (Frontend) | Projeto React, estrutura, roteamento, temas                              |  3  |
| A03   | Modelo de Usuário                   | Campos: perfil, papel, senha segura                                      |  3  |
| A04   | Registro de Usuário (API)           | `POST /api/auth/register` com validação e hash de senha                  |  5  |
| A05   | Login (API)                         | `POST /api/auth/login`, retorna JWT                                      |  5  |
| A06   | Middleware de Autenticação          | Verificação de JWT                                                       |  3  |
| A07   | Tela de Registro                    | Formulário com integração à API                                          |  3  |
| A08   | Tela de Login                       | Formulário, salvamento de token, redirecionamento                        |  3  |
| A09   | Gerenciamento de Sessão             | Lógica de login/logout no frontend                                       |  2  |
| A10   | Perfil de Usuário (API)             | `GET` e `PUT` em `/api/users/profile`                                    |  5  |
| A11   | Tela de Perfil                      | Visualizar e editar informações do usuário                               |  3  |
|       | **Total**                           |                                                                           | **40** |

### ✍️ Módulo 2: Criação e Publicação de Conteúdo

| ID    | Tarefa                              | Descrição                                                                 | PF  |
|:-----:|:------------------------------------|:--------------------------------------------------------------------------|:---:|
| C01   | Modelo de Obras                     | Campos: título, sinopse, tipo, autor, status                             |  5  |
| C02   | Modelo de Capítulos/Episódios       | Texto (novels) e Imagens (webtoons)                                      |  5  |
| C03   | CRUD de Obras                       | Endpoints completos                                                      |  8  |
| C04   | CRUD de Capítulos/Episódios         | Endpoints completos                                                      |  8  |
| C05   | Upload de Imagens (Backend)         | Recebimento, otimização e salvamento                                     |  8  |
| C06   | Dashboard do Criador                | Página de gerenciamento de obras                                         |  5  |
| C07   | Formulário de Obra                  | Frontend para criar/editar metadados                                     |  5  |
| C08   | Editor de Texto (Novels)            | Editor WYSIWYG                                                           |  8  |
| C09   | Uploader de Imagens (Webtoons)      | Interface de upload e ordenação                                          |  8  |
| C10   | Agendamento de Publicações          | Sistema com cron job                                                     | 13  |
|       | **Total**                           |                                                                           | **73** |

### 📖 Módulo 3: Consumo e Descoberta de Conteúdo

| ID    | Tarefa                              | Descrição                                                                 | PF  |
|:-----:|:------------------------------------|:--------------------------------------------------------------------------|:---:|
| R01   | Endpoints Públicos de Leitura       | Obras, capítulos e episódios                                             |  8  |
| R02   | Página de Detalhes da Obra          | Sinopse, capa, capítulos/episódios                                       |  5  |
| R03   | Interface de Leitura (Novels)       | Formatação: fonte, tamanho, cor                                          |  5  |
| R04   | Interface de Leitura (Webtoons)     | Rolagem vertical de imagens                                              |  5  |
| R05   | Página Inicial                       | Destaques, populares, lançamentos                                        |  5  |
| R06   | Sistema de Busca (Backend)          | `GET /api/search`                                                        |  8  |
| R07   | Página de Exploração/Busca          | Filtros, categorias                                                      |  5  |
|       | **Total**                           |                                                                           | **41** |

### 💬 Módulo 4: Interação e Comunidade

| ID    | Tarefa                              | Descrição                                                                 | PF  |
|:-----:|:------------------------------------|:--------------------------------------------------------------------------|:---:|
| I01   | Modelo de Comentários               | Suporte a respostas em threads                                           |  3  |
| I02   | CRUD de Comentários                 | Postar, editar, apagar                                                   |  8  |
| I03   | Seção de Comentários (Frontend)     | Exibição e envio                                                         |  5  |
| I04   | Curtidas (Backend)                  | Endpoints para curtir/descurtir                                          |  5  |
| I05   | Botão de Curtir (Frontend)          | Com feedback visual                                                      |  2  |
| I06   | Sistema de Seguidores               | Seguir autores/obras                                                     |  5  |
| I07   | Notificações (Backend)              | Novos capítulos, comentários etc.                                        |  8  |
|       | **Total**                           |                                                                           | **36** |

### 🛡️ Módulo 5: Administração e Moderação

| ID    | Tarefa                              | Descrição                                                                 | PF  |
|:-----:|:------------------------------------|:--------------------------------------------------------------------------|:---:|
| M01   | Middleware de Autorização           | Restringir rotas a admins                                                |  2  |
| M02   | Fila de Moderação                   | Listagem de conteúdo pendente                                            |  5  |
| M03   | Endpoints de Aprovação/Rejeição     | Moderação de conteúdo                                                    |  5  |
| M04   | Dashboard do Admin                  | Interface para aprovação/rejeição                                        |  8  |
| M05   | Sistema de Denúncias                | Envio e visualização de denúncias                                        |  8  |
|       | **Total**                           |                                                                           | **28** |

---

## 3. Priorização e Planejamento por Fases

Abaixo, um resumo com duração estimada para cada fase do projeto:

| Fase                                      | Pontos | Duração Estimada |
|-------------------------------------------|:------:|:----------------:|
| **Fase 1: MVP (Autenticação + Novel)**     | 40     | ~10–12 semanas   |
| **Fase 2: Conteúdo (Publicação/Leitura)** | 73     | ~20–25 semanas   |
| **Fase 3: Webtoon + Interação Social**     | 44     | ~11–13 semanas   |
| **Fase 4: Admin + Notificações**           | 49     | ~12–15 semanas   |
| **Total Geral**                            | **206**| **~53–65 semanas** |

---

## 4. Considerações Finais

Este plano é um guia **flexível** e deve ser revisto semanalmente.

**Recomendações:**

- 🕒 **Consistência diária**: 1h/dia vale mais que maratonas ocasionais
- 🔁 **Revisões semanais**: para ajustar prioridades e expectativas
- ✅ **Foco no MVP**: entregue algo funcional o quanto antes
- 📚 **Documente tudo**: mesmo que seja para você mesmo
- ✨ **Mantenha a motivação com pequenas vitórias**

---

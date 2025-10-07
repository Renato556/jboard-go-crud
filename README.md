# JBoard - CRUD API

Uma API REST robusta e escalável para gerenciamento de vagas de emprego, desenvolvida em Go com MongoDB, focada especialmente em oportunidades "Brazilian Friendly" para desenvolvedores brasileiros.

## 🎯 Objetivo da Aplicação

O JBoard CRUD API é uma aplicação backend completa que permite:
- Gerenciar vagas de emprego com recursos avançados de categorização
- Controlar usuários com diferentes níveis de permissão (FREE/PREMIUM)
- Associar e gerenciar habilidades técnicas dos usuários
- Identificar e categorizar vagas "Brazilian Friendly" para desenvolvedores brasileiros
- Fornecer arquitetura escalável preparada para deploy em Azure Container Apps

## 🚀 Funcionalidades

### Endpoints da API

#### **Gerenciamento de Vagas (Jobs)**
- **POST** `/v1/jobs` - Criar nova vaga de emprego
- **GET** `/v1/jobs` - Listar todas as vagas disponíveis

**Campos Suportados:**
- Título, empresa e URL da vaga
- Tipo de emprego e modalidade (remoto/presencial/híbrido)
- Nível de senioridade e área de atuação
- Localização do escritório
- Prazo de inscrição e data de expiração
- **Brazilian Friendly**: Indicador especial para vagas amigáveis a brasileiros

#### **Gerenciamento de Usuários (Users)**
- **POST** `/v1/users` - Criar novo usuário no sistema
- **GET** `/v1/users` - Buscar informações de usuário
- **PUT** `/v1/users` - Atualizar dados do usuário
- **DELETE** `/v1/users` - Remover usuário do sistema

**Tipos de Usuário:**
- **FREE**: Funcionalidades básicas
- **PREMIUM**: Recursos avançados e análises

#### **Gerenciamento de Habilidades (Skills)**
- **GET** `/v1/skills` - Listar todas as habilidades dos usuários
- **POST** `/v1/skills` - Adicionar nova habilidade a um usuário
- **PUT** `/v1/skills` - Remover habilidade específica
- **DELETE** `/v1/skills` - Deletar todas as habilidades de um usuário

### Características Técnicas

#### **Arquitetura Limpa**
- **Controllers**: Handlers HTTP para processamento de requisições
- **Services**: Lógica de negócio centralizada
- **Repositories**: Camada de acesso a dados
- **Models**: Estruturas de dados e entidades
- **Routers**: Definição e organização de rotas

#### **Recursos Avançados**
- **Brazilian Friendly**: Sistema especializado para vagas brasileiras
- **Validação de Dados**: Validação robusta em todas as camadas
- **Error Handling**: Tratamento de erros padronizado
- **Logging**: Sistema de logs estruturado

## 🔧 Tecnologias Utilizadas

- **Linguagem**: Go 1.25
- **Framework Web**: Gorilla Mux
- **Banco de Dados**: MongoDB / Azure Cosmos DB
- **Containerização**: Docker
- **Cloud**: Azure Container Apps
- **Testes**: Go testing + testify
- **CI/CD**: Azure DevOps / GitHub Actions

## 📦 Instalação e Execução

### Pré-requisitos
- Go 1.25 ou superior
- MongoDB (local) ou Azure Cosmos DB (produção)
- Docker e Docker Compose
- Git

### Instalação Local

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/Renato556/jboard-go-crud.git
   cd jboard-go-crud
   ```

2. **Configurar variáveis de ambiente:**
   ```bash
   # Criar arquivo .env na raiz do projeto
   MONGODB_URI=mongodb://localhost:27017
   MONGODB_DATABASE_NAME=jobboard
   MONGODB_JOB_COLLECTION=jobs
   MONGODB_USER_COLLECTION=users
   ```

3. **Instalar dependências:**
   ```bash
   go mod download
   ```

4. **Executar em modo desenvolvimento:**
   ```bash
   go run main.go
   ```

5. **Acessar a aplicação:**
   ```
   http://localhost:8080
   ```

### Execução com Docker

```bash
# Build da imagem
docker build -t jboard-api .

# Executar container
docker run -p 8080:8080 --env-file .env jboard-api
```

## 🧪 Execução de Testes

### Testes Unitários

```bash
# Executar todos os testes
go test ./...

# Executar testes com coverage
go test -cover ./...

# Executar testes com coverage detalhado
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Testes por Módulo

```bash
# Testes dos controllers
go test ./internal/controllers/...

# Testes dos services
go test ./internal/services/...

# Testes dos repositories
go test ./internal/repositories/...

# Testes dos routers
go test ./internal/routers/...
```

### Cobertura de Testes

```bash
# Gerar relatório de cobertura
go test -coverprofile=coverage.out ./...

# Visualizar relatório no navegador
go tool cover -html=coverage.out -o coverage.html
```

## 🗄️ Configuração de Banco de Dados

### MongoDB Local (Desenvolvimento)

**Instalação e Configuração:**
- **Download**: [MongoDB Community Server](https://www.mongodb.com/try/download/community)
- **Conexão**: `mongodb://localhost:27017`
- **Database**: `jobboard`

**Collections:**
- `jobs`: Armazena as vagas de emprego
- `users`: Dados dos usuários do sistema
- `skills`: Habilidades associadas aos usuários

### Azure Cosmos DB (Produção)

**Características:**
- **API**: MongoDB compatible
- **Alta disponibilidade**: Distribuição global
- **Escalabilidade**: Automática baseada em demanda
- **Backup**: Automático com recuperação point-in-time
- **Performance**: Baixa latência garantida

**Configuração:**
```env
MONGODB_URI=mongodb://seu-cosmos-account.mongo.cosmos.azure.com:10255/
MONGODB_DATABASE_NAME=jobboard-prod
```

## 🔄 Deploy e Workflows

### Azure Deploy Workflow

O projeto utiliza GitHub Actions para deploy automático no Azure Container Apps.

#### Triggers
- **Push**: Branches `main` e `master`
- **Manual**: `workflow_dispatch` para deploy sob demanda

#### Etapas do Pipeline

##### Test and Setup
- Checkout do código
- Setup Go 1.21
- Cache de dependências Go
- Download de dependências (`go mod download`)
- Execução de testes (`go test -v ./...`)

##### Build and Deploy
- Setup Docker Buildx
- Login no Azure
- Obtenção da string de conexão MongoDB
- Login no Azure Container Registry
- Build e push da imagem Docker
- Deploy no Azure Container Apps
- Geração de relatórios de deployment

#### Variáveis de Ambiente
- **AZURE_CONTAINER_REGISTRY**: `jboardregistry`
- **CONTAINER_APP_NAME**: `jboard-go-crud`
- **RESOURCE_GROUP**: `jboard-microservices`
- **IMAGE_NAME**: `jboard-go-crud`
- **MONGODB_RESOURCE_NAME**: `jboard-mongodb`

#### Secrets Necessários
- **AZURE_CREDENTIALS**: Credenciais de service principal
- **ACR_USERNAME**: Usuário do Container Registry
- **ACR_PASSWORD**: Senha do Container Registry
- **MONGODB_DATABASE_NAME**: Nome do banco de dados MongoDB

#### Características Avançadas
- **Dependency Injection**: String de conexão MongoDB obtida dinamicamente
- **Multi-stage Build**: Otimização de imagem Docker
- **Health Checks**: Verificação de deploy bem-sucedido
- **Connection String Enhancement**: Adição automática do parâmetro `w=0` para MongoDB
- **Relatórios Detalhados**: Summary com status e informações de deployment

### 🔒 Segurança em Produção

**⚠️ IMPORTANTE**: Arquitetura de segurança enterprise implementada.

**Características de Segurança:**
- ✅ **Rede Privada**: Execução em VNet privada do Azure
- ✅ **Zero Internet Exposure**: Nenhum endpoint público direto
- ✅ **Comunicação Interna**: Acesso apenas via rede interna
- ✅ **API Gateway**: Acesso externo controlado via gateway
- ✅ **Isolamento Completo**: Máxima segurança por isolamento
- ✅ **SSL/TLS**: Criptografia em todas as comunicações

### CI/CD Pipeline

**Deploy Automático:**
```bash
# Build e push da imagem
docker build -t jobboard-api .
docker tag jobboard-api your-registry.azurecr.io/jobboard-api
docker push your-registry.azurecr.io/jobboard-api
```

## 🤝 Colaboração e Desenvolvimento

### Padrões de Código
- **gofmt**: Formatação automática do código Go
- **golint**: Análise estática de código
- **Conventional Commits**: Mensagens padronizadas
- **Clean Architecture**: Organização em camadas bem definidas

### Estrutura de Pastas
```
internal/
├── controllers/       # Handlers HTTP
├── services/         # Lógica de negócio
├── repositories/     # Acesso a dados
├── models/          # Estruturas de dados
│   └── enums/       # Enumerações
├── routers/         # Definição de rotas
└── config/          # Configurações da aplicação
```

### Contribuindo
1. Fork o projeto
2. Crie uma branch feature (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add some amazing-feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

### Code Review
- Todos os PRs devem passar nos testes
- Cobertura mínima de 80%
- Aprovação de pelo menos 1 reviewer
- Validação automática do CI/CD

## 📄 Licenciamento

### Licença MIT

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

**Resumo da Licença:**
- ✅ Uso comercial
- ✅ Modificação
- ✅ Distribuição
- ✅ Uso privado
- ❌ Responsabilidade
- ❌ Garantia

### Direitos de Uso
- Permitido uso em projetos comerciais
- Permitida modificação do código
- Créditos aos autores originais apreciados
- Não há garantias de funcionamento

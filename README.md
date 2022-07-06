# desafio tempo certo
## Code Challenge: API de solicitação de agendamentos

Desafio da tempo certo, é um "code challenge" oferecido por eles como step de um processo seletivo de Engenharia de Software.

A aplicação API tem como responsabilidade fornecer uma lista de agendamentos.

## Stack

Challenge tempo certo usa como tech stack:
- Golang


## Instalação

A aplicação necessite de um ambiente com [Golang](https://go.dev/doc/install) 1.17 para rodar.
Necessita também do [Docker](https://docs.docker.com/engine/install/ubuntu/) e [Docker Compose](https://docs.docker.com/compose/install/)

Instale as dependências e para rodar a aplicação use o passo-a-passo abaixo:
(caso tenha o [Make](https://linuxhint.com/install-make-ubuntu/) instalado em sua máquina, pule para o exemplo com o Makefile)

### Passo 1:
Crie suas variáveis de ambiente de acordo com o que está no arquivo: `env/application.env`

### Passo 2:
Rode o comando :
```sh
docker-compose up -d
```

### Passo 3:
Em um novo terminal, rode o comando: 
```sh
 go run ./cmd/api main.go
```

Agora basta seguir o tutorial de uso da API.

## Makefile

Uma alternativa para rodar, é usando o Makefile que encontra-se na raiz do projeto.

Em um terminal rode o comando:
```sh
make docker
```

Em um novo terminal, rode o comando:
```sh
make run
```

Agora basta seguir o tutorial de uso da API.

## Uso da API

# Métodos
Requisições para a API devem seguir os padrões:
| Método | Descrição |
|---|---|
| `GET` | Retorna informações de um ou mais agendamentos e disponibilidade|
| `POST` | Utilizado para criar ou atualizar um agendamento. |


## Respostas

| Código | Descrição |
|---|---|
| `200` | Requisição executada com sucesso (success).|
| `400` | Dados enviados inválidos.|

## Criar Lista de horario disponibilidade
A ação de `lista disponibilidade` não precisa enviar nenhum parâmetro, basta chamar o endpoint e criará a lista:

| Método | Endpoint |
|---|---|
| `POST` |  /agendamento/lista |


## Criar agendamentos
A ação de `criar` necessita do envio dos parâmetros como `json` no `body` da request:

| Parâmetro | Descrição |
|---|---|
| `fantasia` | fantasia da empresa que deseja agendar. Ex.: "tempo certo" (campo opcional) |
| `cnpj` |  cnpj da empresa que queira agendar/atualizar. Ex.: 26488705000193(cnpj tempo certo)( se o cnpj já tiver agendado upt só servira para atualizar o horario do agendamento) |
| `inicio` | horario inicio do agendamento. Ex.: "10:20" (caso não passe nada ele pegara o tempo atual) |


| Método | Endpoint |
|---|---|
| `POST` | /agendamento/agendas |

## Buscar agendamento pelo cnpj
A ação de `Buscar agendamento pelo cnpj` necessita do envio do parâmetro:

| Método | Endpoint | Descrição |
|---|---|---|
| `GET` | /agendamento/agendas/{cnpj} | No `{cnpj}` deve-se substituir pelo cnpj da empresa que queira buscar que tenha um agendamento. Ex.: 26488705000193(cnpj tempo certo) |

## Listar agendamentos
A ação de `listar agendamentos` não precisa enviar nenhum parâmetro, basta chamar o endpoint e retonara todos os agendamentos:

| Método | Endpoint |
|---|---|
| `GET` |  /agendamento/agendas/ |

## Listar disponibilidade
A ação de `listar disponibilidade` não precisa enviar nenhum parâmetro, basta chamar o endpoint e retonara todos os horarios disponiveis:

| Método | Endpoint |
|---|---|
| `GET` |  /agendamento/agendas/disponibilidade |

## observações

Caso de 400 nos get é por causa da api  [receitaWS](https://developers.receitaws.com.br/#/operations/queryCNPJFree) que é free, ela não permite fazer muitas request e não coloquei uma validação para não ficar chato, basta esperar alguns segundos e tentar de novo que vai funcionar.

## Notas

Por se tratar de uma linguagem em que não há uma "regra" de arquitetura, utilizei algumas premissas da comunidade e aderente à algumas boas práticas de mercado, nas quais venho aprimorando desde 2019 quando tive o primeiro contato com a linguagem em um monolito.


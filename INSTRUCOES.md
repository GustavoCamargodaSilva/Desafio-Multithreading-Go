# Como Executar o Projeto

## Requisitos

- Go 1.22 ou superior instalado ([download](https://go.dev/dl/))
- Conexão com a internet (para consultar as APIs de CEP)

## Clonando o Repositório

```bash
git clone https://github.com/gustavo/desafio-multithreading.git
cd desafio-multithreading
```

## Executando o Programa

Com o CEP padrão (01153000):

```bash
go run main.go
```

Informando um CEP específico:

```bash
go run main.go 01001000
```

### Exemplo de saída

```
API mais rápida: BrasilAPI
Endereço: CEP: 01001000, Rua: Praça da Sé, Bairro: Sé, Cidade: São Paulo, Estado: SP
```

Se nenhuma API responder em 1 segundo:

```
Erro: timeout - nenhuma API respondeu em 1 segundo
```

## Executando os Testes

Rodar todos os testes:

```bash
go test -v ./...
```

### Testes disponíveis

| Teste | O que valida |
|-------|-------------|
| TestBuscaBrasilAPI | Busca via BrasilAPI com servidor mock |
| TestBuscaViaCep | Busca via ViaCEP com servidor mock |
| TestAPIMaisRapidaVence | A API mais rápida é retornada primeiro |
| TestTimeout | Timeout é acionado quando a API demora |
| TestParseBrasilAPI | Parsing correto do JSON da BrasilAPI |
| TestParseViaCep | Parsing correto do JSON da ViaCEP |
| TestParseJSONInvalido | Tratamento de JSON inválido |

## Estrutura do Projeto

```
.
├── main.go          # Código principal com busca concorrente de CEP
├── main_test.go     # Testes unitários
├── go.mod           # Definição do módulo Go
├── README.md        # Descrição do desafio
└── INSTRUCOES.md    # Este arquivo
```

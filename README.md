# gochain

Uma biblioteca Go para execução sequencial de funções (handlers) com estado e resultado compartilhados, usando generics.

## 📋 Descrição

Este projeto implementa um mecanismo de "chain of responsibility" genérico, onde múltiplos handlers são executados em sequência, compartilhando um contexto e um resultado. Cada handler pode modificar o contexto ou resultado, e a execução para no primeiro erro.

## ✨ Características

- **Genérico**: Usa generics do Go para qualquer tipo de contexto e resultado
- **Sequencial**: Executa handlers em ordem, parando ao primeiro erro
- **Compartilhamento de Estado**: Contexto e resultado compartilhados entre handlers
- **Reflexão**: Permite atualizar campos do contexto/resultado por nome

## 🚀 Como Usar

### Exemplo Básico

```go
package main

import (
	"fmt"
	"github.com/filipewelton/gochain"
)

type Context struct {
	Name string
}

type Result struct {
	Message string
}

func main() {
	chain := gochain.NewChain[Context, Result]()

	chain.Add(func(c *gochain.Chain[Context, Result]) error {
		c.UpdateContext("Name", "João")
		return nil
	}).Add(func(c *gochain.Chain[Context, Result]) error {
		ctx := c.GetContext()
		c.UpdateResult("Message", fmt.Sprintf("Olá, %s!", ctx.Name))
		return nil
	})

	err := chain.Run()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}
	fmt.Println(chain.GetResult().Message) // Output: Olá, João!
}
```

### Tratamento de Erros

```go
chain := gochain.NewChain[Context, Result]()
chain.Add(func(c *gochain.Chain[Context, Result]) error {
	ctx := c.GetContext()
	if ctx.Name == "" {
		return errors.New("nome vazio")
	}
	return nil
})
err := chain.Run()
if err != nil {
	// Tratamento do erro
}
```

## 🔧 API

### Tipos

```go
type Chain[T, U any] struct{}
type Handler[T, U any] func(chain *Chain[T, U]) error
```

- `Chain[T, U]`: Estrutura principal do chain
- `Handler[T, U]`: Função que processa um estágio do chain

### Métodos

- `Add(handler Handler[T, U]) *Chain[T, U]`: Adiciona um handler à cadeia
- `Run() error`: Executa os handlers em ordem
- `GetContext() T`: Retorna o contexto atual
- `GetResult() U`: Retorna o resultado atual
- `UpdateContext(fieldName string, value any) error`: Atualiza um campo do contexto por nome
- `UpdateResult(fieldName string, value any) error`: Atualiza um campo do resultado por nome

## 📦 Dependências

## ⚠️ Sobre os Generics T e U

Os tipos genéricos `T` (contexto) e `U` (resultado) **devem ser structs**. Isso é necessário porque a biblioteca utiliza a biblioteca `reflect` para atualizar campos por nome, o que só é possível com structs em Go. O uso de outros tipos (como tipos primitivos, slices ou maps) não é suportado e pode causar panics ou comportamentos inesperados.

- `github.com/onsi/gomega`: Para testes
- `github.com/onsi/ginkgo/v2`: Para testes BDD

## 🧪 Testes

Execute os testes com:

```bash
go test ./...
```

## 📝 Licença

Veja o arquivo [LICENSE](LICENSE) para detalhes.

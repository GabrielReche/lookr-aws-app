package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda" // Pacote para AWS Lambda
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// LambdaCmd define o comando `lambda` para o CLI
var LambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Query AWS Lambda functions in different regions", // Descrição breve do comando
	Run:   queryLambda, // Função a ser executada quando o comando `lambda` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(LambdaCmd) // Adiciona o comando `lambda` como um subcomando do comando raiz
}

// queryLambda é a função que executa a lógica para consultar funções Lambda
func queryLambda(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Function Name", "Region", "Runtime", "Handler", "Memory (MB)", "Timeout (s)", "ARN"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		lambdaClient := lambda.New(sess) // Cria um novo cliente Lambda com a sessão configurada

		input := &lambda.ListFunctionsInput{} // Cria um input para listar funções Lambda
		result, err := lambdaClient.ListFunctions(input) // Lista as funções Lambda na região atual
		if err != nil {
			fmt.Println("failed to list AWS Lambda functions,", err) // Imprime erro se a listagem de funções falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, function := range result.Functions { // Itera sobre cada função Lambda na lista de funções
			row := []string{
				*function.FunctionName,         // Nome da função Lambda
				regionName,                     // Nome da região
				*function.Runtime,              // Runtime da função Lambda (ex: nodejs, python)
				*function.Handler,              // Handler da função Lambda
				fmt.Sprintf("%d", *function.MemorySize), // Tamanho da memória em MB
				fmt.Sprintf("%d", *function.Timeout),    // Timeout da função em segundos
				*function.FunctionArn,          // ARN da função Lambda
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

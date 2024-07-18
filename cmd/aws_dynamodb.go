package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb" // Pacote para Amazon DynamoDB
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// DynamoDBCmd define o comando `dynamodb` para o CLI
var DynamoDBCmd = &cobra.Command{
	Use:   "dynamodb",
	Short: "Query Amazon DynamoDB tables in different regions", // Descrição breve do comando
	Run:   queryDynamoDB, // Função a ser executada quando o comando `dynamodb` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(DynamoDBCmd) // Adiciona o comando `dynamodb` como um subcomando do comando raiz
}

// queryDynamoDB é a função que executa a lógica para consultar tabelas DynamoDB
func queryDynamoDB(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Table Name", "Region", "Status", "Item Count", "Size (Bytes)", "Provisioned Throughput", "arn"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		dynamoDBClient := dynamodb.New(sess) // Cria um novo cliente DynamoDB com a sessão configurada

		input := &dynamodb.ListTablesInput{} // Cria um input para listar tabelas DynamoDB

		result, err := dynamoDBClient.ListTables(input) // Lista as tabelas DynamoDB na região atual
		if err != nil {
			fmt.Println("failed to list Amazon DynamoDB tables,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, tableName := range result.TableNames { // Itera sobre cada nome de tabela listado
			describeInput := &dynamodb.DescribeTableInput{
				TableName: tableName,
			}

			tableDetails, err := dynamoDBClient.DescribeTable(describeInput) // Descreve a tabela DynamoDB atual
			if err != nil {
				fmt.Println("failed to describe DynamoDB table,", err) // Imprime erro se a descrição falhar
				return
			}

			provisionedThroughput := ""
			if tableDetails.Table.ProvisionedThroughput != nil { // Verifica se há throughput provisionado
				provisionedThroughput = fmt.Sprintf("Read: %d, Write: %d",
					*tableDetails.Table.ProvisionedThroughput.ReadCapacityUnits,
					*tableDetails.Table.ProvisionedThroughput.WriteCapacityUnits) // Define o throughput provisionado
			}

			// Cria uma linha com os detalhes da tabela para adicionar à tabela
			row := []string{
				*tableName,
				regionName,
				*tableDetails.Table.TableStatus,
				fmt.Sprintf("%d", *tableDetails.Table.ItemCount),
				fmt.Sprintf("%d", *tableDetails.Table.TableSizeBytes),
				provisionedThroughput,
				*tableDetails.Table.TableArn,
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

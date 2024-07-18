package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs" // Pacote para AWS SQS
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// SqsCmd define o comando `sqs` para o CLI
var SqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Query Amazon SQS queues in different regions", // Descrição breve do comando
	Run:   querySQS, // Função a ser executada quando o comando `sqs` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(SqsCmd) // Adiciona o comando `sqs` como um subcomando do comando raiz
}

// querySQS é a função que executa a lógica para consultar filas Amazon SQS
func querySQS(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Queue Name", "Region", "Visibility Timeout", "Approximate Messages", "Created Timestamp", "Arn"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		sqsClient := sqs.New(sess) // Cria um novo cliente SQS com a sessão configurada

		input := &sqs.ListQueuesInput{} // Cria um input para listar filas SQS
		result, err := sqsClient.ListQueues(input) // Lista as filas SQS na região atual
		if err != nil {
			fmt.Println("failed to list Amazon SQS queues,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, queueURL := range result.QueueUrls { // Itera sobre cada URL de fila na lista de URLs
			getQueueAttributesInput := &sqs.GetQueueAttributesInput{
				QueueUrl: queueURL, // URL da fila atual
				AttributeNames: []*string{
					aws.String("VisibilityTimeout"), // Atributo: VisibilityTimeout
					aws.String("ApproximateNumberOfMessages"), // Atributo: ApproximateNumberOfMessages
					aws.String("CreatedTimestamp"), // Atributo: CreatedTimestamp
					aws.String("Arn"), // Atributo: Arn
				},
			}

			attributes, err := sqsClient.GetQueueAttributes(getQueueAttributesInput)
			if err != nil {
				fmt.Println("failed to get queue attributes,", err) // Imprime erro se falhar ao obter atributos da fila
				return
			}

			row := []string{
				queueNameFromURL(*queueURL), // Nome da fila obtido da URL
				regionName, // Nome da região
				*attributes.Attributes["VisibilityTimeout"], // Timeout de visibilidade
				*attributes.Attributes["ApproximateNumberOfMessages"], // Número aproximado de mensagens
				timestampToTimeString(*attributes.Attributes["CreatedTimestamp"]), // Timestamp de criação formatado
				*attributes.Attributes["Arn"], // ARN da fila
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

// queueNameFromURL extrai o nome da fila a partir da URL da fila
func queueNameFromURL(url string) string {
	parts := splitLast(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return url
}

// splitLast divide a string `s` pelo separador `sep` e retorna a última parte
func splitLast(s, sep string) []string {
	parts := strings.Split(s, sep)
	if len(parts) == 0 {
		return nil
	}
	return parts
}

// timestampToTimeString converte um timestamp em formato string para uma string de tempo formatada
func timestampToTimeString(timestamp string) string {
	timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return timestamp // Retorna o timestamp original se a conversão falhar
	}
	return time.Unix(timestampInt64, 0).String() // Converte o timestamp UNIX para string de tempo formatada
}

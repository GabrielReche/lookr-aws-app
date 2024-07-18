package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront" // Pacote para Amazon CloudFront
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// CloudFrontCmd define o comando `cloudfront` para o CLI
var CloudFrontCmd = &cobra.Command{
	Use:   "cloudfront",
	Short: "Query Amazon CloudFront distributions in different regions", // Descrição breve do comando
	Run:   queryCloudFront, // Função a ser executada quando o comando `cloudfront` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(CloudFrontCmd) // Adiciona o comando `cloudfront` como um subcomando do comando raiz
}

// queryCloudFront é a função que executa a lógica para consultar distribuições CloudFront
func queryCloudFront(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Distribution ID", "Region", "Domain Name", "Status", "Default Cache Behavior", "arn"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		cfClient := cloudfront.New(sess) // Cria um novo cliente CloudFront com a sessão configurada

		input := &cloudfront.ListDistributionsInput{} // Cria um input para listar distribuições CloudFront

		result, err := cfClient.ListDistributions(input) // Lista as distribuições CloudFront na região atual
		if err != nil {
			fmt.Println("failed to list Amazon CloudFront distributions,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, distribution := range result.DistributionList.Items { // Itera sobre cada distribuição listada
			defaultCacheBehavior := "N/A"
			if distribution.DefaultCacheBehavior != nil { // Verifica se há comportamento de cache padrão
				defaultCacheBehavior = *distribution.DefaultCacheBehavior.TargetOriginId // Define o comportamento de cache padrão
			}

			// Cria uma linha com os detalhes da distribuição para adicionar à tabela
			row := []string{
				*distribution.Id,
				regionName,
				*distribution.DomainName,
				*distribution.Status,
				defaultCacheBehavior,
				*distribution.ARN,
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

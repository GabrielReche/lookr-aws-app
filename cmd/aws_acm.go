package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm" // Pacote para AWS Certificate Manager (ACM)
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// AcmCmd define o comando `acm` para o CLI
var AcmCmd = &cobra.Command{
	Use:   "acm",
	Short: "Query AWS Certificate Manager certificates in different regions", // Descrição breve do comando
	Run:   queryACM, // Função a ser executada quando o comando `acm` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(AcmCmd) // Adiciona o comando `acm` como um subcomando do comando raiz
}

// queryACM é a função que executa a lógica para consultar certificados ACM da AWS
func queryACM(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Certificate ARN", "Region", "Domain Name", "Status", "Type", "Validation Method"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		acmClient := acm.New(sess) // Cria um novo cliente ACM com a sessão configurada

		input := &acm.ListCertificatesInput{} // Cria um input para listar certificados ACM

		result, err := acmClient.ListCertificates(input) // Lista os certificados ACM na região atual
		if err != nil {
			fmt.Println("failed to list AWS ACM certificates,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, certificate := range result.CertificateSummaryList { // Itera sobre cada certificado listado
			describeInput := &acm.DescribeCertificateInput{
				CertificateArn: certificate.CertificateArn, // Configura input para descrever o certificado
			}

			certificateDetails, err := acmClient.DescribeCertificate(describeInput) // Descreve o certificado ACM
			if err != nil {
				fmt.Println("failed to describe ACM certificate,", err) // Imprime erro se a descrição falhar
				return
			}

			// Cria uma linha com os detalhes do certificado para adicionar à tabela
			row := []string{
				*certificate.CertificateArn,
				regionName,
				*certificate.DomainName,
				*certificateDetails.Certificate.Status,
				*certificateDetails.Certificate.Type,
				*certificateDetails.Certificate.DomainValidationOptions[0].ValidationMethod,
				*certificate.CertificateArn,
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

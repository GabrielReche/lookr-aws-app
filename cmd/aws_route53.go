package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53" // Pacote para AWS Route 53
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// Route53Cmd define o comando `route53` para o CLI
var Route53Cmd = &cobra.Command{
	Use:   "route53",
	Short: "Query Route 53 hosted zones in different regions", // Descrição breve do comando
	Run:   queryRoute53, // Função a ser executada quando o comando `route53` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(Route53Cmd) // Adiciona o comando `route53` como um subcomando do comando raiz
}

// queryRoute53 é a função que executa a lógica para consultar zonas hospedadas do Route 53
func queryRoute53(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Hosted Zone Name", "Region", "Private", "Record Count"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		route53Client := route53.New(sess) // Cria um novo cliente Route 53 com a sessão configurada

		input := &route53.ListHostedZonesInput{} // Cria um input para listar zonas hospedadas do Route 53
		result, err := route53Client.ListHostedZones(input) // Lista as zonas hospedadas do Route 53 na região atual
		if err != nil {
			fmt.Println("failed to list Route 53 hosted zones,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, hostedZone := range result.HostedZones { // Itera sobre cada zona hospedada na lista de zonas
			isPrivate := "No"
			if *hostedZone.Config.PrivateZone {
				isPrivate = "Yes"
			}

			getHostedZoneInput := &route53.GetHostedZoneInput{
				Id: hostedZone.Id,
			}

			getHostedZoneOutput, err := route53Client.GetHostedZone(getHostedZoneInput)
			if err != nil {
				fmt.Println("failed to get hosted zone,", err)
				continue
			}

			row := []string{
				*hostedZone.Name, // Nome da zona hospedada
				regionName, // Nome da região
				isPrivate, // Indica se é privada
				fmt.Sprintf("%d", *getHostedZoneOutput.HostedZone.ResourceRecordSetCount), // Contagem de registros
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

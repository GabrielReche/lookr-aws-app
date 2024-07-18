package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2" // Pacote para Amazon EC2
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// EbsCmd define o comando `ebs` para o CLI
var EbsCmd = &cobra.Command{
	Use:   "ebs",
	Short: "Query Amazon EBS volumes in different regions", // Descrição breve do comando
	Run:   queryEBS, // Função a ser executada quando o comando `ebs` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(EbsCmd) // Adiciona o comando `ebs` como um subcomando do comando raiz
}

// queryEBS é a função que executa a lógica para consultar volumes EBS
func queryEBS(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Volume ID", "Region", "az", "Size (GB)", "Type", "Status", "IOPS", "Encryption"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		ec2Client := ec2.New(sess) // Cria um novo cliente EC2 com a sessão configurada

		input := &ec2.DescribeVolumesInput{} // Cria um input para descrever volumes EBS

		result, err := ec2Client.DescribeVolumes(input) // Descreve os volumes EBS na região atual
		if err != nil {
			fmt.Println("failed to describe Amazon EBS volumes,", err) // Imprime erro se a descrição falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, volume := range result.Volumes { // Itera sobre cada volume listado
			encryption := "No"
			if volume.Encrypted != nil && *volume.Encrypted { // Verifica se o volume está criptografado
				encryption = "Yes"
			}

			iops := ""
			if volume.Iops != nil { // Verifica se há IOPS configurados para o volume
				iops = fmt.Sprintf("%d", *volume.Iops)
			}

			// Cria uma linha com os detalhes do volume para adicionar à tabela
			row := []string{
				*volume.VolumeId,
				regionName,
				*volume.AvailabilityZone,
				fmt.Sprintf("%d", *volume.Size),
				*volume.VolumeType,
				*volume.State,
				iops,
				encryption,
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

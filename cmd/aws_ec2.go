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

// EC2Cmd define o comando `ec2` para o CLI
var EC2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Query EC2 instances in different regions", // Descrição breve do comando
	Run:   queryEC2, // Função a ser executada quando o comando `ec2` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(EC2Cmd) // Adiciona o comando `ec2` como um subcomando do comando raiz
}

// queryEC2 é a função que executa a lógica para consultar instâncias EC2
func queryEC2(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Instance ID", "Region", "Instance Type", "State", "Private IP", "Public IP"}) // Define cabeçalhos da tabela

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

		input := &ec2.DescribeInstancesInput{} // Cria um input para descrever instâncias EC2

		result, err := ec2Client.DescribeInstances(input) // Descreve as instâncias EC2 na região atual
		if err != nil {
			fmt.Println("failed to describe EC2 instances,", err) // Imprime erro se a descrição falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, reservation := range result.Reservations { // Itera sobre cada reserva de instâncias
			for _, instance := range reservation.Instances { // Itera sobre cada instância na reserva
				row := []string{
					*instance.InstanceId, // ID da instância
					regionName, // Nome da região
					*instance.InstanceType, // Tipo da instância
					*instance.State.Name, // Estado da instância
					*instance.PrivateIpAddress, // Endereço IP privado da instância
					*instance.PublicIpAddress, // Endereço IP público da instância
				}
				table.Append(row) // Adiciona a linha à tabela
			}
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

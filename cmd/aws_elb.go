package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2" // Pacote para ELBv2 (Elastic Load Balancing)
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// ElbCmd define o comando `elb` para o CLI
var ElbCmd = &cobra.Command{
	Use:   "elb",
	Short: "Query ELB Load Balancers in different regions", // Descrição breve do comando
	Run:   queryELB, // Função a ser executada quando o comando `elb` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(ElbCmd) // Adiciona o comando `elb` como um subcomando do comando raiz
}

// queryELB é a função que executa a lógica para consultar ELB Load Balancers
func queryELB(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Load Balancer Name", "Region", "DNS Name", "Scheme", "Type", "State", "ARN"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		elbv2Client := elbv2.New(sess) // Cria um novo cliente ELBv2 com a sessão configurada

		input := &elbv2.DescribeLoadBalancersInput{} // Cria um input para descrever ELB Load Balancers

		result, err := elbv2Client.DescribeLoadBalancers(input) // Descreve os ELB Load Balancers na região atual
		if err != nil {
			fmt.Println("failed to describe ELB Load Balancers,", err) // Imprime erro se a descrição falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, lb := range result.LoadBalancers { // Itera sobre cada ELB Load Balancer na lista de Load Balancers
			row := []string{
				*lb.LoadBalancerName, // Nome do Load Balancer
				regionName,           // Nome da região
				*lb.DNSName,          // DNS Name do Load Balancer
				*lb.Scheme,           // Scheme do Load Balancer
				*lb.Type,             // Tipo do Load Balancer
				*lb.State.Code,       // Estado do Load Balancer
				*lb.LoadBalancerArn,  // ARN do Load Balancer
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

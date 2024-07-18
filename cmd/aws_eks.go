package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks" // Pacote para Amazon EKS
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// EksCmd define o comando `eks` para o CLI
var EksCmd = &cobra.Command{
	Use:   "eks",
	Short: "Query EKS clusters in different regions", // Descrição breve do comando
	Run:   queryEKS, // Função a ser executada quando o comando `eks` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(EksCmd) // Adiciona o comando `eks` como um subcomando do comando raiz
}

// queryEKS é a função que executa a lógica para consultar clusters EKS
func queryEKS(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Cluster Name", "Region", "Status", "Endpoint", "Kubernetes Version", "Arn"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		eksClient := eks.New(sess) // Cria um novo cliente EKS com a sessão configurada

		input := &eks.ListClustersInput{} // Cria um input para listar clusters EKS

		result, err := eksClient.ListClusters(input) // Lista os clusters EKS na região atual
		if err != nil {
			fmt.Println("failed to list EKS clusters,", err) // Imprime erro se a listagem falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, clusterName := range result.Clusters { // Itera sobre cada nome de cluster na lista de clusters
			describeInput := &eks.DescribeClusterInput{
				Name: aws.String(*clusterName), // Nome do cluster a ser descrito
			}

			clusterDetails, err := eksClient.DescribeCluster(describeInput) // Descreve o cluster EKS
			if err != nil {
				fmt.Println("failed to describe EKS cluster,", err) // Imprime erro se a descrição falhar
				return
			}

			cluster := clusterDetails.Cluster // Obtém detalhes do cluster

			row := []string{
				*cluster.Name,         // Nome do cluster
				regionName,            // Nome da região
				*cluster.Status,       // Estado do cluster
				*cluster.Endpoint,     // Endpoint do cluster
				*cluster.Version,      // Versão do Kubernetes do cluster
				*cluster.Arn,          // ARN do cluster
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

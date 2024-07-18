package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds" // Pacote para Amazon RDS (Relational Database Service)
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// AuroraCmd define o comando `aurora` para o CLI
var AuroraCmd = &cobra.Command{
	Use:   "aurora",
	Short: "Query Amazon Aurora clusters in different regions", // Descrição breve do comando
	Run:   queryAurora, // Função a ser executada quando o comando `aurora` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(AuroraCmd) // Adiciona o comando `aurora` como um subcomando do comando raiz
}

// queryAurora é a função que executa a lógica para consultar clusters Aurora da Amazon RDS
func queryAurora(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Cluster Name", "Region", "Status", "Engine", "Engine Version", "DB Instances", "Replicas", "arn"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		rdsClient := rds.New(sess) // Cria um novo cliente RDS com a sessão configurada

		input := &rds.DescribeDBClustersInput{} // Cria um input para descrever clusters Aurora

		result, err := rdsClient.DescribeDBClusters(input) // Descreve os clusters Aurora na região atual
		if err != nil {
			fmt.Println("failed to describe Amazon Aurora clusters,", err) // Imprime erro se a descrição falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, cluster := range result.DBClusters { // Itera sobre cada cluster listado
			dbInstances := ""
			for _, instance := range cluster.DBClusterMembers { // Itera sobre cada instância do cluster
				dbInstances += *instance.DBInstanceIdentifier + ", "
			}

			replicaIdentifier := ""
			if len(cluster.ReadReplicaIdentifiers) > 0 { // Verifica se há réplicas
				replicaIdentifier = *cluster.ReadReplicaIdentifiers[0] // Define o identificador da primeira réplica
			}

			// Cria uma linha com os detalhes do cluster para adicionar à tabela
			row := []string{
				*cluster.DBClusterIdentifier,
				regionName,
				*cluster.Status,
				*cluster.Engine,
				*cluster.EngineVersion,
				dbInstances,
				replicaIdentifier,
				*cluster.DBClusterArn,
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

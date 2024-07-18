package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticache" // Pacote para Amazon ElastiCache
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// ElastiCacheCmd define o comando `elasticache` para o CLI
var ElastiCacheCmd = &cobra.Command{
	Use:   "elasticache",
	Short: "Query Amazon ElastiCache clusters in different regions", // Descrição breve do comando
	Run:   queryElastiCache, // Função a ser executada quando o comando `elasticache` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(ElastiCacheCmd) // Adiciona o comando `elasticache` como um subcomando do comando raiz
}

// queryElastiCache é a função que executa a lógica para consultar clusters ElastiCache
func queryElastiCache(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Cluster ID", "Region", "Engine", "Engine Version", "Status", "Cluster Mode", "Node Type", "Nodes", "ARN"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		elastiCacheClient := elasticache.New(sess) // Cria um novo cliente ElastiCache com a sessão configurada

		input := &elasticache.DescribeCacheClustersInput{} // Cria um input para descrever clusters ElastiCache

		result, err := elastiCacheClient.DescribeCacheClusters(input) // Descreve os clusters ElastiCache na região atual
		if err != nil {
			fmt.Println("failed to describe Amazon ElastiCache clusters,", err) // Imprime erro se a descrição falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, cluster := range result.CacheClusters { // Itera sobre cada cluster na lista de clusters
			row := []string{
				*cluster.CacheClusterId, // ID do cluster
				regionName,              // Nome da região
				*cluster.Engine,         // Engine do cluster
				*cluster.EngineVersion,  // Versão da engine do cluster
				*cluster.CacheClusterStatus, // Status do cluster
				*cluster.CacheNodeType, // Tipo de nó do cluster
				fmt.Sprintf("%d", *cluster.NumCacheNodes), // Número de nós do cluster formatado como string
				*cluster.ARN, // ARN do cluster
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

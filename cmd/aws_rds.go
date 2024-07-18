package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds" // Pacote para AWS RDS
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// RdsCmd define o comando `rds` para o CLI
var RdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Query RDS in different regions", // Descrição breve do comando
	Run:   queryRDS, // Função a ser executada quando o comando `rds` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(RdsCmd) // Adiciona o comando `rds` como um subcomando do comando raiz
}

// queryRDS é a função que executa a lógica para consultar instâncias RDS
func queryRDS(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"DB Name", "Region", "AZ", "Status", "Instance Type", "Engine", "Version", "Port", "Storage Type", "Storage Size", "Multi-AZ", "Replica", "ARN"}) // Define cabeçalhos da tabela

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

		input := &rds.DescribeDBInstancesInput{} // Cria um input para descrever instâncias RDS
		result, err := rdsClient.DescribeDBInstances(input) // Descreve as instâncias RDS na região atual
		if err != nil {
			fmt.Println("failed to describe db instances,", err) // Imprime erro se a descrição de instâncias falhar
			return
		}

		for _, dbInstance := range result.DBInstances { // Itera sobre cada instância RDS na lista de instâncias
			hasReadReplica := "No"
			if len(dbInstance.ReadReplicaDBInstanceIdentifiers) > 0 {
				hasReadReplica = "Yes"
			}

			multiAZ := "No"
			if dbInstance.MultiAZ != nil && *dbInstance.MultiAZ {
				multiAZ = "Yes"
			}

			regionName := deps.GetRegionName(region) // Obtém o nome da região atual

			row := []string{
				*dbInstance.DBInstanceIdentifier, // Identificador da instância RDS
				regionName,                       // Nome da região
				*dbInstance.AvailabilityZone,     // Zona de disponibilidade
				*dbInstance.DBInstanceStatus,     // Status da instância
				*dbInstance.DBInstanceClass,      // Tipo da instância
				*dbInstance.Engine,               // Engine do RDS (ex: mysql, postgres)
				*dbInstance.EngineVersion,        // Versão da engine
				fmt.Sprintf("%d", *dbInstance.Endpoint.Port), // Porta da instância
				*dbInstance.StorageType,          // Tipo de armazenamento (ex: gp2)
				fmt.Sprintf("%d", *dbInstance.AllocatedStorage), // Tamanho do armazenamento alocado
				multiAZ,                          // Indica se é Multi-AZ
				hasReadReplica,                   // Indica se tem réplica de leitura
				*dbInstance.DBInstanceArn,        // ARN da instância RDS
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

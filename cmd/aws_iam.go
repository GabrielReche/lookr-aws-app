package cmd

import (
	"fmt"
	"lookr/deps" // Importação de pacotes locais ou dependências
	"os"

	"github.com/aws/aws-sdk-go/aws" // Pacote AWS SDK para Go
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam" // Pacote para IAM (Identity and Access Management)
	"github.com/olekukonko/tablewriter" // Pacote para formatação de tabelas
	"github.com/spf13/cobra" // Pacote para criação de CLI usando Cobra
)

// IAMCmd define o comando `iam` para o CLI
var IAMCmd = &cobra.Command{
	Use:   "iam",
	Short: "Query AWS IAM groups, users, and roles in different regions", // Descrição breve do comando
	Run:   queryIAM, // Função a ser executada quando o comando `iam` é chamado
}

// init é chamado antes da execução do programa principal
func init() {
	rootCmd.AddCommand(IAMCmd) // Adiciona o comando `iam` como um subcomando do comando raiz
}

// queryIAM é a função que executa a lógica para consultar IAM groups, users e roles
func queryIAM(cmd *cobra.Command, args []string) {
	table := tablewriter.NewWriter(os.Stdout) // Cria um novo escritor de tabela que escreve para os.Stdout
	table.SetHeader([]string{"Name", "Type", "Region", "Creation Time", "ARN"}) // Define cabeçalhos da tabela

	AuthRegions := deps.AuthRegions() // Obtém as regiões autorizadas para autenticação
	for _, region := range AuthRegions { // Itera sobre cada região autorizada
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region), // Configura a sessão com a região atual
		})

		if err != nil {
			fmt.Println("failed to create session,", err) // Imprime erro se a sessão não puder ser criada
			return
		}

		iamClient := iam.New(sess) // Cria um novo cliente IAM com a sessão configurada

		listGroupsInput := &iam.ListGroupsInput{} // Cria um input para listar IAM groups
		groupsResult, err := iamClient.ListGroups(listGroupsInput) // Lista os IAM groups na região atual
		if err != nil {
			fmt.Println("failed to list IAM groups,", err) // Imprime erro se a listagem de groups falhar
			return
		}

		listUsersInput := &iam.ListUsersInput{} // Cria um input para listar IAM users
		usersResult, err := iamClient.ListUsers(listUsersInput) // Lista os IAM users na região atual
		if err != nil {
			fmt.Println("failed to list IAM users,", err) // Imprime erro se a listagem de users falhar
			return
		}

		listRolesInput := &iam.ListRolesInput{} // Cria um input para listar IAM roles
		rolesResult, err := iamClient.ListRoles(listRolesInput) // Lista os IAM roles na região atual
		if err != nil {
			fmt.Println("failed to list IAM roles,", err) // Imprime erro se a listagem de roles falhar
			return
		}

		regionName := deps.GetRegionName(region) // Obtém o nome da região atual

		for _, group := range groupsResult.Groups { // Itera sobre cada IAM group na lista de groups
			row := []string{
				*group.GroupName,            // Nome do group
				"Group",                     // Tipo do objeto (Group)
				regionName,                  // Nome da região
				group.CreateDate.String(),   // Data de criação formatada para string
				*group.Arn,                  // ARN do group
			}
			table.Append(row) // Adiciona a linha à tabela
		}

		for _, user := range usersResult.Users { // Itera sobre cada IAM user na lista de users
			row := []string{
				*user.UserName,              // Nome do user
				"User",                      // Tipo do objeto (User)
				regionName,                  // Nome da região
				user.CreateDate.String(),   // Data de criação formatada para string
				*user.Arn,                  // ARN do user
			}
			table.Append(row) // Adiciona a linha à tabela
		}

		for _, role := range rolesResult.Roles { // Itera sobre cada IAM role na lista de roles
			row := []string{
				*role.RoleName,              // Nome do role
				"Role",                      // Tipo do objeto (Role)
				regionName,                  // Nome da região
				role.CreateDate.String(),   // Data de criação formatada para string
				*role.Arn,                  // ARN do role
			}
			table.Append(row) // Adiciona a linha à tabela
		}
	}
	table.Render() // Renderiza a tabela com os resultados
}

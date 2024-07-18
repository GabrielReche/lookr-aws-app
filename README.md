# lookr - A Aws CLI

# O `lookr` CLI é uma ferramenta que permite consultar informações sobre vários serviços da Amazon Web Services (AWS) em diferentes regiões. 
# Ele utiliza o SDK oficial da AWS para Go para interagir com esses serviços e apresenta os resultados em formato tabular.

## Installation

# Para instalar o `lookr` CLI, siga estes passos:

# Certifique-se de ter o Go instalado em sua máquina. Caso contrário, você pode baixá-lo em https://golang.org/dl/

# Clone este repositório ou baixe o arquivo ZIP.

# Navegue até o diretório CLI usando o terminal.

# Execute o seguinte comando para construir o CLI:

```bash

# Constrói o executável `lookr` a partir do arquivo `main.go` localizado no diretório `cmd/lookr`.

go build -o lookr cmd/lookr/main.go

```

Agora você pode executar a CLI usando ./lookr.

ou Mova o executável para uma pasta que está no PATH do seu sistema:

```bash

sudo mv goact /usr/local/bin

```

# Comandos

O CLI lookr fornece os seguintes comandos para consultar informações sobre vários serviços da AWS:

ec2: Consulta informações sobre instâncias EC2.

rds: Consulta informações sobre bancos de dados RDS.

sqs: Consulta informações sobre filas Amazon SQS.

lambda: Consulta informações sobre funções AWS Lambda.

iam: Consulta informações sobre grupos, usuários e funções IAM.

ebs: Consulta informações sobre volumes Amazon EBS.

acm: Consulta informações sobre certificados do AWS Certificate Manager.

cloudfront: Consulta informações sobre distribuições Amazon CloudFront.

elasticache: Consulta informações sobre clusters Amazon ElastiCache.

dynamodb: Consulta informações sobre tabelas Amazon DynamoDB.

#  Uso

Para usar o CLI lookr e consultar informações sobre um serviço específico, execute o seguinte comando:

```shell

./lookr <command>

```

Substitua <command> pelo nome do serviço desejado. Por exemplo, para consultar informações sobre instâncias EC2, use:

```shell

./lookr ec2

```

O CLI exibirá os resultados em formato tabular, mostrando detalhes sobre o serviço na região atual e em outras regiões configuradas.

Licença

Este projeto está licenciado sob a Licença MIT - consulte o arquivo LICENSE para mais detalhes.


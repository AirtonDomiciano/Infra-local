# EnvScout

EnvScout é um CLI leve em Go para varrer sua rede local, encontrar hosts com porta aberta (ex.: SQL Server na 1433) e, opcionalmente, atualizar automaticamente uma chave no seu `.env` com o IP encontrado.

Ele foi criado para resolver dores comuns de desenvolvimento quando o IP muda por DHCP e o app perde o endereço do serviço.

## Recursos

- Varredura rápida de rede local (1-254) com concorrência configurável
- Suporte a timeout por host
- Seleção automática do IP encontrado ou preferência por um IP específico
- Atualização automática de chave no `.env`
- Abertura automática de sessão RDP após encontrar o IP

## Requisitos

- Go 1.22+ (o módulo usa `go 1.22.5`)
- `xfreerdp3` instalado (para abertura automática do RDP)

## Uso rápido

Entre na pasta do módulo e rode:

```bash
cd envscout

go run ./cmd/envscout \
  --base 192.168.1. \
  --port 1433
```

Saída típica:

```
✅ Hosts encontrados: [192.168.1.7]
🎯 Escolhido: 192.168.1.7
```

## Atualizar o .env automaticamente

```bash
go run ./cmd/envscout \
  --base 192.168.1. \
  --port 1433 \
  --env /caminho/para/.env \
  --key SERVER_DB
```

Se a chave já existir, ela é atualizada. Se não existir, é adicionada no final do arquivo.

## Preferir um IP específico (com fallback automático)

```bash
go run ./cmd/envscout \
  --base 192.168.1. \
  --port 1433 \
  --prefer 192.168.1.7 \
  --env /caminho/para/.env \
  --key SERVER_DB
```

Se o IP preferido estiver na lista encontrada, ele será escolhido. Caso contrário, o primeiro IP encontrado é usado.

## Build do binário

```bash
cd envscout

go build -o envscout ./cmd/envscout
./envscout --base 192.168.1. --port 1433
```

## Parâmetros

- `--base` (string, padrão `192.168.1.`): base da rede. Use com ponto final, ex.: `192.168.0.`
- `--port` (int, padrão `1433`): porta para testar
- `--timeout` (int, padrão `300`): timeout por host em ms
- `-c` (int, padrão `120`): concorrência (goroutines simultâneas)
- `--env` (string, padrão vazio): caminho do `.env` para atualizar
- `--key` (string, padrão `SERVER_DB`): chave do `.env` a ser atualizada
- `--prefer` (string, padrão vazio): IP preferido, ex.: `192.168.1.7`

## Como funciona

- Varre os IPs de `base + 1` até `base + 254`.
- Para cada IP, tenta abrir conexão TCP na porta informada dentro do timeout.
- Junta os hosts com porta aberta, ordena e escolhe um IP.
- Se `--prefer` estiver definido e o IP estiver na lista, ele é escolhido.
- Após escolher o IP, executa o comando RDP via `xfreerdp3`.
- Se `--env` estiver definido, atualiza `KEY=VALUE` no arquivo.

## Códigos de saída

- `0`: execução bem-sucedida
- `1`: nenhum host encontrado ou falha ao atualizar o `.env`

## Dicas

- Em redes mais lentas, aumente `--timeout` para reduzir falsos negativos.
- Se a rede ficar muito pesada, reduza `-c`.
- Certifique-se de que a base termina com ponto, ex.: `192.168.1.`

## RDP automático (xfreerdp3)

Depois de localizar o IP, o EnvScout executa automaticamente:

```bash
xfreerdp3 /v:<IP_ENCONTRADO> /u:airtonprg\\airton /p:asd123 /sec:nla /cert:ignore /dynamic-resolution /clipboard /network:auto
```

Se quiser alterar usuário, senha ou parâmetros, edite diretamente o arquivo `cmd/envscout/main.go`.

## Executar o RDP manualmente

Caso queira abrir manualmente, você pode usar o alias abaixo e depois executar `win`:

```bash
alias win='xfreerdp3 /v:192.168.1.4 /u:airtonprg\\airton /p:asd123 /sec:nla /cert:ignore /dynamic-resolution /clipboard /network:auto'
win
```

## Alias para acesso remoto (IP local)

Se você estiver na rede local e o IP for `192.168.1.4`, pode usar este alias:

```bash
alias win='xfreerdp3 /v:192.168.1.4 /u:airtonprg\\airton /p:asd123 /sec:nla /cert:ignore /dynamic-resolution /clipboard /network:auto'
```

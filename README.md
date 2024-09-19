# Lab5-conc-BitTorrent

Este é um sistema de cliente-servidor que permite a comunicação entre várias máquinas para compartilhar e buscar arquivos de acordo com seus hashes. O cliente calcula automaticamente os hashes dos arquivos presentes em um diretório específico e os envia ao servidor. Além disso, o cliente pode procurar arquivos por seus hashes no servidor.

## Funcionalidades

- **Atualização automática de hashes**: Quando o cliente se conecta ao servidor, ele envia automaticamente o hash de todos os arquivos presentes no diretório `/tmp/dataset`.
- **Busca de arquivos por hash**: O cliente pode enviar o comando `search <hash>` para buscar um hash específico no servidor. Se o hash for encontrado, o servidor retorna os IPs das máquinas que possuem um arquivo com o hash especificado.
- **Gerenciamento de arquivos**: O servidor mantém um registro dos hashes e das máquinas (IPs) que os possuem.

## Requisitos

- Go instalado nas máquinas servidor e cliente.

## Como Usar

### 1. Configurar e Rodar o Servidor

1. Navegue até o diretório onde o servidor (`server.go`) está localizado.
2. Compile e execute o servidor:

    ```bash
    go run server.go
    ```

   O servidor ficará escutando em um IP e porta específicos (definidos no código).

### 2. Configurar e Rodar o Cliente

1. Navegue até o diretório onde o cliente (`client.go`) está localizado.
2. Certifique-se de que há arquivos no diretório `/tmp/dataset` na máquina cliente. Você pode usar o script `makedataset.sh` para gerar arquivos de teste:

    ```bash
    ./makedataset.sh <número_de_arquivos>
    ```

3. Compile e execute o cliente:

    ```bash
    go run client.go
    ```

   O cliente se conectará automaticamente ao servidor.

### 3. Comandos do Cliente

- **Enviar Hashes**: O cliente calcula e envia automaticamente os hashes dos arquivos ao executar o comando:

    ```
    update
    ```

- **Buscar Hash**: Para buscar um hash específico no servidor, use o comando:

    ```
    search <hash>
    ```

   O servidor retornará os IPs das máquinas que possuem um arquivo com o hash especificado.
- **Atualizar Hashes**: Ao atualizar, remover ou criar novos arquivos em seu diretório /dataset, use novamente o comando update.  

## Observações

- Certifique-se de que o servidor esteja em execução antes de iniciar os clientes.
- O servidor deve estar configurado para escutar no IP e na porta corretos conforme especificado no código.

## Licença

Este projeto é licenciado sob a [MIT License](LICENSE).

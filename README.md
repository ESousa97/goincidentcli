# incident CLI 🛡️

Uma ferramenta moderna e modular em Go para gerenciamento de incidentes diretamente do terminal. O `incident` auxilia no rastreamento, declaração e isolamento de estado de ocorrências técnicas de forma rápida e segura.

## 🚀 Funcionalidades

- **Identificação Híbrida**: IDs gerados no formato `INC-YYYYMMDD-XXXX` para fácil ordenação e referência.
- **Isolamento de Estado**: Cada incidente possui sua própria pasta em `.incidents/[ID]`.
- **Git Friendly**: Criação automática de `.gitignore` para evitar o versionamento acidental de logs e estados locais.
- **Configuração Flexível**: Gerenciamento de tokens e URLs via `~/.incident.yaml` utilizando Viper.
- **Arquitetura Stateless**: Projetado para ser leve e escalável.

## 📦 Instalação

Certifique-se de ter o Go 1.21+ instalado.

```powershell
# Clone o repositório
git clone https://github.com/ESousa97/goincidentcli.git
cd goincidentcli

# Build do binário
go build -o incident.exe
```

## 🛠️ Uso

### Declarar um Incidente

```powershell
.\incident.exe declare --title "Falha na latência da API de Autenticação"
```

**Resultado esperado:**
- Criado diretório `.incidents/INC-20260405-b7x2/`.
- Configurações carregadas de `~/.incident.yaml`.

## ⚙️ Configuração

No primeiro uso, o CLI criará automaticamente um template em seu diretório de usuário (`~/.incident.yaml`). Edite este arquivo para preencher suas credenciais:

```yaml
api_token: "seu_token_aqui"
base_url: "https://api.exemplo.com"
```

## 🏗️ Estrutura do Projeto

- `/cmd`: Comandos Cobra (CLI entrypoints).
- `/internal/config`: Lógica de carregamento e tipagem de configurações.
- `/internal/incident`: Domínio de negócio e manipulação de arquivos do incidente.

## 📄 Licença

Distribuído sob a licença MIT. Veja `LICENSE` para mais informações.

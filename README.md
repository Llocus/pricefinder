## 📄 Notas
Esse projeto é apenas para praticar habilidades com webscrapping e nas tecnologias React.js e Golang, **não serve para uso comercial ou de qualquer outro tipo.**

### 📌 Objetivo do projeto
Esse projeto procura o menor preço entre os produtos de várias lojas diferentes, permitindo que você economize tempo de procurar o produto que você queira com o melhor preço.

## ⚡️ Como rodar o projeto

### 🌟 Frontend:

Recomendo utilizar a versão 18 do node para evitar conflitos de versão. **"nvm use 18"**.

```bash
npm install
```

Não esqueça de configurar o **".env"**, caso não exista, crie o **".env"** com a configuração:

```bash
REACT_APP_SERVER_URL=http://localhost:8000
```

E então rode a aplicação:

```bash
npm run start
```


### 🌟 Backend:

Instale as dependencias necessárias

```bash
go mod download
```

Então pode iniciar a aplicação

```bash
go run main.go
```

### 📌 Demonstrativo

<video controls src="https://raw.githubusercontent.com/Llocus/pricefinder/refs/heads/main/example/pricefinder.mp4" title="Price Finder"></video>
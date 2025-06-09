package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Extension struct {
	XMLName      xml.Name `xml:"extension"`
	Name         string   `xml:"name"`
	Version      string   `xml:"version"`
	Author       string   `xml:"author"`
	AuthorEmail  string   `xml:"authorEmail"`
	AuthorUrl    string   `xml:"authorUrl"`
	License      string   `xml:"license"`
	CreationDate string   `xml:"creationDate"`
	Description  string   `xml:"description"`
	ScriptFile   string   `xml:"scriptfile"`
	Update       struct {
		Schemas []SchemaPath `xml:"schemas>schemapath"`
	} `xml:"update"`
	FileSet struct {
		Files struct {
			Folders []string `xml:"folder"`
			Files   []string `xml:"file"`
		} `xml:"files"`
	} `xml:"fileset"`
	UpdateServers []UpdateServer `xml:"updateservers>server"`
}

type SchemaPath struct {
	Type string `xml:"type,attr"`
	Path string `xml:",chardata"`
}

type UpdateServer struct {
	Name string `xml:"name,attr"`
	URL  string `xml:",chardata"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: ./joomla_explorer https://alvo.com/administrator/manifests/files/joomla.xml")
		return
	}

	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("❌ Erro ao acessar %s\n", url)
		return
	}
	defer resp.Body.Close()

	var ext Extension
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler XML")
		return
	}

	if err := xml.Unmarshal(data, &ext); err != nil {
		fmt.Println("Erro ao parsear XML:", err)
		return
	}

	reportFile, _ := os.Create("report.md")
	wordlistFile, _ := os.Create("wordlist.txt")
	defer reportFile.Close()
	defer wordlistFile.Close()

	write := func(s string) {
		fmt.Print(s)
		reportFile.WriteString(s)
	}

	write("## 🔍 Joomla Manifest Analysis Report\n\n")
	write(fmt.Sprintf("**Nome**: %s\n\n", ext.Name))
	write(fmt.Sprintf("**Versão**: %s\n", ext.Version))
	write(fmt.Sprintf("**Autor**: %s <%s> (%s)\n", ext.Author, ext.AuthorEmail, ext.AuthorUrl))
	write(fmt.Sprintf("**Criado em**: %s\n", ext.CreationDate))
	write(fmt.Sprintf("**Licença**: %s\n\n", ext.License))

	write("### 📁 Folders Expostas:\n")
	for _, folder := range ext.FileSet.Files.Folders {
		write(fmt.Sprintf("- /%s/\n", folder))
		wordlistFile.WriteString("/" + folder + "/\n")
	}
	write("\n")

	write("### 📄 Arquivos Expostos:\n")
	for _, file := range ext.FileSet.Files.Files {
		write(fmt.Sprintf("- /%s\n", file))
		wordlistFile.WriteString("/" + file + "\n")
	}
	write("\n")

	write("### 🔧 Script Instalador Detectado:\n")
	if ext.ScriptFile != "" {
		write(fmt.Sprintf("- %s\n", ext.ScriptFile))
		wordlistFile.WriteString("/" + ext.ScriptFile + "\n")
	}

	write("\n### 🧬 Schemas de Banco:\n")
	for _, schema := range ext.Update.Schemas {
		write(fmt.Sprintf("- (%s): %s\n", schema.Type, schema.Path))
		wordlistFile.WriteString("/" + schema.Path + "\n")
	}

	write("\n### 🌐 Servidores de Update:\n")
	for _, server := range ext.UpdateServers {
		write(fmt.Sprintf("- %s → %s\n", server.Name, server.URL))
	}

	write("\n### 🔥 Avaliação de Criticidade:\n")
	critical := false
	if strings.Contains(ext.ScriptFile, ".php") {
		write("- ⚠️ script.php detectado: pode ser usado para execuções internas.\n")
		critical = true
	}
	for _, f := range ext.FileSet.Files.Folders {
		if f == "tmp" || f == "logs" || f == "cache" {
			write(fmt.Sprintf("- 🔥 Diretório sensível exposto: /%s/\n", f))
			critical = true
		}
	}
	if critical {
		write("\n**🔴 Criticidade Alta: Estrutura interna e arquivos sensíveis expostos.**\n")
	} else {
		write("\n🟡 Criticidade Moderada: Estrutura revelada, mas sem arquivos perigosos acessíveis diretamente.\n")
	}

	write("\n### 🧪 Sugestões de Exploração:\n")
	write("- Tente LFI com /index.php?page=../../administrator/components/...\n")
	write("- Use ffuf com wordlist.txt gerado:\n")
	write("```bash\nffuf -u https://alvo.com/FUZZ -w wordlist.txt -fc 403,404\n```\n")
	write("- Verifique permissões de upload via /tmp, /cache, /logs\n")
	write("- Teste se script.php executa comandos internos via CSRF ou ACL\n")

	write("\n### 🌐 Testando URLs derivadas...\n")
	base := strings.TrimSuffix(url, "/administrator/manifests/files/joomla.xml")
	testPaths := []string{}
	for _, folder := range ext.FileSet.Files.Folders {
		testPaths = append(testPaths, "/"+folder+"/")
	}
	for _, file := range ext.FileSet.Files.Files {
		testPaths = append(testPaths, "/"+file)
	}
	if ext.ScriptFile != "" {
		testPaths = append(testPaths, "/"+ext.ScriptFile)
	}
	for _, schema := range ext.Update.Schemas {
		testPaths = append(testPaths, "/"+schema.Path)
	}

	client := &http.Client{}
	for _, p := range testPaths {
		full := base + p
		req, _ := http.NewRequest("GET", full, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
		line := fmt.Sprintf("- [%d] %s\n", resp.StatusCode, full)
		if resp.StatusCode != 404 {
			fmt.Print(line)
			reportFile.WriteString(line)
		}
	}

	fmt.Println("\n✅ Relatório salvo como: report.md")
	fmt.Println("✅ Wordlist salva como: wordlist.txt")
}

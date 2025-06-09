## 🚀 JoomlaGO

**Automate locating and analyzing Joomla manifest XML files!**

### 🔍 What is the Joomla Manifest XML?

The `joomla.xml` file is the manifest for Joomla components, modules, or plugins, containing metadata such as name, version, author, license, and installation/update instructions.

### 📖 How It Works

1. **Locate:** The XML is usually at `administrator/manifests/files/joomla.xml` in the Joomla installation.
2. **Parse:** Extract key fields: name, version, author, license, directory and file paths.
3. **Report:** Generate a detailed report with emojis highlighting critical items.

### ⚙️ Installation

```bash
# Clone the repository
git clone https://github.com/KingOfBugbounty/jomlaGO/
cd jomlaGO
# Build the binary
go build -o joml main.go
```

### 📦 Usage

```bash
./joml <manifest_path>
# Example:
./joml site.com/administrator/manifests/files/joomla.xml
```

### ✨ Features

* 🔍 **Metadata Extraction:** Name, version, author, license
* 📁 **Exposed Directories & Files:** Lists common folders and files
* 🔧 **Installer Script Detection:** Looks for `script.php`
* 🧬 **Database Schemas:** Identifies SQL update paths
* 🌐 **Update Servers:** Finds core update URLs
* 🔥 **Criticality Assessment:** Flags sensitive directories and scripts

### 💡 Sample Output

```text
## 🔍 Joomla Manifest Analysis Report

**Name**: files_joomla
**Version**: 3.8.6
**Author**: Joomla! Project <admin@joomla.org>
...
**🔴 High Criticality:** Internal structure and sensitive files exposed.
```

### 🤝 Contributing

1. Fork this repo
2. Create a branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to your branch (`git push origin feature/new-feature`)
5. Open a Pull Request

### 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

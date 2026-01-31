# Lucas Kit (`lucaskit`)

> The ultimate domain toolkit containing **UltraDNS**, **SiteStress**, and **UltraCrack**. Made by Lucas Mangroelal | lucasmangroelal.nl

**Lucas Kit** is een collectie van krachtige tools voor DNS/Domain information gathering en security testing. Het bevat:

1. **UltraDNS** (voorheen LucasDNS): Info gathering (DNS, WHOIS, Mail Security, Subdomains).
2. **SiteStress** (voorheen Lucaskill): Advanced HTTP stress test / load test tool.
3. **UltraCrack**: SecLists integrated brute-force tool.

## Install

### Kali Linux / macOS / Linux (aanbevolen)

**Automatische installatie (detecteert architecture):**

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/Lucas_Kit/main/scripts/install.sh | sh
```

Dit installeert `ultradns`, `sitestress`, en `ultracrack` naar `/usr/local/bin` (of `~/.local/bin`).

### Windows

**PowerShell:**

```powershell
.\scripts\install.ps1 -Repo "lucasenlucas/Lucas_Kit"
```

## Tools

### 1. UltraDNS (`ultradns`)

Info gathering tool.

```bash
ultradns -d <domein> [flags]
```

**Features:**
- DNS Records (A, AAAA, MX, NS, TXT, SOA, CAA, SRV)
- Mail Security (SPF, DMARC, DKIM, MTA-STS)
- WHOIS informatie
- Certificate Transparency (Subdomeinen)

**Voorbeelden:**
```bash
ultradns -d example.com -inf -n
ultradns -d example.com -subs
```

### 2. SiteStress (`sitestress`)

HTTP stress/load test tool.

```bash
sitestress -d <domein> -t <minuten> [flags]
```

**Features:**
- High performance HTTP flooding
- "Keep-Down" logic: monitort site en valt automatisch weer aan als hij online komt
- Real-time dashboard
- Reporting (`-o output_dir`)

**Voorbeelden:**
```bash
sitestress -d example.com -t 10
```

### 3. UltraCrack (`ultracrack`)

Brute-force tool met SecLists integratie en auto-detectie.

```bash
ultracrack [flags]
```

**Features:**
- **Auto Analyze**: Vindt automatisch inlogvelden (`--analyze`).
- **SecLists**: Download populaire lijsten (`--dl-seclists`).
- **Supports**: HTTP Basic Auth & HTML Forms.

**Hoe gebruik je het?**

1. **Vind de velden en command:**
   Gebruik `--analyze` op je target pagina. De tool vertelt je welk commando je moet gebruiken!
   ```bash
   ultracrack --analyze -url http://example.com/login
   ```

2. **Run de aanval:**
   Kopieer het gesuggereerde commando, of bouw het zelf:
   ```bash
   ultracrack -u admin -pl top-10000.txt -url http://example.com/login -m form -uf username -pf password -fail-text "Invalid password"
   ```

**Flags:**
- `-m`: Method (`basic` of `form`).
- `-uf`: Username veld naam (bijv. `email`, `user_id`).
- `-pf`: Password veld naam (bijv. `pass`, `pwd`).
- `-fail-text`: Tekst die op de pagina staat als de login faalt (belangrijk voor form mode!).

> **⚠️ DISCLAIMER:** Gebruik deze tools alleen op systemen waar je expliciete toestemming voor hebt.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Additional Notice on Naming, Forks, and Liability

Use of the names "Lucas Mangroelal", "Lucas DNS", "Lucas Kit", or any related project names associated with the original version of this Software does not imply endorsement by the original author.
Any redistributed, modified, or forked versions must make it clear that they are unofficial versions if they are not directly maintained by Lucas Mangroelal.
Lucas Mangroelal is not responsible or liable for any misuse, damages, or consequences resulting from third-party copies, forks, or modified versions of this Software.

For more information, permissions regarding naming, or official inquiries, contact:
kit@lucasmangroelal.nl

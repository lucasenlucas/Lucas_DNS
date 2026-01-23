# Lucas DNS (`lucasdns`)

Cross-platform DNS/Domain info tool voor **Kali Linux**, **macOS** en **Windows** terminals.

## Install

### Kali Linux (aanbevolen)

**Automatische installatie (detecteert automatisch architecture - amd64/arm64):**

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/Lucas_DNS/main/scripts/install.sh | sh
```

De installer:
- ✅ Detecteert automatisch je architecture (amd64 of arm64)
- ✅ Downloadt de juiste binary
- ✅ Installeert naar `/usr/local/bin` (vereist sudo)
- ✅ Test of alles werkt

**Na installatie:**
```bash
lucasdns --help
```

### macOS / Andere Linux distributies

**Automatische installatie:**

```bash
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/Lucas_DNS/main/scripts/install.sh | sh
```

**⚠️ Als je geen sudo hebt:** De installer gebruikt dan `~/.local/bin`. Voeg dit toe aan je PATH:

**Voor zsh (macOS standaard):**
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Voor bash (Linux):**
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Via Go (als je Go geïnstalleerd hebt)

```bash
go install github.com/lucasenlucas/Lucas_DNS@latest
lucasdns --help
```

### Windows

**PowerShell:**
```powershell
.\scripts\install.ps1 -Repo "lucasenlucas/Lucas_DNS"
```

Of download handmatig de `.zip` voor Windows vanaf de [releases pagina](https://github.com/lucasenlucas/Lucas_DNS/releases).

## Gebruik

```bash
lucasdns -d <domein> <flag(s)>
```

Voorbeelden:

```bash
lucasdns -d lucasmangroelal.nl -subs
lucasdns -d lucasmangroelal.nl -inf -n
lucasdns -d lucasmangroelal.nl -whois
```

## Flags (kort)

- `-inf`: info mode (als je verder niks specificeert: DNS + mail checks + WHOIS)
- `-n`: alle DNS records + mail checks (A/AAAA/CNAME/MX/NS/TXT/SOA/CAA/SRV)
- `-whois`: WHOIS (registratie/expiratie/nameservers waar mogelijk)
- `-subs`: subdomeinen (certificate transparency via `crt.sh`)

Record-only:

- `-a`, `-aaaa`, `-cname`, `-mx`, `-ns`, `-txt`, `-soa`, `-caa`, `-srv`

Extra:

- `-r <ip[:port]>`: custom DNS resolver (default: systeem of `8.8.8.8:53`)
- `-timeout 5s`: timeout per query

## Notes

- `-subs` gebruikt certificate transparency (CT). Dit vindt niet “alles”, maar is vaak een goede eerste bron.
- WHOIS parsing verschilt per TLD/registrar; als parsing faalt print `lucasdns` de raw WHOIS.
- `-srv` toont SRV records voor een lijst met **bekende services** (SRV kan je niet “globaal” op één domein opvragen zonder te weten welke `_service._proto` je zoekt).


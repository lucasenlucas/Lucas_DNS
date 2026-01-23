# Lucas DNS (`lucasdns`)

Cross-platform DNS/Domain info tool voor **Kali Linux**, **macOS** en **Windows** terminals.

## Install

### Optie 1 — via Go (makkelijkst)

Als je Go geïnstalleerd hebt:

```bash
<<<<<<< HEAD
go install github.com/lucasenlucas/Lucas_DNS@latest
=======
go install github.com/lucasenlucas/LucasDNS@latest
>>>>>>> c284afc9acca5d209529308c1db34565f2f0d0ec
```

Daarna:

```bash
lucasdns --help
```

### Optie 2 — via GitHub Releases (aanrader voor non-Go users)

<<<<<<< HEAD
Download de laatste release binaries voor jouw platform (Linux/macOS/Windows) vanaf [GitHub Releases](https://github.com/lucasenlucas/Lucas_DNS/releases).
=======
Download de laatste release binaries voor jouw platform (Linux/macOS/Windows) vanaf [GitHub Releases](https://github.com/lucasenlucas/LucasDNS/releases).
>>>>>>> c284afc9acca5d209529308c1db34565f2f0d0ec

**Linux/macOS installer (automatisch):**

```bash
<<<<<<< HEAD
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/Lucas_DNS/main/scripts/install.sh | sh
=======
curl -fsSL https://raw.githubusercontent.com/lucasenlucas/LucasDNS/main/scripts/install.sh | sh
>>>>>>> c284afc9acca5d209529308c1db34565f2f0d0ec
```

Of handmatig:

```bash
<<<<<<< HEAD
REPO="lucasenlucas/Lucas_DNS" sh scripts/install.sh
=======
REPO="lucasenlucas/LucasDNS" sh scripts/install.sh
>>>>>>> c284afc9acca5d209529308c1db34565f2f0d0ec
```

**Windows PowerShell:**

```powershell
<<<<<<< HEAD
.\scripts\install.ps1 -Repo "lucasenlucas/Lucas_DNS"
```

Of download handmatig de `.zip` voor Windows vanaf de [releases pagina](https://github.com/lucasenlucas/Lucas_DNS/releases).
=======
.\scripts\install.ps1 -Repo "lucasenlucas/LucasDNS"
```

Of download handmatig de `.zip` voor Windows vanaf de [releases pagina](https://github.com/lucasenlucas/LucasDNS/releases).
>>>>>>> c284afc9acca5d209529308c1db34565f2f0d0ec

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


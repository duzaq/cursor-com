## 17. Nginx Reverse Proxy

### 17.1 Vis?o Geral
O sistema utiliza Nginx como reverse proxy para rotear requisi??es HTTP/HTTPS dos subdom?nios criados via Cloudflare para os containers Docker correspondentes.

### 17.2 Fluxo de Requisi??o
```
Cliente ? Cloudflare DNS ? Cloudflare Proxy ? Nginx Reverse Proxy ? Container Docker
```

### 17.3 Configura??o Din?mica
- Cada servi?o criado gera automaticamente:
  - Configura??o de upstream (apontando para o container)
  - Server block (com o subdom?nio configurado)
  - Configura??o de SSL/TLS (via Cloudflare)
  - Rate limiting (se configurado)
  - Logs separados por servi?o

### 17.4 Gerenciamento
- Cria??o autom?tica de configura??es
- Atualiza??o sem downtime (nginx -s reload)
- Remo??o autom?tica ao deletar servi?o
- Valida??o antes de aplicar mudan?as

### 17.5 Seguran?a
- Rate limiting por servi?o
- Headers de seguran?a
- SSL/TLS autom?tico via Cloudflare
- Isolamento de logs por servi?o

---

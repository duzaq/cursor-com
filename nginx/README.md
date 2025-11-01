## Nginx Reverse Proxy

O sistema utiliza **Nginx como reverse proxy** para rotear requisi??es baseadas em subdom?nios para os containers Docker correspondentes.

### Por que Nginx Reverse Proxy?

1. **Roteamento por Subdom?nio**: Cada servi?o criado recebe um subdom?nio ?nico (ex: `automation.example.com`, `blog.example.com`)
2. **SSL/TLS**: Integra??o com Cloudflare para SSL autom?tico
3. **Load Balancing**: Distribui??o de carga entre m?ltiplas inst?ncias
4. **Rate Limiting**: Prote??o contra abuso por dom?nio
5. **Cache**: Otimiza??o de performance
6. **Headers de Seguran?a**: Configura??o centralizada de seguran?a

### Arquitetura

```
Internet ? Cloudflare ? Nginx Reverse Proxy ? Container Docker (n8n/WordPress/etc)
                         ?
                    API Manager (porta 8080)
```

### Fluxo de Funcionamento

1. **Cria??o de Servi?o**:
   - Servi?o ? criado via API
   - Container Docker ? iniciado
   - Subdom?nio ? criado no Cloudflare
   - Configura??o do Nginx ? gerada automaticamente
   - Nginx ? recarregado

2. **Requisi??o ao Servi?o**:
   - Cliente acessa `automation.example.com`
   - Cloudflare roteia para Nginx
   - Nginx verifica configura??o do subdom?nio
   - Nginx faz proxy para o container Docker correspondente

3. **Remo??o de Servi?o**:
   - Servi?o ? removido via API
   - Container Docker ? parado
   - Subdom?nio ? removido do Cloudflare
   - Configura??o do Nginx ? removida
   - Nginx ? recarregado

### Estrutura de Arquivos

```
nginx/
??? nginx.conf              # Configura??o principal do Nginx
??? conf.d/
    ??? service-template.conf.example  # Template de exemplo
    ??? automation.conf     # Configura??o gerada para servi?o "automation"
    ??? blog.conf           # Configura??o gerada para servi?o "blog"
    ??? ...                 # Uma configura??o por servi?o
```

### Configura??o Din?mica

O sistema gera automaticamente arquivos de configura??o do Nginx para cada servi?o criado. As configura??es s?o:

- **Geradas automaticamente** quando um servi?o ? criado
- **Atualizadas** quando recursos do servi?o s?o modificados
- **Removidas** quando um servi?o ? deletado
- **Validadas** antes de aplicar mudan?as

### Exemplo de Configura??o Gerada

```nginx
server {
    listen 80;
    server_name automation.example.com;
    
    limit_req zone=service_limit burst=10 nodelay;
    
    location / {
        proxy_pass http://172.17.0.1:5678;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Recarregamento Autom?tico

Quando configura??es s?o alteradas, o Nginx ? recarregado automaticamente usando `nginx -s reload` sem downtime.

### Pr?ximos Passos

- [ ] Integra??o completa com cria??o de servi?os
- [ ] Descoberta autom?tica de IPs dos containers
- [ ] Load balancing entre m?ltiplas inst?ncias
- [ ] Cache configur?vel por servi?o
- [ ] Rate limiting configur?vel por cliente

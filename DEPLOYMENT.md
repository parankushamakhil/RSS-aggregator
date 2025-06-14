# Deployment Guide

## Production Deployment Instructions

### Backend Deployment Options

#### Option 1: Traditional Server Deployment
1. **Server Setup**:
   ```bash
   # Install Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   
   # Install PostgreSQL
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   ```

2. **Database Setup**:
   ```bash
   sudo -u postgres createdb rssaggregator_prod
   sudo -u postgres psql -c "CREATE USER rssapp WITH PASSWORD 'secure_password';"
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE rssaggregator_prod TO rssapp;"
   ```

3. **Application Setup**:
   ```bash
   # Clone and build
   git clone <your-repo>
   cd backend-go
   go build -o rssaggregator
   
   # Create production .env
   cat > .env << EOF
   DB_USER=rssapp
   DB_PASSWORD=secure_password
   DB_NAME=rssaggregator_prod
   DB_HOST=localhost
   DB_PORT=5432
   EOF
   
   # Run with systemd
   sudo cp rssaggregator /usr/local/bin/
   ```

4. **Systemd Service**:
   ```ini
   # /etc/systemd/system/rssaggregator.service
   [Unit]
   Description=RSS Aggregator Backend
   After=network.target postgresql.service
   
   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/opt/rssaggregator
   ExecStart=/usr/local/bin/rssaggregator
   Restart=always
   RestartSec=5
   
   [Install]
   WantedBy=multi-user.target
   ```

#### Option 2: Docker Deployment
1. **Backend Dockerfile**:
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o rssaggregator
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/rssaggregator .
   EXPOSE 8080
   CMD ["./rssaggregator"]
   ```

2. **Docker Compose**:
   ```yaml
   version: '3.8'
   services:
     db:
       image: postgres:14
       environment:
         POSTGRES_DB: rssaggregator
         POSTGRES_USER: rssapp
         POSTGRES_PASSWORD: secure_password
       volumes:
         - postgres_data:/var/lib/postgresql/data
       ports:
         - "5432:5432"
   
     backend:
       build: ./backend-go
       environment:
         DB_USER: rssapp
         DB_PASSWORD: secure_password
         DB_NAME: rssaggregator
         DB_HOST: db
         DB_PORT: 5432
       ports:
         - "8080:8080"
       depends_on:
         - db
   
     frontend:
       build: ./frontend
       ports:
         - "80:80"
       depends_on:
         - backend
   
   volumes:
     postgres_data:
   ```

### Frontend Deployment

#### Option 1: Static Hosting (Recommended)
The frontend is already deployed at: **https://jrejpupe.manus.space**

For your own deployment:
1. Build the React app:
   ```bash
   cd frontend
   npm run build
   ```

2. Deploy to static hosting services:
   - **Netlify**: Drag and drop the `dist` folder
   - **Vercel**: Connect your GitHub repository
   - **AWS S3 + CloudFront**: Upload to S3 bucket with static hosting

#### Option 2: Nginx Deployment
1. **Build and copy**:
   ```bash
   npm run build
   sudo cp -r dist/* /var/www/html/
   ```

2. **Nginx Configuration**:
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       root /var/www/html;
       index index.html;
   
       location / {
           try_files $uri $uri/ /index.html;
       }
   
       location /api {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

### Environment Variables

#### Backend (.env)
```
DB_USER=your_db_user
DB_PASSWORD=your_secure_password
DB_NAME=rssaggregator_prod
DB_HOST=your_db_host
DB_PORT=5432
JWT_SECRET=your_jwt_secret_key
PORT=8080
```

#### Frontend (if needed)
```
VITE_API_URL=https://your-api-domain.com
```

### Security Considerations

1. **Database Security**:
   - Use strong passwords
   - Limit database access to application only
   - Regular backups
   - SSL connections in production

2. **Application Security**:
   - Use HTTPS in production
   - Set secure JWT secret
   - Implement rate limiting
   - Regular security updates

3. **Server Security**:
   - Firewall configuration
   - Regular OS updates
   - SSH key authentication
   - Fail2ban for brute force protection

### Monitoring and Maintenance

1. **Logging**:
   ```bash
   # View application logs
   journalctl -u rssaggregator -f
   
   # Database logs
   sudo tail -f /var/log/postgresql/postgresql-14-main.log
   ```

2. **Health Checks**:
   ```bash
   # Check if backend is running
   curl http://localhost:8080/health
   
   # Check database connection
   psql -h localhost -U rssapp -d rssaggregator_prod -c "SELECT 1;"
   ```

3. **Backup Strategy**:
   ```bash
   # Database backup
   pg_dump -h localhost -U rssapp rssaggregator_prod > backup_$(date +%Y%m%d).sql
   
   # Automated daily backups
   echo "0 2 * * * pg_dump -h localhost -U rssapp rssaggregator_prod > /backups/rss_$(date +\%Y\%m\%d).sql" | crontab -
   ```

### Scaling Considerations

1. **Database Scaling**:
   - Read replicas for heavy read workloads
   - Connection pooling (PgBouncer)
   - Database indexing optimization

2. **Application Scaling**:
   - Load balancer (nginx, HAProxy)
   - Multiple backend instances
   - Redis for session storage

3. **RSS Fetching Optimization**:
   - Separate worker processes
   - Queue system (Redis/RabbitMQ)
   - Distributed fetching

### Troubleshooting

1. **Common Issues**:
   - Database connection errors: Check credentials and network
   - CORS errors: Verify frontend-backend URL configuration
   - RSS fetch failures: Check feed URLs and network connectivity

2. **Performance Issues**:
   - Monitor database query performance
   - Check RSS fetch frequency
   - Optimize frontend bundle size

3. **Security Issues**:
   - Monitor failed login attempts
   - Check for SQL injection attempts
   - Verify JWT token validation


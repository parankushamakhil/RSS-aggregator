# RSS Aggregator - Full Stack Web Application

A modern, full-stack RSS aggregator application with user authentication and database integration. Users can register accounts, add RSS feeds, and view aggregated posts from their subscribed feeds.

## Features

- **User Authentication**: Secure registration and login with JWT tokens
- **Personal Feed Management**: Add, view, and delete RSS feeds per user
- **Automatic RSS Fetching**: Background service fetches RSS posts every 30 minutes
- **Real-time Updates**: Latest posts displayed in chronological order
- **Responsive Design**: Works on desktop and mobile devices
- **Secure Database**: PostgreSQL database with proper user isolation

## Architecture

### Backend (Go)
- **Framework**: Go with Chi router
- **Database**: PostgreSQL with proper schema design
- **Authentication**: JWT tokens with secure password hashing (bcrypt)
- **RSS Processing**: gofeed library for parsing RSS feeds
- **CORS**: Enabled for frontend-backend communication

### Frontend (React)
- **Framework**: React with modern hooks
- **Styling**: CSS with responsive design
- **State Management**: React Context for authentication
- **API Integration**: Fetch API with JWT authentication

### Database Schema
```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Feeds table
CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    last_fetched TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Posts table
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Local Development Setup

### Prerequisites
- Go 1.19 or later
- Node.js 18 or later
- PostgreSQL 14 or later

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend-go
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up PostgreSQL database:
   ```bash
   sudo -u postgres createdb rssaggregator
   sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"
   ```

4. Create `.env` file:
   ```
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=rssaggregator
   DB_HOST=localhost
   DB_PORT=5432
   ```

5. Run the backend:
   ```bash
   go run .
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The frontend will start on `http://localhost:5173`

## API Endpoints

### Authentication
- `POST /register` - Register a new user
- `POST /login` - Login user and receive JWT token

### Feeds Management
- `GET /feeds` - Get user's RSS feeds
- `POST /feeds` - Add a new RSS feed
- `DELETE /feeds/{id}` - Delete a RSS feed

### Posts
- `GET /posts` - Get latest posts from user's feeds

## Deployment

### Frontend Deployment
The frontend is deployed at: **https://jrejpupe.manus.space**

### Backend Deployment
For production deployment, consider:
1. Using environment variables for database configuration
2. Setting up SSL/TLS certificates
3. Using a process manager like systemd or Docker
4. Setting up a reverse proxy with nginx

## Usage

1. **Registration**: Create a new account with username and password
2. **Login**: Sign in with your credentials
3. **Add Feeds**: Use the "Add RSS Feed" form to subscribe to RSS feeds
4. **View Posts**: Latest posts from your feeds appear in the "Latest Posts" section
5. **Manage Feeds**: Delete feeds you no longer want to follow

## Security Features

- Passwords are hashed using bcrypt
- JWT tokens for secure authentication
- SQL injection protection with parameterized queries
- CORS configuration for secure cross-origin requests
- User data isolation (users can only see their own feeds and posts)

## Technical Improvements Made

1. **Database Integration**: Replaced file-based storage with PostgreSQL
2. **User Authentication**: Added secure registration and login system
3. **User Isolation**: Each user has their own feeds and posts
4. **Background Processing**: Automatic RSS fetching every 30 minutes
5. **Modern Frontend**: React-based UI with authentication flow
6. **API Design**: RESTful API with proper HTTP status codes
7. **Error Handling**: Comprehensive error handling throughout the application

## Future Enhancements

- Email notifications for new posts
- Feed categorization and tagging
- Search functionality across posts
- Export/import OPML files
- Mobile app development
- Social features (sharing, comments)

## License

This project is open source and available under the MIT License.


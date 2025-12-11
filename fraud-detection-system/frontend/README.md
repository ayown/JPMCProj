# Fraud Detection Frontend

Modern React/TypeScript frontend for the Banking Fraud Detection System.

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Start development server
npm run dev
```

Visit `http://localhost:3000`

## 📦 Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Redux Toolkit** - State management
- **React Router** - Routing
- **Axios** - HTTP client
- **Recharts** - Data visualization
- **Lucide React** - Icons
- **React Hot Toast** - Notifications

## 🏗️ Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── components/      # React components
│   │   ├── Layout/     # Layout components
│   │   ├── Common/     # Reusable components
│   │   ├── Verification/ # Verification components
│   │   ├── Dashboard/  # Dashboard components
│   │   ├── Reports/    # Report components
│   │   └── Education/  # Education components
│   ├── pages/          # Page components
│   ├── services/       # API services
│   ├── store/          # Redux store
│   ├── hooks/          # Custom hooks
│   ├── types/          # TypeScript types
│   ├── utils/          # Utility functions
│   └── styles/         # Global styles
├── Dockerfile          # Docker configuration
└── nginx.conf          # Nginx configuration
```

## 🎯 Features

### ✅ Implemented
- User authentication (register/login)
- JWT token management
- Message fraud verification
- Real-time fraud detection results
- Verification history
- Statistics dashboard
- Report submission
- Responsive design
- Error handling
- Loading states
- Toast notifications

### 🔄 To Be Implemented
- UI components (in progress)
- Dashboard visualizations
- Advanced filtering
- Export functionality
- Real-time alerts (WebSocket)
- Dark mode
- Multi-language support

## 🔧 Development

### Available Scripts

```bash
npm run dev          # Start dev server (port 3000)
npm run build        # Build for production
npm run preview      # Preview production build
npm run lint         # Run ESLint
npm run type-check   # TypeScript type checking
```

### Environment Variables

Create `.env` file:

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Fraud Detection System
```

## 🐳 Docker

### Build and Run

```bash
# Build image
docker build -t fraud-detection-frontend .

# Run container
docker run -p 80:80 fraud-detection-frontend
```

### Docker Compose

The frontend is integrated in the main `docker-compose.yml`:

```yaml
frontend:
  build: ./frontend
  ports:
    - "3000:80"
  depends_on:
    - api-gateway
```

## 🔗 API Integration

The frontend integrates with the backend API:

### Endpoints
- `POST /auth/register` - Register user
- `POST /auth/login` - Login
- `POST /verify` - Verify message
- `GET /verify/history` - Get history
- `GET /verify/stats` - Get statistics
- `POST /reports` - Submit report

### Features
- Automatic token refresh
- Request interceptors
- Error handling
- Loading states
- Retry logic

## 🎨 Styling

### Tailwind CSS

Custom configuration in `tailwind.config.js`:

```javascript
theme: {
  extend: {
    colors: {
      primary: { /* blue shades */ },
      danger: { /* red shades */ },
      success: { /* green shades */ },
    },
  },
}
```

### Global Styles

Custom CSS classes in `src/styles/globals.css`:

- `.btn` - Button styles
- `.card` - Card container
- `.input` - Form input
- `.label` - Form label

## 🔐 Authentication

### Flow
1. User logs in
2. JWT token stored in localStorage
3. Token added to all API requests
4. Auto-refresh on expiry
5. Redirect to login on auth failure

### Protected Routes

```typescript
<ProtectedRoute>
  <Dashboard />
</ProtectedRoute>
```

## 📱 Responsive Design

Breakpoints:
- Mobile: 320px+
- Tablet: 768px+
- Desktop: 1024px+
- Large: 1440px+

## 🧪 Testing

```bash
# Unit tests (to be added)
npm run test

# E2E tests (to be added)
npm run test:e2e

# Coverage
npm run test:coverage
```

## 🚀 Deployment

### Production Build

```bash
npm run build
# Output in dist/
```

### Nginx Configuration

Production nginx config included in `nginx.conf`:
- Gzip compression
- Security headers
- React Router support
- Static asset caching
- API proxy

## 📚 Documentation

- [API Documentation](../docs/API.md)
- [Architecture](../docs/ARCHITECTURE.md)
- [Deployment Guide](../docs/DEPLOYMENT.md)

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Run linter and type check
4. Test thoroughly
5. Submit PR

## 📄 License

MIT License - see LICENSE file


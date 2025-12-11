import { Link } from 'react-router-dom'

function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="text-6xl font-bold text-danger-600">404</div>
        <div className="mt-4 text-gray-700">Page not found</div>
        <Link to="/" className="btn btn-primary mt-6">Go Home</Link>
      </div>
    </div>
  )
}

export default NotFound

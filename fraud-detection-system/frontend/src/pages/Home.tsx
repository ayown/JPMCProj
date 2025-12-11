import { Link } from 'react-router-dom'

function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-primary-50 to-gray-50">
      <div className="max-w-2xl text-center px-6">
        <h1 className="text-4xl font-bold text-primary-700">AI Fraud Detection</h1>
        <p className="mt-4 text-gray-700">Verify suspicious messages, view risk analysis and stay protected from phishing, vishing and KYC fraud.</p>
        <div className="mt-8 flex gap-4 justify-center">
          <Link to="/verify" className="btn btn-primary">Verify a Message</Link>
          <Link to="/login" className="btn btn-secondary">Login</Link>
        </div>
      </div>
    </div>
  )
}

export default Home

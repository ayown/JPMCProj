import { useEffect } from 'react'
import { useVerification } from '@/hooks/useVerification'

function Dashboard() {
  const { stats, history, loadStats, loadHistory } = useVerification()
  useEffect(()=>{ loadStats(); loadHistory(10,0) },[])
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="card"><div className="text-sm text-gray-600">Total</div><div className="text-2xl font-semibold">{stats?.total_verifications ?? 0}</div></div>
        <div className="card"><div className="text-sm text-gray-600">Fraud Detected</div><div className="text-2xl font-semibold">{stats?.fraud_detected ?? 0}</div></div>
        <div className="card"><div className="text-sm text-gray-600">Fraud Rate</div><div className="text-2xl font-semibold">{stats? Math.round(stats.fraud_rate*100):0}%</div></div>
        <div className="card"><div className="text-sm text-gray-600">Avg Score</div><div className="text-2xl font-semibold">{stats? Math.round(stats.avg_fraud_score*100):0}%</div></div>
      </div>
      <div className="card">
        <div className="text-lg font-semibold mb-3">Recent Verifications</div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="text-left">
                <th className="py-2 pr-4">Verified At</th>
                <th className="py-2 pr-4">Fraud</th>
                <th className="py-2 pr-4">Score</th>
                <th className="py-2 pr-4">Risk</th>
              </tr>
            </thead>
            <tbody>
              {(history ?? []).map(v=> (
                <tr key={v.id} className="border-t">
                  <td className="py-2 pr-4">{new Date(v.verified_at).toLocaleString()}</td>
                  <td className={`py-2 pr-4 ${v.is_fraud?'text-danger-600':'text-success-600'}`}>{v.is_fraud?'Yes':'No'}</td>
                  <td className="py-2 pr-4">{Math.round(v.fraud_score*100)}%</td>
                  <td className="py-2 pr-4">{v.risk_level}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

export default Dashboard

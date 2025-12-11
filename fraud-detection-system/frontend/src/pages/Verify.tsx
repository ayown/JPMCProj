import { useState } from 'react'
import { useVerification } from '@/hooks/useVerification'

function Verify() {
  const { verify, currentVerification, clearVerification } = useVerification()
  const [content,setContent]=useState('')
  const [sender,setSender]=useState('')
  const [messageType,setMessageType]=useState<'SMS'|'WhatsApp'|'Email'>('SMS')
  const submit=async(e:React.FormEvent)=>{
    e.preventDefault()
    await verify({content, sender_header: sender, message_type: messageType})
  }
  return (
    <div className="max-w-3xl mx-auto py-8">
      <h2 className="text-2xl font-semibold mb-4">Verify a Message</h2>
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="label">Message Content</label>
          <textarea className="input h-32" value={content} onChange={e=>setContent(e.target.value)} required/>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="label">Sender Header</label>
            <input className="input" value={sender} onChange={e=>setSender(e.target.value)} required/>
          </div>
          <div>
            <label className="label">Message Type</label>
            <select className="input" value={messageType} onChange={e=>setMessageType(e.target.value as any)}>
              <option value="SMS">SMS</option>
              <option value="WhatsApp">WhatsApp</option>
              <option value="Email">Email</option>
            </select>
          </div>
        </div>
        <div className="flex gap-3">
          <button className="btn btn-primary" type="submit">Analyze</button>
          <button className="btn btn-secondary" type="button" onClick={clearVerification}>Clear</button>
        </div>
      </form>
      {currentVerification && (
        <div className="card mt-6">
          <div className="flex items-center justify-between">
            <div className="text-lg font-semibold">Result</div>
            <div className="text-sm text-gray-600">Risk: {currentVerification.risk_level}</div>
          </div>
          <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-600">Fraud</div>
              <div className={`text-xl font-bold ${currentVerification.is_fraud?'text-danger-600':'text-success-600'}`}>{currentVerification.is_fraud?'Yes':'No'}</div>
            </div>
            <div>
              <div className="text-sm text-gray-600">Score</div>
              <div className="text-xl font-bold">{(currentVerification.fraud_score*100).toFixed(0)}%</div>
            </div>
          </div>
          <div className="mt-4">
            <div className="text-sm text-gray-600">Explanation</div>
            <div className="mt-2 text-gray-800">{currentVerification.explanation}</div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Verify

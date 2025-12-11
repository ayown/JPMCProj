import { useState } from 'react'
import { useAuth } from '@/hooks/useAuth'

function Login() {
  const { login, register } = useAuth()
  const [mode,setMode]=useState<'login'|'register'>('login')
  const [email,setEmail]=useState('')
  const [password,setPassword]=useState('')
  const [fullName,setFullName]=useState('')
  const [phone,setPhone]=useState('')

  const submit=async(e:React.FormEvent)=>{
    e.preventDefault()
    if(mode==='login'){
      await login({email,password})
    }else{
      await register({email,password,full_name:fullName,phone_number:phone})
      setMode('login')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="card w-full max-w-md">
        <h2 className="text-2xl font-semibold mb-6">{mode==='login'?'Login':'Register'}</h2>
        <form onSubmit={submit} className="space-y-4">
          <div>
            <label className="label">Email</label>
            <input className="input" type="email" value={email} onChange={e=>setEmail(e.target.value)} required/>
          </div>
          <div>
            <label className="label">Password</label>
            <input className="input" type="password" value={password} onChange={e=>setPassword(e.target.value)} required/>
          </div>
          {mode==='register' && (
            <>
              <div>
                <label className="label">Full Name</label>
                <input className="input" value={fullName} onChange={e=>setFullName(e.target.value)} required/>
              </div>
              <div>
                <label className="label">Phone Number</label>
                <input className="input" value={phone} onChange={e=>setPhone(e.target.value)} required/>
              </div>
            </>
          )}
          <button className="btn btn-primary w-full" type="submit">{mode==='login'?'Login':'Create Account'}</button>
        </form>
        <div className="mt-4 text-sm text-gray-600">
          {mode==='login'? (
            <button className="text-primary-700" onClick={()=>setMode('register')}>Create an account</button>
          ) : (
            <button className="text-primary-700" onClick={()=>setMode('login')}>Already have an account? Login</button>
          )}
        </div>
      </div>
    </div>
  )
}

export default Login

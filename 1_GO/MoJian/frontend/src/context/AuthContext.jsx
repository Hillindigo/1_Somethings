import { createContext, useContext, useState, useEffect } from 'react'
import client from '../api/client'

const AuthContext = createContext(null)

// AuthProvider 提供全局认证状态管理
export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  // 初始化：从 localStorage 恢复登录状态
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      client.get('/user/profile')
        .then((res) => setUser(res.data))
        .catch(() => localStorage.removeItem('token'))
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [])

  // 登录：保存 Token 并获取用户信息
  const login = async (username, password, captchaId, captchaCode) => {
    const res = await client.post('/auth/login', { username, password, captcha_id: captchaId, captcha_code: captchaCode })
    localStorage.setItem('token', res.data.token)
    setUser(res.data.user)
    return res.data
  }

  // 注册
  const register = async (username, password, email, captchaId, captchaCode) => {
    const res = await client.post('/auth/register', { username, password, email, captcha_id: captchaId, captcha_code: captchaCode })
    return res.data
  }

  // 登出
  const logout = () => {
    localStorage.removeItem('token')
    setUser(null)
  }

  const isAdmin = user?.role === 1

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout, isAdmin }}>
      {children}
    </AuthContext.Provider>
  )
}

// useAuth Hook：在组件中获取认证上下文
export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

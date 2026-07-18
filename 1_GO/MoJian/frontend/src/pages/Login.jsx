import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import client from '../api/client'

// Login 登录页面，集成图片验证码
export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ username: '', password: '', captchaCode: '' })
  const [captchaId, setCaptchaId] = useState('')
  const [captchaImage, setCaptchaImage] = useState('')
  const [error, setError] = useState('')

  // fetchCaptcha 获取验证码图片
  const fetchCaptcha = async () => {
    try {
      const res = await client.get('/auth/captcha')
      setCaptchaId(res.data.captcha_id)
      setCaptchaImage(res.data.image)
    } catch {
      setError('获取验证码失败，请刷新页面重试')
    }
  }

  // 组件挂载时获取验证码
  useEffect(() => {
    fetchCaptcha()
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      await login(form.username, form.password, captchaId, form.captchaCode)
      navigate('/')
    } catch (err) {
      // 登录失败后刷新验证码
      fetchCaptcha()
      setForm({ ...form, captchaCode: '' })
      setError(err.message || '登录失败')
    }
  }

  return (
    <div className="max-w-md mx-auto px-6 pt-30">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-2">
          登录<span className="text-accent">.</span>
        </h1>
        <p className="text-muted mb-10">回到你的文字世界</p>

        {error && (
          <div className="bg-accent/10 text-accent px-4 py-3 rounded mb-6 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
              用户名
            </label>
            <input
              type="text"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              required
              className="w-full border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
              placeholder="输入用户名"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
              密码
            </label>
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              required
              className="w-full border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
              placeholder="输入密码"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
              验证码
            </label>
            <div className="flex items-center gap-3">
              <input
                type="text"
                value={form.captchaCode}
                onChange={(e) => setForm({ ...form, captchaCode: e.target.value })}
                required
                maxLength={4}
                className="flex-1 border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
                placeholder="输入验证码"
              />
              {captchaImage && (
                <img
                  src={captchaImage}
                  alt="验证码"
                  onClick={fetchCaptcha}
                  className="h-12 rounded cursor-pointer border border-border hover:opacity-80 transition-opacity"
                  title="点击刷新验证码"
                />
              )}
            </div>
            <p className="text-xs text-muted mt-1">点击图片可刷新验证码</p>
          </div>
          <button
            type="submit"
            className="w-full bg-ink text-paper py-3 rounded-lg text-sm font-semibold tracking-wide hover:bg-accent transition-colors duration-200"
          >
            登录
          </button>
        </form>

        <p className="text-center text-sm text-muted mt-8">
          还没有账号？{' '}
          <Link to="/register" className="text-accent hover:underline">立即注册</Link>
        </p>
      </div>
    </div>
  )
}

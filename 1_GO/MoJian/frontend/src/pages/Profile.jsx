import { useState } from 'react'
import { useAuth } from '../context/AuthContext'
import client from '../api/client'

// Profile 个人信息页面
export default function Profile() {
  const { user, logout } = useAuth()
  const [email, setEmail] = useState(user?.email || '')
  const [avatar, setAvatar] = useState(user?.avatar || '')
  const [saved, setSaved] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await client.put('/user/profile', { email, avatar })
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } catch (err) {
      // 错误已在拦截器中处理
    }
  }

  return (
    <div className="max-w-md mx-auto px-6 pt-18">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-2">
          个人设置<span className="text-accent">.</span>
        </h1>
        <p className="text-muted mb-10">{user?.username}</p>

        {saved && (
          <div className="bg-green-50 text-green-700 px-4 py-3 rounded mb-6 text-sm">
            保存成功
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
              邮箱
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full border border-border rounded-lg px-4 py-3 bg-white/50 text-ink focus:outline-none focus:border-accent transition-colors"
            />
          </div>
          <div>
            <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
              头像 URL
            </label>
            <input
              type="url"
              value={avatar}
              onChange={(e) => setAvatar(e.target.value)}
              className="w-full border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
              placeholder="https://..."
            />
          </div>
          <button
            type="submit"
            className="bg-ink text-paper px-8 py-3 rounded-lg text-sm font-semibold tracking-wide hover:bg-accent transition-colors duration-200"
          >
            保存
          </button>
        </form>
      </div>
    </div>
  )
}

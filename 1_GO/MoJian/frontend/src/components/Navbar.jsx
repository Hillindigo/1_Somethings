import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

// Navbar 顶部导航栏，包含品牌标识、导航链接和用户状态
export default function Navbar() {
  const { user, logout, isAdmin } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const isActive = (path) => location.pathname === path

  const navLink = (to, label) => (
    <Link
      to={to}
      className={`text-sm tracking-wide uppercase transition-colors duration-200 ${
        isActive(to) ? 'text-accent font-semibold' : 'text-muted hover:text-ink'
      }`}
    >
      {label}
    </Link>
  )

  return (
    <header className="border-b border-border">
      <nav className="max-w-5xl mx-auto px-6 py-5 flex items-baseline justify-between">
        {/* 品牌：墨笺 MoJian 标识 */}
        <Link to="/" className="font-display text-2xl font-bold tracking-tight text-ink">
          墨笺
        </Link>

        {/* 导航链接 */}
        <div className="flex items-baseline gap-8">
          {navLink('/', '文章')}
          {user && navLink('/write', '撰写')}
          {user && navLink('/admin/categories', '分类')}
          {isAdmin && navLink('/admin', '管理')}
          {user ? (
            <div className="flex items-baseline gap-6">
              {navLink('/profile', user.username)}
              <button
                onClick={() => { logout(); navigate('/') }}
                className="text-sm tracking-wide uppercase text-muted hover:text-accent transition-colors"
              >
                退出
              </button>
            </div>
          ) : (
            <div className="flex items-baseline gap-6">
              {navLink('/login', '登录')}
              {navLink('/register', '注册')}
            </div>
          )}
        </div>
      </nav>
    </header>
  )
}

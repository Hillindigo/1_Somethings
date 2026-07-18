import { Link } from 'react-router-dom'

// Dashboard 管理后台首页
export default function Dashboard() {
  return (
    <div className="max-w-5xl mx-auto px-6 pt-18 pb-10">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-10">
          管理<span className="text-accent">.</span>
        </h1>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Link
            to="/admin/users"
            className="ink-bar py-8 bg-white/30 rounded-r-lg hover:bg-white/60 transition-colors duration-200"
          >
            <h2 className="font-display text-2xl font-bold mb-2">用户管理</h2>
            <p className="text-muted text-sm">管理用户角色和账号</p>
          </Link>
          <Link
            to="/admin/categories"
            className="ink-bar py-8 bg-white/30 rounded-r-lg hover:bg-white/60 transition-colors duration-200"
          >
            <h2 className="font-display text-2xl font-bold mb-2">分类管理</h2>
            <p className="text-muted text-sm">创建、编辑和删除文章分类</p>
          </Link>
          <Link
            to="/admin/tags"
            className="ink-bar py-8 bg-white/30 rounded-r-lg hover:bg-white/60 transition-colors duration-200"
          >
            <h2 className="font-display text-2xl font-bold mb-2">标签管理</h2>
            <p className="text-muted text-sm">创建、编辑和删除文章标签</p>
          </Link>
        </div>
      </div>
    </div>
  )
}

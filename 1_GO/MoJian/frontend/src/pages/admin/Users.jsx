import { useState, useEffect } from 'react'
import client from '../../api/client'

// Users 用户管理页面（管理员），管理所有用户信息
export default function Users() {
  const [users, setUsers] = useState([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const pageSize = 10

  // 加载用户列表
  const loadUsers = async () => {
    const res = await client.get('/admin/users', { params: { page, page_size: pageSize } })
    setUsers(res.data?.list || [])
    setTotal(res.data?.total || 0)
  }

  useEffect(() => { loadUsers() }, [page])

  // 切换用户角色
  const handleToggleRole = async (user) => {
    const newRole = user.role === 1 ? 0 : 1
    const label = newRole === 1 ? '管理员' : '普通用户'
    if (!confirm(`确定将用户 ${user.username} 设为${label}？`)) return
    await client.put(`/admin/users/${user.id}/role`, { role: newRole })
    loadUsers()
  }

  // 删除用户
  const handleDelete = async (user) => {
    if (!confirm(`确定删除用户 ${user.username}？此操作不可恢复！`)) return
    await client.delete(`/admin/users/${user.id}`)
    loadUsers()
  }

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="max-w-5xl mx-auto px-6 pt-18 pb-10">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-10">
          用户管理<span className="text-accent">.</span>
        </h1>

        {/* 用户列表 */}
        <div className="divide-y divide-border">
          {users.map((user) => (
            <div key={user.id} className="flex items-center justify-between py-5">
              <div className="flex items-center gap-4">
                {/* 头像占位 */}
                <div className="w-10 h-10 rounded-full bg-accent/10 flex items-center justify-center text-accent font-bold text-sm">
                  {user.username?.[0]?.toUpperCase() || '?'}
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-display text-lg font-bold">{user.username}</span>
                    <span
                      className={`text-xs px-2 py-0.5 rounded-full ${
                        user.role === 1
                          ? 'bg-accent/10 text-accent border border-accent/30'
                          : 'bg-muted/10 text-muted border border-border'
                      }`}
                    >
                      {user.role === 1 ? '管理员' : '普通用户'}
                    </span>
                  </div>
                  <p className="text-sm text-muted mt-0.5">{user.email}</p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <button
                  onClick={() => handleToggleRole(user)}
                  className="text-sm text-muted hover:text-accent transition-colors"
                >
                  {user.role === 1 ? '设为普通用户' : '设为管理员'}
                </button>
                <button
                  onClick={() => handleDelete(user)}
                  className="text-sm text-muted hover:text-red-500 transition-colors"
                >
                  删除
                </button>
              </div>
            </div>
          ))}
          {users.length === 0 && (
            <p className="text-muted italic py-8 text-center font-display">暂无用户</p>
          )}
        </div>

        {/* 分页 */}
        {totalPages > 1 && (
          <div className="flex items-center justify-center gap-4 mt-12">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page <= 1}
              className="text-sm px-5 py-2 border border-border rounded hover:border-ink disabled:opacity-30 disabled:hover:border-border transition-colors"
            >
              上一页
            </button>
            <span className="text-sm text-muted">{page} / {totalPages}</span>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="text-sm px-5 py-2 border border-border rounded hover:border-ink disabled:opacity-30 disabled:hover:border-border transition-colors"
            >
              下一页
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

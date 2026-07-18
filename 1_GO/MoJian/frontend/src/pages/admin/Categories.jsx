import { useState, useEffect } from 'react'
import client from '../../api/client'

// Categories 分类管理页面
export default function Categories() {
  const [categories, setCategories] = useState([])
  const [name, setName] = useState('')
  const [editingId, setEditingId] = useState(null)
  const [editName, setEditName] = useState('')

  // 加载分类列表
  const loadCategories = async () => {
    const res = await client.get('/categories')
    setCategories(res.data || [])
  }

  useEffect(() => { loadCategories() }, [])

  // 创建分类
  const handleCreate = async (e) => {
    e.preventDefault()
    if (!name.trim()) return
    await client.post('/categories', { name: name.trim() })
    setName('')
    loadCategories()
  }

  // 更新分类
  const handleUpdate = async (id) => {
    if (!editName.trim()) return
    await client.put(`/categories/${id}`, { name: editName.trim() })
    setEditingId(null)
    loadCategories()
  }

  // 删除分类
  const handleDelete = async (id) => {
    if (!confirm('确定删除该分类？')) return
    await client.delete(`/categories/${id}`)
    loadCategories()
  }

  return (
    <div className="max-w-3xl mx-auto px-6 pt-18 pb-10">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-10">
          分类管理<span className="text-accent">.</span>
        </h1>

        {/* 新建分类表单 */}
        <form onSubmit={handleCreate} className="flex gap-3 mb-10">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="新分类名称"
            className="flex-1 border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
          />
          <button
            type="submit"
            className="bg-accent text-paper px-6 py-3 rounded-lg text-sm font-semibold tracking-wide hover:bg-ink transition-colors duration-200"
          >
            创建
          </button>
        </form>

        {/* 分类列表 */}
        <div className="divide-y divide-border">
          {categories.map((cat) => (
            <div key={cat.id} className="flex items-center justify-between py-4">
              {editingId === cat.id ? (
                <div className="flex gap-3 flex-1">
                  <input
                    type="text"
                    value={editName}
                    onChange={(e) => setEditName(e.target.value)}
                    className="flex-1 border border-accent rounded px-3 py-1.5 text-ink focus:outline-none"
                    autoFocus
                  />
                  <button
                    onClick={() => handleUpdate(cat.id)}
                    className="text-sm text-accent hover:underline"
                  >
                    保存
                  </button>
                  <button
                    onClick={() => setEditingId(null)}
                    className="text-sm text-muted hover:underline"
                  >
                    取消
                  </button>
                </div>
              ) : (
                <>
                  <span className="font-display text-lg font-bold">{cat.name}</span>
                  <div className="flex gap-4">
                    <button
                      onClick={() => { setEditingId(cat.id); setEditName(cat.name) }}
                      className="text-sm text-muted hover:text-ink transition-colors"
                    >
                      编辑
                    </button>
                    <button
                      onClick={() => handleDelete(cat.id)}
                      className="text-sm text-muted hover:text-accent transition-colors"
                    >
                      删除
                    </button>
                  </div>
                </>
              )}
            </div>
          ))}
          {categories.length === 0 && (
            <p className="text-muted italic py-8 text-center font-display">暂无分类</p>
          )}
        </div>
      </div>
    </div>
  )
}

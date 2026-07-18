import { useState, useEffect } from 'react'
import client from '../../api/client'

// Tags 标签管理页面
export default function Tags() {
  const [tags, setTags] = useState([])
  const [name, setName] = useState('')
  const [editingId, setEditingId] = useState(null)
  const [editName, setEditName] = useState('')

  // 加载标签列表
  const loadTags = async () => {
    const res = await client.get('/tags')
    setTags(res.data || [])
  }

  useEffect(() => { loadTags() }, [])

  // 创建标签
  const handleCreate = async (e) => {
    e.preventDefault()
    if (!name.trim()) return
    await client.post('/admin/tags', { name: name.trim() })
    setName('')
    loadTags()
  }

  // 更新标签
  const handleUpdate = async (id) => {
    if (!editName.trim()) return
    await client.put(`/admin/tags/${id}`, { name: editName.trim() })
    setEditingId(null)
    loadTags()
  }

  // 删除标签
  const handleDelete = async (id) => {
    if (!confirm('确定删除该标签？')) return
    await client.delete(`/admin/tags/${id}`)
    loadTags()
  }

  return (
    <div className="max-w-3xl mx-auto px-6 pt-18 pb-10">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-10">
          标签管理<span className="text-accent">.</span>
        </h1>

        {/* 新建标签表单 */}
        <form onSubmit={handleCreate} className="flex gap-3 mb-10">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="新标签名称"
            className="flex-1 border border-border rounded-lg px-4 py-3 bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none focus:border-accent transition-colors"
          />
          <button
            type="submit"
            className="bg-accent text-paper px-6 py-3 rounded-lg text-sm font-semibold tracking-wide hover:bg-ink transition-colors duration-200"
          >
            创建
          </button>
        </form>

        {/* 标签列表 */}
        <div className="flex flex-wrap gap-3">
          {tags.map((tag) => (
            editingId === tag.id ? (
              <div key={tag.id} className="flex items-center gap-2">
                <input
                  type="text"
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  className="border border-accent rounded px-3 py-1.5 text-sm text-ink focus:outline-none"
                  autoFocus
                />
                <button onClick={() => handleUpdate(tag.id)} className="text-xs text-accent hover:underline">保存</button>
                <button onClick={() => setEditingId(null)} className="text-xs text-muted hover:underline">取消</button>
              </div>
            ) : (
              <div key={tag.id} className="group flex items-center gap-2 border border-border rounded-full px-4 py-2 hover:border-muted transition-colors">
                <span className="text-sm">{tag.name}</span>
                <button
                  onClick={() => { setEditingId(tag.id); setEditName(tag.name) }}
                  className="text-xs text-muted opacity-0 group-hover:opacity-100 hover:text-ink transition-all"
                >
                  编辑
                </button>
                <button
                  onClick={() => handleDelete(tag.id)}
                  className="text-xs text-muted opacity-0 group-hover:opacity-100 hover:text-accent transition-all"
                >
                  删除
                </button>
              </div>
            )
          ))}
          {tags.length === 0 && (
            <p className="text-muted italic py-8 w-full text-center font-display">暂无标签</p>
          )}
        </div>
      </div>
    </div>
  )
}

import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import client from '../api/client'

// WriteArticle 撰写/编辑文章页面，支持 Markdown 实时预览
export default function WriteArticle() {
  const { id } = useParams() // 有 id 则为编辑模式
  const navigate = useNavigate()
  const isEdit = !!id

  const [form, setForm] = useState({
    title: '',
    content: '',
    summary: '',
    status: 0,
    category_id: 0,
    tag_ids: [],
  })
  const [categories, setCategories] = useState([])
  const [tags, setTags] = useState([])
  const [showPreview, setShowPreview] = useState(false)
  const [newCategoryName, setNewCategoryName] = useState('')
  const [showNewCategory, setShowNewCategory] = useState(false)

  // 加载分类和标签
  useEffect(() => {
    client.get('/categories').then((res) => setCategories(res.data || []))
    client.get('/tags').then((res) => setTags(res.data || []))
  }, [])

  // 内联创建分类
  const handleCreateCategory = async () => {
    if (!newCategoryName.trim()) return
    try {
      const res = await client.post('/categories', { name: newCategoryName.trim() })
      const created = res.data
      setCategories((prev) => [...prev, created])
      setForm((prev) => ({ ...prev, category_id: created.id }))
      setNewCategoryName('')
      setShowNewCategory(false)
    } catch (err) {
      alert(err.message || '创建分类失败')
    }
  }

  // 编辑模式：加载已有文章数据
  useEffect(() => {
    if (isEdit) {
      client.get(`/articles/${id}`).then((res) => {
        const a = res.data
        setForm({
          title: a.title,
          content: a.content,
          summary: a.summary || '',
          status: a.status,
          category_id: a.category_id || 0,
          tag_ids: a.tags?.map((t) => t.id) || [],
        })
      })
    }
  }, [id, isEdit])

  // 切换标签选中
  const toggleTag = (tagId) => {
    setForm((prev) => ({
      ...prev,
      tag_ids: prev.tag_ids.includes(tagId)
        ? prev.tag_ids.filter((t) => t !== tagId)
        : [...prev.tag_ids, tagId],
    }))
  }

  // 提交文章
  const handleSubmit = async (status) => {
    const payload = { ...form, status }
    try {
      if (isEdit) {
        await client.put(`/articles/${id}`, payload)
      } else {
        await client.post('/articles', payload)
      }
      navigate('/')
    } catch (err) {
      // 错误已在拦截器中处理
    }
  }

  return (
    <div className="max-w-4xl mx-auto px-6 pt-18 pb-10">
      <div className="article-enter">
        <h1 className="font-display text-5xl font-black mb-10">
          {isEdit ? '编辑' : '撰写'}<span className="text-accent">.</span>
        </h1>

        <div className="space-y-8">
          {/* 标题 */}
          <div>
            <input
              type="text"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              placeholder="文章标题"
              className="w-full font-display text-3xl md:text-4xl font-bold bg-transparent border-b-2 border-border pb-3 text-ink placeholder:text-muted/30 focus:outline-none focus:border-accent transition-colors"
            />
          </div>

          {/* 摘要 */}
          <div>
            <input
              type="text"
              value={form.summary}
              onChange={(e) => setForm({ ...form, summary: e.target.value })}
              placeholder="文章摘要（选填）"
              className="w-full bg-transparent border-b border-border pb-2 text-lg text-ink placeholder:text-muted/30 focus:outline-none focus:border-accent transition-colors"
            />
          </div>

          {/* 分类 + 标签 */}
          <div className="flex flex-wrap gap-6">
            <div>
              <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
                分类
              </label>
              <div className="flex items-center gap-2">
                <select
                  value={form.category_id}
                  onChange={(e) => setForm({ ...form, category_id: Number(e.target.value) })}
                  className="border border-border rounded-lg px-4 py-2.5 bg-white/50 text-ink focus:outline-none focus:border-accent transition-colors"
                >
                  <option value={0}>无分类</option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
                {/* 内联创建分类按钮 */}
                {showNewCategory ? (
                  <div className="flex items-center gap-2">
                    <input
                      type="text"
                      value={newCategoryName}
                      onChange={(e) => setNewCategoryName(e.target.value)}
                      onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), handleCreateCategory())}
                      placeholder="分类名称"
                      autoFocus
                      className="border border-accent rounded-lg px-3 py-2 text-sm bg-white/50 text-ink placeholder:text-muted/40 focus:outline-none transition-colors w-28"
                    />
                    <button
                      onClick={handleCreateCategory}
                      className="text-xs text-accent hover:underline whitespace-nowrap"
                    >
                      确定
                    </button>
                    <button
                      onClick={() => { setShowNewCategory(false); setNewCategoryName('') }}
                      className="text-xs text-muted hover:underline"
                    >
                      取消
                    </button>
                  </div>
                ) : (
                  <button
                    onClick={() => setShowNewCategory(true)}
                    className="text-xs text-muted hover:text-accent transition-colors whitespace-nowrap"
                    title="新建分类"
                  >
                    + 新建
                  </button>
                )}
              </div>
            </div>
            <div>
              <label className="block text-xs font-semibold tracking-widest uppercase text-muted mb-2">
                标签
              </label>
              <div className="flex flex-wrap gap-2">
                {tags.map((tag) => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => toggleTag(tag.id)}
                    className={`text-xs px-3 py-1.5 rounded-full border transition-all duration-200 ${
                      form.tag_ids.includes(tag.id)
                        ? 'bg-accent text-paper border-accent'
                        : 'border-border text-muted hover:text-ink hover:border-muted'
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* 预览切换 */}
          <div className="flex gap-4 border-b border-border pb-2">
            <button
              onClick={() => setShowPreview(false)}
              className={`text-sm font-semibold tracking-wide pb-1 transition-colors ${
                !showPreview ? 'text-accent border-b-2 border-accent' : 'text-muted hover:text-ink'
              }`}
            >
              编辑
            </button>
            <button
              onClick={() => setShowPreview(true)}
              className={`text-sm font-semibold tracking-wide pb-1 transition-colors ${
                showPreview ? 'text-accent border-b-2 border-accent' : 'text-muted hover:text-ink'
              }`}
            >
              预览
            </button>
          </div>

          {/* 内容区 */}
          {showPreview ? (
            <div className="prose-ink min-h-[400px] border border-border rounded-lg p-6 bg-white/30">
              {form.content || <p className="text-muted/30 italic">暂无内容</p>}
            </div>
          ) : (
            <textarea
              value={form.content}
              onChange={(e) => setForm({ ...form, content: e.target.value })}
              placeholder="用 Markdown 写下你的文章..."
              rows={20}
              className="w-full border border-border rounded-lg p-6 bg-white/30 text-ink font-mono text-sm leading-relaxed placeholder:text-muted/30 focus:outline-none focus:border-accent transition-colors resize-y"
            />
          )}

          {/* 操作按钮 */}
          <div className="flex gap-4 pt-4">
            <button
              onClick={() => handleSubmit(1)}
              className="bg-accent text-paper px-8 py-3 rounded-lg text-sm font-semibold tracking-wide hover:bg-ink transition-colors duration-200"
            >
              发布文章
            </button>
            <button
              onClick={() => handleSubmit(0)}
              className="border border-border text-muted px-8 py-3 rounded-lg text-sm font-semibold tracking-wide hover:border-ink hover:text-ink transition-colors duration-200"
            >
              保存草稿
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

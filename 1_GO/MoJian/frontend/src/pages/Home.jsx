import { useState, useEffect } from 'react'
import client from '../api/client'
import ArticleCard from '../components/ArticleCard'

// Home 首页，展示已发布文章列表，带分类筛选
export default function Home() {
  const [articles, setArticles] = useState([])
  const [categories, setCategories] = useState([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [selectedCategory, setSelectedCategory] = useState(null)
  const pageSize = 10

  // 加载文章列表
  useEffect(() => {
    const params = { page, page_size: pageSize, status: 1 }
    if (selectedCategory) params.category_id = selectedCategory
    client.get('/articles', { params }).then((res) => {
      setArticles(res.data.list || [])
      setTotal(res.data.total || 0)
    })
  }, [page, selectedCategory])

  // 加载分类列表
  useEffect(() => {
    client.get('/categories').then((res) => setCategories(res.data || []))
  }, [])

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="max-w-5xl mx-auto px-6 pt-18 pb-10">
      {/* 页面标题区 */}
      <div className="mb-16 article-enter">
        <h1 className="font-display text-6xl md:text-8xl font-black leading-none tracking-tight">
          最新<span className="text-accent">文章</span>
        </h1>
        <p className="text-muted mt-4 text-lg">以墨为记，记录思考与灵感</p>
      </div>

      {/* 分类筛选条 */}
      {categories.length > 0 && (
        <div className="flex flex-wrap gap-3 mb-12 article-enter" style={{ animationDelay: '0.1s' }}>
          <button
            onClick={() => { setSelectedCategory(null); setPage(1) }}
            className={`text-sm px-4 py-1.5 rounded-full border transition-all duration-200 ${
              !selectedCategory
                ? 'bg-accent text-paper border-accent'
                : 'border-border text-muted hover:text-ink hover:border-muted'
            }`}
          >
            全部
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => { setSelectedCategory(cat.id); setPage(1) }}
              className={`text-sm px-4 py-1.5 rounded-full border transition-all duration-200 ${
                selectedCategory === cat.id
                  ? 'bg-accent text-paper border-accent'
                : 'border-border text-muted hover:text-ink hover:border-muted'
              }`}
            >
              {cat.name}
            </button>
          ))}
        </div>
      )}

      {/* 文章列表 */}
      <div className="stagger divide-y divide-border">
        {articles.map((article, i) => (
          <ArticleCard key={article.id} article={article} index={i} />
        ))}
      </div>

      {/* 空状态 */}
      {articles.length === 0 && (
        <div className="text-center py-30">
          <p className="text-muted text-xl font-display italic">暂无文章</p>
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-4 mt-16">
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
  )
}

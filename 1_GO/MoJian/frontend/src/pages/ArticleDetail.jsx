import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import dayjs from 'dayjs'
import client from '../api/client'
import { useAuth } from '../context/AuthContext'

// ArticleDetail 文章详情页，沉浸式阅读体验
export default function ArticleDetail() {
  const { id } = useParams()
  const { user } = useAuth()
  const [article, setArticle] = useState(null)
  const [comments, setComments] = useState([])
  const [commentContent, setCommentContent] = useState('')
  const [replyTo, setReplyTo] = useState(null)
  const [guestLimit, setGuestLimit] = useState(0)
  const [isLimited, setIsLimited] = useState(false)

  // 加载文章详情
  useEffect(() => {
    client.get(`/articles/${id}`).then((res) => setArticle(res.data))
    client.get(`/articles/${id}/comments`).then((res) => {
      const data = res.data || {}
      setComments(data.comments || [])
      setGuestLimit(data.guest_comment_limit || 0)
      setIsLimited(data.is_limited || false)
    })
  }, [id])

  // 提交评论
  const handleSubmitComment = async (e) => {
    e.preventDefault()
    if (!commentContent.trim()) return
    const payload = { content: commentContent, article_id: Number(id) }
    if (replyTo) payload.parent_id = replyTo
    await client.post('/comments', payload)
    setCommentContent('')
    setReplyTo(null)
    const res = await client.get(`/articles/${id}/comments`)
    const data = res.data || {}
    setComments(data.comments || [])
  }

  if (!article) {
    return (
      <div className="max-w-3xl mx-auto px-6 py-30 text-center">
        <p className="text-muted font-display italic text-xl">加载中...</p>
      </div>
    )
  }

  return (
    <article className="max-w-3xl mx-auto px-6 pt-18 pb-10">
      {/* 文章头部 — 大标题 + 元信息 */}
      <header className="mb-16 article-enter">
        {article.category?.name && (
          <span className="text-xs font-semibold tracking-widest uppercase text-accent">
            {article.category.name}
          </span>
        )}
        <h1 className="font-display text-4xl md:text-6xl font-black leading-tight mt-3 mb-6">
          {article.title}
        </h1>
        <div className="flex items-center gap-4 text-sm text-muted">
          <span>{article.user?.username}</span>
          <span>&middot;</span>
          <span>{dayjs(article.published_at || article.created_at).format('YYYY年M月D日')}</span>
          <span>&middot;</span>
          <span>{article.view_count} 次阅读</span>
        </div>
        {article.tags?.length > 0 && (
          <div className="flex gap-2 mt-4">
            {article.tags.map((tag) => (
              <span key={tag.id} className="text-xs text-muted/70 border border-border rounded px-2.5 py-1">
                {tag.name}
              </span>
            ))}
          </div>
        )}
      </header>

      {/* 分隔线 */}
      <hr className="border-border mb-12" />

      {/* 文章正文 — Markdown 渲染 */}
      <div className="prose-ink article-enter" style={{ animationDelay: '0.15s' }}>
        <ReactMarkdown remarkPlugins={[remarkGfm]}>
          {article.content}
        </ReactMarkdown>
      </div>

      {/* 分隔线 */}
      <hr className="border-border mt-16 mb-12" />

      {/* 评论区 */}
      <section className="article-enter" style={{ animationDelay: '0.2s' }}>
        <h2 className="font-display text-3xl font-bold mb-8">
          评论<span className="text-accent">.</span>
        </h2>

        {/* 评论输入框 */}
        {user ? (
          <form onSubmit={handleSubmitComment} className="mb-10">
            {replyTo && (
              <div className="flex items-center gap-2 mb-3 text-sm text-muted">
                <span>回复评论 #{replyTo}</span>
                <button
                  type="button"
                  onClick={() => setReplyTo(null)}
                  className="text-accent hover:underline"
                >
                  取消
                </button>
              </div>
            )}
            <textarea
              value={commentContent}
              onChange={(e) => setCommentContent(e.target.value)}
              placeholder="写下你的想法..."
              rows={4}
              className="w-full border border-border rounded-lg p-4 bg-white/50 text-ink placeholder:text-muted/50 focus:outline-none focus:border-accent transition-colors resize-none"
            />
            <button
              type="submit"
              className="mt-3 bg-accent text-paper px-6 py-2.5 rounded text-sm font-semibold tracking-wide hover:bg-ink transition-colors duration-200"
            >
              发表评论
            </button>
          </form>
        ) : (
          <p className="text-muted mb-10 text-sm">
            <Link to="/login" className="text-accent hover:underline">登录</Link> 后即可评论
          </p>
        )}

        {/* 游客评论限制提示 */}
        {!user && isLimited && (
          <div className="mb-8 px-4 py-3 bg-accent/5 border border-accent/20 rounded-lg text-sm text-muted">
            游客最多只能查看 {guestLimit} 条评论，<Link to="/login" className="text-accent hover:underline">登录</Link> 后可查看全部评论
          </div>
        )}

        {/* 评论列表 */}
        <div className="space-y-6">
          {comments.filter((c) => !c.parent_id).map((comment) => (
            <div key={comment.id} className="ink-bar">
              <div className="flex items-center gap-3 mb-2">
                <span className="text-sm font-semibold">{comment.user?.username}</span>
                <span className="text-xs text-muted">
                  {dayjs(comment.created_at).format('M月D日 HH:mm')}
                </span>
              </div>
              <p className="text-ink leading-relaxed">{comment.content}</p>
              {user && (
                <button
                  onClick={() => setReplyTo(comment.id)}
                  className="text-xs text-muted hover:text-accent mt-2 transition-colors"
                >
                  回复
                </button>
              )}
              {/* 子评论 */}
              {comments.filter((c) => c.parent_id === comment.id).map((reply) => (
                <div key={reply.id} className="ml-6 mt-4 pl-4 border-l-2 border-border">
                  <div className="flex items-center gap-3 mb-1">
                    <span className="text-sm font-semibold">{reply.user?.username}</span>
                    <span className="text-xs text-muted">
                      {dayjs(reply.created_at).format('M月D日 HH:mm')}
                    </span>
                  </div>
                  <p className="text-ink text-sm leading-relaxed">{reply.content}</p>
                </div>
              ))}
            </div>
          ))}
          {comments.length === 0 && (
            <p className="text-muted italic font-display">暂无评论，来写第一条吧</p>
          )}
        </div>
      </section>
    </article>
  )
}

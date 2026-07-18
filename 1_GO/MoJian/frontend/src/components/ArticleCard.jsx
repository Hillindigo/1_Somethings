import { Link } from 'react-router-dom'
import dayjs from 'dayjs'

// ArticleCard 文章卡片组件，带左侧 ink-bar 标识
export default function ArticleCard({ article, index = 0 }) {
  return (
    <article
      className="article-enter group"
      style={{ animationDelay: `${index * 0.06}s` }}
    >
      <Link to={`/article/${article.id}`} className="block">
        <div className="ink-bar py-6 hover:bg-white/40 transition-colors duration-300 rounded-r-sm">
          {/* 分类 + 日期 */}
          <div className="flex items-center gap-3 mb-2">
            {article.category?.name && (
              <span className="text-xs font-semibold tracking-widest uppercase text-accent">
                {article.category.name}
              </span>
            )}
            <span className="text-xs text-muted">
              {dayjs(article.published_at || article.created_at).format('YYYY年M月D日')}
            </span>
          </div>

          {/* 标题 — 大号衬线体，视觉锚点 */}
          <h2 className="font-display text-2xl md:text-3xl font-bold leading-snug mb-2 group-hover:text-accent transition-colors duration-200">
            {article.title}
          </h2>

          {/* 摘要 */}
          {article.summary && (
            <p className="text-muted text-base leading-relaxed line-clamp-2">
              {article.summary}
            </p>
          )}

          {/* 标签 + 浏览 */}
          <div className="flex items-center gap-4 mt-3">
            {article.tags?.map((tag) => (
              <span key={tag.id} className="text-xs text-muted/70 border border-border rounded px-2 py-0.5">
                {tag.name}
              </span>
            ))}
            <span className="text-xs text-muted/50 ml-auto">
              {article.view_count} 次阅读
            </span>
          </div>
        </div>
      </Link>
    </article>
  )
}

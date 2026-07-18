import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import { ProtectedRoute, AdminRoute } from './components/ProtectedRoute'
import Home from './pages/Home'
import ArticleDetail from './pages/ArticleDetail'
import Login from './pages/Login'
import Register from './pages/Register'
import Profile from './pages/Profile'
import WriteArticle from './pages/WriteArticle'
import Dashboard from './pages/admin/Dashboard'
import Categories from './pages/admin/Categories'
import Tags from './pages/admin/Tags'
import Users from './pages/admin/Users'

// App 根组件，定义所有路由规则
export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        {/* 公开页面 */}
        <Route path="/" element={<Home />} />
        <Route path="/article/:id" element={<ArticleDetail />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* 需要登录 */}
        <Route path="/profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
        <Route path="/write" element={<ProtectedRoute><WriteArticle /></ProtectedRoute>} />
        <Route path="/article/:id/edit" element={<ProtectedRoute><WriteArticle /></ProtectedRoute>} />

        {/* 管理员 */}
        <Route path="/admin" element={<AdminRoute><Dashboard /></AdminRoute>} />
        <Route path="/admin/users" element={<AdminRoute><Users /></AdminRoute>} />
        <Route path="/admin/categories" element={<ProtectedRoute><Categories /></ProtectedRoute>} />
        <Route path="/admin/tags" element={<AdminRoute><Tags /></AdminRoute>} />
      </Route>
    </Routes>
  )
}

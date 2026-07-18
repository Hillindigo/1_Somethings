// Footer 底部栏，极简风格
export default function Footer() {
  return (
    <footer className="border-t border-border mt-30">
      <div className="max-w-5xl mx-auto px-6 py-10 flex justify-between items-center">
        <p className="text-sm text-muted font-body">
          墨笺 &copy; {new Date().getFullYear()}
        </p>
        <p className="text-xs text-muted/60 tracking-wider uppercase">
          以墨为记
        </p>
      </div>
    </footer>
  )
}

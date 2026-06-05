import type { ReactNode } from "react"

interface SectionHeaderProps {
  title: ReactNode
  action?: ReactNode
  className?: string
}

// 统一的区块标题行：左标题、右可选操作。替换 Home/History 中重复的 flex 标题行
export default function SectionHeader({ title, action, className = "" }: SectionHeaderProps) {
  return (
    <div className={`flex items-center justify-between mb-3.5 ${className}`}>
      <div className="text-lg font-bold">{title}</div>
      {action}
    </div>
  )
}

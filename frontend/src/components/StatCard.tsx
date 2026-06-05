import { cardShadow } from "@/components/Card"

type StatColor = "primary" | "mint" | "pink" | "yellow"

const barClass: Record<StatColor, string> = {
  primary: "after:bg-primary",
  mint: "after:bg-mint",
  pink: "after:bg-pink",
  yellow: "after:bg-yellow",
}

const valueClass: Record<StatColor, string> = {
  primary: "text-primary",
  mint: "text-mint",
  pink: "text-pink",
  yellow: "text-yellow",
}

interface StatCardProps {
  value: string | number
  label: string
  color: StatColor
  /** 数值偏长（如菜名）时用更小字号 */
  small?: boolean
}

// 记录页顶部统计卡：顶部彩条 + 数值 + 标签
export default function StatCard({ value, label, color, small }: StatCardProps) {
  return (
    <div className={`bg-card rounded-2xl p-4 ${cardShadow} text-center relative overflow-hidden border border-border after:content-[''] after:absolute after:top-0 after:left-0 after:right-0 after:h-[3px] ${barClass[color]}`}>
      <div className={`${small ? "text-base" : "text-2xl"} font-extrabold mb-0.5 truncate ${valueClass[color]}`}>{value}</div>
      <div className="text-xs text-text2">{label}</div>
    </div>
  )
}

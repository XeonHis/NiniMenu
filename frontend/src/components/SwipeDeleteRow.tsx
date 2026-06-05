import { useRef } from "react"
import { gsap, motionDuration } from "@/lib/gsap"

interface SwipeDeleteRowProps {
  children: React.ReactNode
  onDelete: () => void
  actionWidth: number
  className?: string
}

export default function SwipeDeleteRow({ children, onDelete, actionWidth, className = "" }: SwipeDeleteRowProps) {
  const rowRef = useRef<HTMLDivElement>(null)
  const startX = useRef(0)
  const offset = useRef(0)
  const dragging = useRef(false)

  function apply(value: number, animate = false) {
    offset.current = value
    gsap.to(rowRef.current, {
      x: value,
      duration: animate ? motionDuration(0.2) : 0,
      ease: "power2.out",
      overwrite: true,
      force3D: true,
    })
    if (rowRef.current) rowRef.current.style.zIndex = value < -10 ? "5" : "20"
  }

  function onTouchStart(e: React.TouchEvent) {
    dragging.current = true
    startX.current = e.touches[0].clientX
  }

  function onTouchMove(e: React.TouchEvent) {
    if (!dragging.current) return
    const diff = e.touches[0].clientX - startX.current
    const next = Math.min(0, offset.current + diff)
    apply(Math.max(next, -actionWidth))
    startX.current = e.touches[0].clientX
  }

  function onTouchEnd() {
    dragging.current = false
    apply(offset.current < -30 ? -actionWidth : 0, true)
  }

  return (
    <div className={`relative overflow-hidden ${className}`}>
      <button
        onClick={onDelete}
        onTouchEnd={(e) => { e.preventDefault(); onDelete() }}
        className="absolute bottom-0 right-0 top-0 z-10 flex items-center justify-center bg-red text-sm font-semibold text-white transition-colors active:bg-[#b91c1c]"
        style={{ width: actionWidth }}
      >
        删除
      </button>
      <div
        ref={rowRef}
        className="relative z-20 bg-card will-change-transform"
        style={{ touchAction: "pan-y" }}
        onTouchStart={onTouchStart}
        onTouchMove={onTouchMove}
        onTouchEnd={onTouchEnd}
      >
        {children}
      </div>
    </div>
  )
}

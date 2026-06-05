import gsap from "gsap"
import { useGSAP } from "@gsap/react"
import { Flip } from "gsap/Flip"
import { ScrollToPlugin } from "gsap/ScrollToPlugin"

gsap.registerPlugin(useGSAP, Flip, ScrollToPlugin)
gsap.defaults({ duration: 0.3, ease: "power2.out", overwrite: "auto" })

const REDUCE_MOTION_QUERY = "(prefers-reduced-motion: reduce)"

export { Flip, gsap, useGSAP }

export function prefersReducedMotion() {
  return typeof window !== "undefined" && window.matchMedia(REDUCE_MOTION_QUERY).matches
}

export function motionDuration(duration: number) {
  return prefersReducedMotion() ? 0 : duration
}

function getScrollParent(element: HTMLElement): HTMLElement | Window {
  let node = element.parentElement
  while (node && node !== document.body) {
    const style = window.getComputedStyle(node)
    const overflowY = style.overflowY
    if ((overflowY === "auto" || overflowY === "scroll") && node.scrollHeight > node.clientHeight) {
      return node
    }
    node = node.parentElement
  }
  return window
}

export function scrollToElement(
  element: HTMLElement | null,
  options: { block?: "start" | "center"; offset?: number; duration?: number } = {},
) {
  if (!element) return
  const { block = "start", offset = 0, duration = 0.42 } = options
  const scroller = getScrollParent(element)
  const targetRect = element.getBoundingClientRect()

  if (scroller === window) {
    const maxY = document.documentElement.scrollHeight - window.innerHeight
    const rawY = targetRect.top + window.scrollY - offset
    const y = block === "center"
      ? rawY - window.innerHeight / 2 + targetRect.height / 2
      : rawY
    gsap.to(window, {
      scrollTo: { y: gsap.utils.clamp(0, Math.max(0, maxY), y) },
      duration: motionDuration(duration),
      ease: "power2.out",
    })
    return
  }

  const container = scroller as HTMLElement
  const containerRect = container.getBoundingClientRect()
  const maxY = container.scrollHeight - container.clientHeight
  const rawY = targetRect.top - containerRect.top + container.scrollTop - offset
  const y = block === "center"
    ? rawY - container.clientHeight / 2 + targetRect.height / 2
    : rawY
  gsap.to(container, {
    scrollTo: { y: gsap.utils.clamp(0, Math.max(0, maxY), y) },
    duration: motionDuration(duration),
    ease: "power2.out",
  })
}

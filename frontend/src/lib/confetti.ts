import { gsap, motionDuration, prefersReducedMotion } from "@/lib/gsap"

const colors = ["#FF8C69", "#98D8C8", "#FFB7B2", "#FFD93D", "#8B5CF6", "#6EC6B8", "#F4A8A0"]

export function launchConfetti(container: HTMLDivElement | null, count = 80) {
  if (!container) return
  gsap.killTweensOf(container.children)
  container.innerHTML = ""

  if (prefersReducedMotion()) return

  container.style.perspective = "600px"
  const randomLeft = gsap.utils.random(10, 90, 1, true)
  const randomDrift = gsap.utils.random(-180, 180, 1, true)
  const randomSpin = gsap.utils.random(360, 1080, 1, true)
  const randomDuration = gsap.utils.random(2.2, 4.7, 0.1, true)
  const randomDelay = gsap.utils.random(0, 0.7, 0.02, true)
  const randomSize = gsap.utils.random(4, 11, 0.5, true)
  const pieces: HTMLDivElement[] = []

  for (let i = 0; i < count; i++) {
    const piece = document.createElement("div")
    const size = randomSize()
    const shape = Math.random()
    const width = shape < 0.55 ? size : Math.max(3, size * 0.35)
    const height = shape < 0.25 ? size : shape < 0.55 ? size : size * (1.2 + Math.random() * 0.8)
    const dir = Math.random() > 0.5 ? 1 : -1

    piece.dataset.drift = String(randomDrift())
    piece.dataset.spin = String(dir * randomSpin())
    piece.dataset.flip = String(dir * randomSpin())
    piece.dataset.duration = String(randomDuration())
    piece.dataset.delay = String(randomDelay())

    container.appendChild(piece)
    gsap.set(piece, {
      position: "absolute",
      left: `${randomLeft()}%`,
      top: "-6%",
      width,
      height,
      borderRadius: shape < 0.25 ? "50%" : "2px",
      backgroundColor: colors[i % colors.length],
      autoAlpha: 0,
      x: 0,
      y: "-10vh",
      rotation: 0,
      rotationY: 0,
      willChange: "transform, opacity",
      force3D: true,
    })
    pieces.push(piece)
  }

  gsap.timeline({
    onComplete: () => {
      pieces.forEach((piece) => piece.remove())
      container.style.perspective = ""
    },
  })
    .to(pieces, {
      autoAlpha: 1,
      duration: motionDuration(0.14),
      stagger: { each: 0.006, from: "random" },
    }, 0)
    .to(pieces, {
      x: (_index, target) => Number((target as HTMLElement).dataset.drift || 0),
      y: "118vh",
      rotation: (_index, target) => Number((target as HTMLElement).dataset.spin || 540),
      rotationY: (_index, target) => Number((target as HTMLElement).dataset.flip || 360),
      duration: (_index, target) => Number((target as HTMLElement).dataset.duration || 3.2),
      delay: (_index, target) => Number((target as HTMLElement).dataset.delay || 0),
      ease: "power2.in",
    }, 0)
    .to(pieces, {
      autoAlpha: 0,
      duration: 0.34,
      stagger: { each: 0.005, from: "random" },
    }, 1.8)
}

import { onUnmounted, ref, type Ref } from 'vue'

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  alpha: number
  alphaSpeed: number
  color: string
}

export function useLyricsParticles(canvasRef: Ref<HTMLCanvasElement | null>) {
  const paused = ref(false)
  let animFrameId: number | null = null
  let particles: Particle[] = []
  let ctx: CanvasRenderingContext2D | null = null

  const colors = [
    'rgba(29, 185, 84, 0.4)',
    'rgba(96, 165, 250, 0.3)',
    'rgba(168, 85, 247, 0.3)',
    'rgba(244, 114, 182, 0.25)',
  ]

  let resizeRaf: number | null = null

  function init() {
    const canvas = canvasRef.value
    if (!canvas) return
    ctx = canvas.getContext('2d')
    if (!ctx) return

    resize()
    window.addEventListener('resize', onResize)

    const count = 40
    particles = Array.from({ length: count }, () => createParticle(true))
    startLoop()
  }

  function onResize() {
    if (resizeRaf !== null) cancelAnimationFrame(resizeRaf)
    resizeRaf = requestAnimationFrame(resize)
  }

  function resize() {
    const canvas = canvasRef.value
    if (!canvas) return
    canvas.width = canvas.offsetWidth * devicePixelRatio
    canvas.height = canvas.offsetHeight * devicePixelRatio
  }

  function createParticle(randomY = false): Particle {
    const canvas = canvasRef.value
    const w = canvas?.offsetWidth || 400
    const h = canvas?.offsetHeight || 400
    return {
      x: Math.random() * w,
      y: randomY ? Math.random() * h : h + 20,
      vx: (Math.random() - 0.5) * 0.3,
      vy: -(0.2 + Math.random() * 0.4),
      size: 1.5 + Math.random() * 3,
      alpha: 0.1 + Math.random() * 0.4,
      alphaSpeed: 0.002 + Math.random() * 0.005,
      color: colors[Math.floor(Math.random() * colors.length)]!,
    }
  }

  function startLoop() {
    function draw() {
      animFrameId = window.requestAnimationFrame(draw)

      if (!ctx || !canvasRef.value || paused.value) return
      const canvas = canvasRef.value
      const w = canvas.offsetWidth
      const h = canvas.offsetHeight

      ctx.clearRect(0, 0, w * devicePixelRatio, h * devicePixelRatio)
      ctx.setTransform(devicePixelRatio, 0, 0, devicePixelRatio, 0, 0)

      for (const p of particles) {
        p.x += p.vx
        p.y += p.vy
        p.alpha += p.alphaSpeed
        if (p.alpha > 0.5 || p.alpha < 0.05) p.alphaSpeed *= -1

        if (p.y < -20) {
          Object.assign(p, createParticle(false))
        }

        ctx!.beginPath()
        ctx!.arc(p.x, p.y, p.size, 0, Math.PI * 2)
        ctx!.fillStyle = p.color.replace(/[\d.]+\)$/, `${p.alpha})`)
        ctx!.fill()
      }
    }
    animFrameId = window.requestAnimationFrame(draw)
  }

  function stop() {
    if (animFrameId !== null) {
      cancelAnimationFrame(animFrameId)
      animFrameId = null
    }
    if (resizeRaf !== null) {
      cancelAnimationFrame(resizeRaf)
      resizeRaf = null
    }
    window.removeEventListener('resize', onResize)
    particles = []
    ctx = null
  }

  onUnmounted(() => {
    stop()
  })

  return { init, stop, paused }
}

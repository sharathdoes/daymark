"use client"

import { useEffect, useRef, useState } from "react"
import type React from "react"
import { annotate } from "rough-notation"
import { type RoughAnnotation } from "rough-notation/lib/model"

type AnnotationAction =
  | "highlight"
  | "underline"
  | "box"
  | "circle"
  | "strike-through"
  | "crossed-off"
  | "bracket"

interface HighlighterProps {
  children: React.ReactNode
  action?: AnnotationAction
  color?: string
  strokeWidth?: number
  animationDuration?: number
  iterations?: number
  padding?: number
  multiline?: boolean
  isView?: boolean
}

export function Highlighter({
  children,
  action = "highlight",
  color = "#ffd1dc",
  strokeWidth = 1.5,
  animationDuration = 600,
  iterations = 2,
  padding = 2,
  multiline = true,
  isView = false,
}: HighlighterProps) {
  const elementRef = useRef<HTMLSpanElement>(null)
  const annotationRef = useRef<RoughAnnotation | null>(null)
  const [isInView, setIsInView] = useState(false)

  // Native IntersectionObserver — replaces motion/react useInView
  useEffect(() => {
    if (!isView) return
    const el = elementRef.current
    if (!el) return
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true)
          observer.disconnect()
        }
      },
      { rootMargin: "-10%" }
    )
    observer.observe(el)
    return () => observer.disconnect()
  }, [isView])

  const shouldShow = !isView || isInView

  useEffect(() => {
    const element = elementRef.current
    let resizeObserver: ResizeObserver | null = null

    if (shouldShow && element) {
      const annotation = annotate(element, {
        type: action,
        color,
        strokeWidth,
        animationDuration,
        iterations,
        padding,
        multiline,
      })

      annotationRef.current = annotation
      annotation.show()

      resizeObserver = new ResizeObserver(() => {
        annotation.hide()
        annotation.show()
      })
      resizeObserver.observe(element)
      resizeObserver.observe(document.body)
    }

    return () => {
      if (annotationRef.current) {
        annotationRef.current.remove()
        annotationRef.current = null
      }
      resizeObserver?.disconnect()
    }
  }, [shouldShow, action, color, strokeWidth, animationDuration, iterations, padding, multiline])

  return (
    <span ref={elementRef} className="relative inline-block bg-transparent">
      {children}
    </span>
  )
}

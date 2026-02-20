import { useState, useEffect, useRef, type RefObject } from 'react';

/**
 * Hook to detect if an element is in the viewport
 * Uses Intersection Observer for efficient detection
 */
export function useInView(
  rootMargin: string = '100px'
): [RefObject<HTMLDivElement | null>, boolean] {
  const ref = useRef<HTMLDivElement | null>(null);
  const [isInView, setIsInView] = useState(false);

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        // Once in view, stay in view (no need to re-observe)
        if (entry.isIntersecting) {
          setIsInView(true);
          observer.disconnect();
        }
      },
      {
        rootMargin,
      }
    );

    observer.observe(element);

    return () => observer.disconnect();
  }, [rootMargin]);

  return [ref, isInView];
}

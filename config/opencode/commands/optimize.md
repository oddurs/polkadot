---
description: Analyze and optimize code for performance
agent: build
---

Optimize the specified code for performance:

1. **Profile first** — Identify the actual bottleneck, don't guess
2. **Measure** — Establish a baseline (bundle size, query time, render count, etc.)
3. **Common optimizations to check**:
   - N+1 queries → batch/join
   - Missing database indexes
   - Unnecessary re-renders → React.memo, useMemo, useCallback (only where measured)
   - Large bundle imports → dynamic import / tree shaking
   - Synchronous blocking → async/streaming
   - Repeated computation → caching/memoization
   - Unoptimized images → next/image, WebP, lazy loading
   - Excessive DOM nodes → virtualization
4. **Apply** — Make targeted changes with explanations
5. **Don't over-optimize** — Only fix what's actually slow

Focus on: $ARGUMENTS

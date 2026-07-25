# Creator Analytics (/studio)

> Enhanced from existing `PageCreatorDashboard.vue`
> Focus: chart system (P10), data visualizations, accessible analytics

---

## Layout (Desktop — Full Width Max-7xl)

```
┌────────────────────────────────────────────────────────────┐
│  STUDIO  (badge)                                            │
│  Creator Dashboard                                          │
│  Analytics, earnings, audience, and content management      │
│                                                             │
│  [S Upload] [Refresh ↻]                                     │
│                                                             │
│  Tabs: [Overview] [Content] [Audience] [Earnings]          │
│                                                              │
│  ════════════════════ OVERVIEW TAB ════════════════════════ │
│                                                              │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐        │
│  │ Total Streams│ │  Listeners   │ │   Followers  │        │
│  │   128,432    │ │   4,892      │ │   1,234      │        │
│  │   +12% ↑     │ │   +8% ↑      │ │   +24% ↑     │        │
│  └──────────────┘ └──────────────┘ └──────────────┘        │
│                                                              │
│  ┌─── Streams Over Time (Line Chart Area) ───────────────┐ │
│  │  ┌──────────────────────────────────────────────────┐  │ │
│  │  │   📈  Streams per day (last 30 days)              │  │ │
│  │  │                                                    │  │ │
│  │  │   5,000 ┤          ╱╲                              │  │ │
│  │  │   4,000 ┤        ╱╱ ╲╲    ╱╲                      │  │ │
│  │  │   3,000 ┤      ╱╱   ╲╲  ╱╱ ╲╲  ╱╲                │  │ │
│  │  │   2,000 ┤   ╱╱      ╲╲╱╱   ╲╲╱╱ ╲╲               │  │ │
│  │  │   1,000 ┤ ╱╱         ╲╱      ╲╱   ╲╲             │  │ │
│  │  │       0 ┤──────────                      ╲╲──     │  │ │
│  │  │          └──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──┬──      │  │ │
│  │  │             6  8  10 12 14 16 18 20 22 24 26 28     │  │ │
│  │  └──────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌──────────────┬────────────────────────────────────────┐  │
│  │  Top Tracks   │   Audience Demographics (Donut)       │  │
│  │   ┌─────────┐│   ┌──────────────────────────────────┐ │  │
│  │   │ 1. Song ││   │        ╭───╮                      │ │  │
│  │   │ 2. Song ││   │   ╭────╯   ╰────╮                 │ │  │
│  │   │ 3. Song ││   │   │ Iran 68%    │                 │ │  │
│  │   │ 4. Song ││   │   │ US   12%    │                 │ │  │
│  │   │ 5. Song ││   │   │ Turkey 8%  │                 │ │  │
│  │   └─────────┘│   │   ╰────╮   ╭────╯                 │ │  │
│  │               │   │        ╰───╯                      │ │  │
│  │               │   └──────────────────────────────────┘ │  │
│  └──────────────┴────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

## Chart Component Library

Muse uses **Lottie + SVG** for charts (lightweight, no heavy charting library).
Each chart is a Vue component with Canvas2D fallback.

### Chart Types Required

| Chart | Component | Use Case | Accessibility |
|-------|-----------|----------|---------------|
| **Line** | `StreamsLineChart.vue` | Streams over time | `aria-label` describing trend + data table |
| **Bar** | `TopContentBar.vue` | Top tracks/albums | `role="img"` + `aria-label` with rank values |
| **Donut** | `DemographicsDonut.vue` | Audience by country | Segment labels + legend with percentages |
| **Heatmap** | `ActivityHeatmap.vue` | Listener activity by hour/day | Data table alternative |

### StreamsLineChart.vue (Spec)
```vue
<template>
  <div class="rounded-2xl bg-white/[0.06] p-6 ring-1 ring-white/[0.10]">
    <div class="mb-4 flex items-center justify-between">
      <div>
        <h3 class="text-sm font-semibold text-white">Streams Over Time</h3>
        <p class="text-xs text-white/40">Last 30 days</p>
      </div>
      <!-- Period selector -->
      <select aria-label="Select time period"
              class="rounded-lg bg-white/[0.06] px-3 py-1.5 text-xs text-white/70
                     outline-none focus-visible:ring-2 focus-visible:ring-[#1db954]"
              v-model="period">
        <option value="7d">7 days</option>
        <option value="30d">30 days</option>
        <option value="90d">90 days</option>
        <option value="1y">1 year</option>
      </select>
    </div>

    <!-- Chart canvas -->
    <div class="relative h-64 w-full" ref="chartContainer">
      <canvas ref="chartCanvas" class="h-full w-full"
              role="img"
              :aria-label="chartDescription" />
    </div>

    <!-- Screen reader data table (visually hidden) -->
    <div class="sr-only" role="table" aria-label="Stream data by date">
      <div role="row" v-for="point in data" :key="point.date">
        <span role="cell">{{ point.date }}</span>
        <span role="cell">{{ point.value }} streams</span>
      </div>
    </div>

    <!-- Percentage change -->
    <div class="mt-2 flex items-center gap-1 text-xs"
         :class="trend >= 0 ? 'text-[#1db954]' : 'text-red-400'">
      <Icon :name="trend >= 0 ? 'trending-up' : 'trending-down'"
            class="h-3 w-3" aria-hidden="true" />
      <span>{{ trend >= 0 ? '+' : '' }}{{ trend }}% vs last period</span>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex h-64 items-center justify-center">
      <SkeletonLoader variant="chart" class="h-full w-full" aria-label="Loading chart data" />
    </div>

    <!-- Empty state -->
    <div v-if="!loading && data.length === 0"
         class="flex h-64 flex-col items-center justify-center gap-2 text-sm text-white/40">
      <Icon name="bar-chart" class="h-8 w-8 text-white/20" aria-hidden="true" />
      <p>No data yet. Upload tracks to see your stats.</p>
    </div>
  </div>
</template>
```

### Chart Rendering Strategy
- **Canvas2D** for line and bar charts (performant, no DOM nodes per data point)
- **SVG** for donut charts (easy accessible labels)
- **Responsive**: `ResizeObserver` to redraw on container resize
- **Dark mode**: Chart colors use CSS variables:
  ```css
  --chart-line: #1db954;
  --chart-fill: rgba(29, 185, 84, 0.1);
  --chart-grid: rgba(255, 255, 255, 0.06);
  --chart-text: rgba(255, 255, 255, 0.5);
  ```
- **Reduced motion**: Static chart (no animated drawing) when `prefers-reduced-motion`

### Empty / Loading States for All Charts

| State | Visual |
|-------|--------|
| **Loading** | SkeletonLoader `variant="chart"` — shimmer rectangle with line placeholders |
| **Empty** (new creator) | Illustration + "Upload your first track to see analytics" |
| **Error** | "Could not load chart data" + Retry button |
| **Partial data** | Chart shows available data with note "Insufficient data for this period" |

## Metric Card Spec
```vue
<div class="rounded-2xl bg-white/[0.06] p-5 ring-1 ring-white/[0.10]">
  <p class="text-xs font-medium text-white/40">{{ label }}</p>
  <p class="mt-1 text-2xl font-black text-white">
    {{ formatNumber(value, locale) }}
  </p>
  <div v-if="change !== undefined"
       class="mt-1 flex items-center gap-1 text-xs"
       :class="change >= 0 ? 'text-[#1db954]' : 'text-red-400'">
    <Icon :name="change >= 0 ? 'arrow-up' : 'arrow-down'"
          class="h-3 w-3" aria-hidden="true" />
    <span>{{ Math.abs(change) }}%</span>
  </div>
</div>
```

## Keyboard Navigation for Charts

| Key | Action |
|-----|--------|
| Tab | Navigate between chart elements (period selector, data points) |
| Arrow keys | Navigate between days/bars within a chart |
| Enter/Space | Select a data point (show tooltip) |
| Escape | Dismiss tooltip |
| Focus visible | All chart controls have `focus-visible` ring |

## Persian/Arabic Specific

| Element | RTL Adjustment |
|---------|---------------|
| Chart labels | Persian month names (فروردین, اردیبهشت, ...) |
| Numbers | Persian digits in metric values |
| Currency | Toman ( تومان ) format via `Intl.NumberFormat('fa-IR')` |
| Tooltip direction | Tooltips open to start (right in RTL, left in LTR) |
| Period selector | Dropdown opens correctly in RTL |

## Data Flow
```
PageCreatorDashboard.vue
├── on mount → fetchCreatorStats(userId)
├── creatorStore.stats → metric cards
├── creatorStore.streamsChart → line chart data (date[] → value[])
├── creatorStore.topTracks → bar chart data
├── creatorStore.demographics → donut chart data
├── creatorStore.activityHeatmap → heatmap grid data
├── Refreshing:
│   └── refreshStats() → creatorStore.fetchAll(userId)
└── Period change:
    └── fetchStreams(userId, period) → redraw chart
```

## Accessibility Validation

| Criterion | Implementation |
|-----------|---------------|
| 1.1.1 Non-text | Each chart has `role="img"` + descriptive `aria-label` |
| 1.4.1 Use of color | Chart uses patterns + labels, not color alone |
| 1.4.11 Non-text contrast | Chart lines at least 3:1 against background |
| 2.1.1 Keyboard | All chart interactions usable without mouse |
| 4.1.2 Name, Role | Data table alternative uses `role="table"`, `role="row"`, `role="cell"` |
| 4.1.3 Status messages | Loading/error states announced via `aria-live="polite"` |

## Stub Chart Implementation (for dev/backlog)

> Charts are tracked as separate Vue components in the storybook backlog.
> Until implemented, each chart section shows a well-designed placeholder:

```vue
<!-- Fallback while chart components are built -->
<div class="flex h-64 items-center justify-center rounded-2xl bg-white/[0.06]">
  <div class="flex flex-col items-center gap-3 text-center">
    <svg class="h-12 w-12 text-white/10" ...><!-- chart icon --></svg>
    <p class="text-sm text-white/30">Analytics chart incoming</p>
    <p class="text-xs text-white/20">Visual stream data will appear here</p>
  </div>
</div>
```

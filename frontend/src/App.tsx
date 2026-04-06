import { useEffect, useMemo, useRef, useState } from 'react'

type SummarySignal = {
  symbol: string
  signal: string
  confidence: number
  reason: string
}

type SummaryResponse = {
  updatedAt: string
  market: {
    symbol: string
    price: number
    changePercent: number
  }
  signals: SummarySignal[]
}

type WatchlistItem = {
  symbol: string
  name: string
  price: number
  changePercent: number
  signal: string
}

type WatchlistResponse = {
  items: WatchlistItem[]
}

type FilingItem = {
  symbol: string
  title: string
  source: string
  publishedAt: string
  url: string
}

type FilingsResponse = {
  items: FilingItem[]
}

type RecommendationScores = {
  technical: number
  fundamental: number
  news: number
  risk: number
}

type RecommendationSources = {
  marketData: string
  news: string
  filings: string
  [key: string]: string
}

type RecommendationResponse = {
  updatedAt: string
  symbol: string
  action: 'buy' | 'hold' | 'sell' | string
  confidence: number
  scores: RecommendationScores
  reasons: string[]
  sources: RecommendationSources
  metadata?: {
    scoreModel?: string
    scoreVersion?: string
    generatedAt?: string
    [key: string]: unknown
  }
  [key: string]: unknown
}

type SectorItem = {
  key: string
  label: string
  symbols: string[]
}

const API_BASE = 'http://localhost:8080/api'
const POLL_INTERVAL_MS = 8000

const sectors: SectorItem[] = [
  {
    key: 'technology',
    label: 'Technology',
    symbols: ['AAPL', 'MSFT', 'NVDA', 'AMD', 'INTC', 'QCOM']
  },
  {
    key: 'energy',
    label: 'Energy',
    symbols: ['XOM', 'CVX', 'COP', 'SLB', 'EOG']
  },
  {
    key: 'oil-gas',
    label: 'Oil & Gas',
    symbols: ['XOM', 'CVX', 'COP', 'MPC', 'VLO', 'OXY']
  }
]

const allSymbols = Array.from(new Set(sectors.flatMap((sector) => sector.symbols)))

const fallbackSummary: SummaryResponse = {
  updatedAt: '2026-04-06T00:00:00Z',
  market: {
    symbol: 'AAPL',
    price: 189.42,
    changePercent: 1.34
  },
  signals: [
    {
      symbol: 'AAPL',
      signal: 'buy',
      confidence: 0.82,
      reason: 'Momentum pendapatan masih kuat'
    },
    {
      symbol: 'MSFT',
      signal: 'hold',
      confidence: 0.64,
      reason: 'Pertumbuhan cloud tetap stabil'
    },
    {
      symbol: 'NVDA',
      signal: 'buy',
      confidence: 0.91,
      reason: 'Permintaan chip AI terus meningkat'
    }
  ]
}

const fallbackWatchlist: WatchlistResponse = {
  items: [
    { symbol: 'AAPL', name: 'Apple Inc.', price: 189.42, changePercent: 1.34, signal: 'buy' },
    { symbol: 'MSFT', name: 'Microsoft', price: 412.18, changePercent: 0.88, signal: 'hold' },
    { symbol: 'NVDA', name: 'NVIDIA', price: 972.56, changePercent: 3.12, signal: 'buy' }
  ]
}

const fallbackFilings: FilingsResponse = {
  items: [
    { symbol: 'AAPL', title: 'Laporan Kuartalan', source: 'SEC', publishedAt: '2026-04-06T00:00:00Z', url: 'https://example.com' },
    { symbol: 'MSFT', title: '8-K: Pengumuman Produk', source: 'SEC', publishedAt: '2026-04-05T18:30:00Z', url: 'https://example.com' },
    { symbol: 'NVDA', title: 'Pernyataan Proksi', source: 'SEC', publishedAt: '2026-04-05T14:15:00Z', url: 'https://example.com' }
  ]
}

const fallbackRecommendation: RecommendationResponse = {
  updatedAt: '2026-04-06T00:00:00Z',
  symbol: 'AAPL',
  action: 'hold',
  confidence: 0,
  scores: {
    technical: 0,
    fundamental: 0,
    news: 0,
    risk: 0
  },
  reasons: ['Menunggu data rekomendasi dari backend'],
  sources: {
    marketData: 'Belum tersedia',
    news: 'Belum tersedia',
    filings: 'Belum tersedia'
  }
}

function createFallbackRecommendation(symbol: string): RecommendationResponse {
  return {
    ...fallbackRecommendation,
    symbol
  }
}

async function fetchJson<T>(path: string, fallback: T): Promise<T> {
  try {
    const response = await fetch(`${API_BASE}${path}`)
    if (!response.ok) {
      throw new Error(`Request failed: ${response.status}`)
    }
    return (await response.json()) as T
  } catch {
    return fallback
  }
}

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2
  }).format(value)
}

function formatPercent(value: number): string {
  const sign = value > 0 ? '+' : ''
  return `${sign}${value.toFixed(2)}%`
}

function formatDate(iso: string): string {
  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short',
    timeZone: 'UTC'
  }).format(new Date(iso))
}

function formatLocalTimestamp(date: Date | null): string {
  if (!date) return 'Belum diperbarui'
  return new Intl.DateTimeFormat('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'medium'
  }).format(date)
}

function signalTone(signal: string): string {
  const normalized = signal.toLowerCase()
  if (normalized === 'buy') return 'text-emerald-700 bg-emerald-50 ring-emerald-200'
  if (normalized === 'sell') return 'text-rose-700 bg-rose-50 ring-rose-200'
  return 'text-amber-700 bg-amber-50 ring-amber-200'
}

function actionTone(action: string): string {
  const normalized = action.toLowerCase()
  if (normalized === 'buy') return 'border-emerald-200 bg-emerald-50 text-emerald-700'
  if (normalized === 'sell') return 'border-rose-200 bg-rose-50 text-rose-700'
  return 'border-amber-200 bg-amber-50 text-amber-700'
}

function progressWidth(value: number): string {
  return `${Math.max(0, Math.min(100, value))}%`
}

function translateReason(reason: string): string {
  const replacements: Array<[RegExp, string]> = [
    [/Technical trend remains positive/gi, 'Tren teknikal masih positif'],
    [/Fundamental profile is stable/gi, 'Profil fundamental stabil'],
    [/Recent news flow is supportive/gi, 'Aliran berita terbaru mendukung'],
    [/Risk remains manageable/gi, 'Risiko masih terkendali'],
    [/Technical score is/gi, 'Skor teknikal adalah'],
    [/based on current price/gi, 'berdasarkan harga saat ini'],
    [/and previous close/gi, 'dan harga penutupan sebelumnya'],
    [/Fundamental score is/gi, 'Skor fundamental adalah'],
    [/for /gi, 'untuk '],
    [/on /gi, 'pada '],
    [/News score is/gi, 'Skor berita adalah'],
    [/from/gi, 'dari'],
    [/recent headlines/gi, 'headline terbaru'],
    [/Risk score is/gi, 'Skor risiko adalah'],
    [/Latest headline:/gi, 'Headline terbaru:'],
    [/Trend remains positive/gi, 'Tren masih positif'],
    [/Earnings momentum remains strong/gi, 'Momentum pendapatan masih kuat'],
    [/Cloud growth remains steady/gi, 'Pertumbuhan cloud tetap stabil'],
    [/Demand for AI accelerators continues to expand/gi, 'Permintaan chip AI terus meningkat']
  ]

  let translated = reason
  for (const [pattern, replacement] of replacements) {
    translated = translated.replace(pattern, replacement)
  }
  return translated
}

function formatRecommendationScore(value: number): string {
  return Number.isFinite(value) ? value.toFixed(0) : '-'
}

export default function App() {
  const [summary, setSummary] = useState<SummaryResponse>(fallbackSummary)
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>(fallbackWatchlist.items)
  const [filings, setFilings] = useState<FilingItem[]>(fallbackFilings.items)
  const [selectedSector, setSelectedSector] = useState(sectors[0].key)
  const [recommendationSymbol, setRecommendationSymbol] = useState(sectors[0].symbols[0])
  const [recommendation, setRecommendation] = useState<RecommendationResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [recommendationLoading, setRecommendationLoading] = useState(false)
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null)
  const [refreshError, setRefreshError] = useState<string | null>(null)
  const selectedSectorData = sectors.find((sector) => sector.key === selectedSector) ?? sectors[0]
  const recommendationSymbolRef = useRef(recommendationSymbol)

  useEffect(() => {
    recommendationSymbolRef.current = recommendationSymbol
  }, [recommendationSymbol])

  const headlineChange = useMemo(() => formatPercent(summary.market.changePercent), [summary.market.changePercent])

  async function loadDashboardSnapshot(symbol = recommendationSymbolRef.current, silent = false) {
    try {
      const [summaryData, watchlistData, filingsData, recommendationData] = await Promise.all([
        fetchJson<SummaryResponse>('/summary', summary),
        fetchJson<WatchlistResponse>('/watchlist', { items: watchlist }),
        fetchJson<FilingsResponse>('/filings', { items: filings }),
        fetchJson<RecommendationResponse>(`/recommendation?symbol=${encodeURIComponent(symbol)}`, createFallbackRecommendation(symbol))
      ])

      setSummary(summaryData)
      setWatchlist(watchlistData.items)
      setFilings(filingsData.items)
      setRecommendation(recommendationData)
      setLastUpdated(new Date())
      setRefreshError(null)
    } catch {
      if (!silent) {
        setRefreshError('Gagal memuat pembaruan terbaru. Menampilkan data terakhir yang tersedia.')
      }
    } finally {
      if (!silent) {
        setLoading(false)
      }
    }
  }

  useEffect(() => {
    let mounted = true

    ;(async () => {
      if (!mounted) return
      await loadDashboardSnapshot(recommendationSymbolRef.current)
    })()

    const interval = window.setInterval(() => {
      void loadDashboardSnapshot(recommendationSymbolRef.current, true)
    }, POLL_INTERVAL_MS)

    return () => {
      mounted = false
      window.clearInterval(interval)
    }
  }, [])

  useEffect(() => {
    if (!selectedSectorData.symbols.includes(recommendationSymbol)) {
      setRecommendationSymbol(selectedSectorData.symbols[0])
    }
  }, [selectedSector, selectedSectorData, recommendationSymbol])

  useEffect(() => {
    void loadDashboardSnapshot(recommendationSymbol, true)
  }, [recommendationSymbol])

  async function handleLoadRecommendation(symbol = recommendationSymbol) {
    setRecommendationLoading(true)
    try {
      const result = await fetchJson<RecommendationResponse>(
        `/recommendation?symbol=${encodeURIComponent(symbol)}`,
        createFallbackRecommendation(symbol)
      )
      setRecommendation(result)
      setLastUpdated(new Date())
      setRefreshError(null)
    } finally {
      setRecommendationLoading(false)
    }
  }

  async function handleSelectSector(value: string) {
    const sector = sectors.find((item) => item.key === value) ?? sectors[0]
    setSelectedSector(sector.key)
    setRecommendationSymbol(sector.symbols[0])
  }

  const displayedRecommendation = recommendation ?? createFallbackRecommendation(recommendationSymbol)

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <div className="app-shell mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <header className="app-header mb-8 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
            <div>
              <p className="app-subtitle text-sm font-medium uppercase tracking-[0.28em] text-sky-700/80">
                Dashboard intel pasar
              </p>
              <h1 className="app-title mt-2 text-3xl font-semibold tracking-tight text-slate-900 sm:text-4xl">
                Finance Agent
              </h1>
              <p className="app-subtitle mt-3 max-w-2xl text-sm leading-6 text-slate-600">
                Pantau watchlist, baca filing, dan lihat rekomendasi beli / tahan / jual secara live dari data backend yang terhubung ke Finnhub.
              </p>
            </div>
            <div className="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-600">
              <div className="font-medium text-slate-900">API Backend</div>
              <div className="font-mono text-xs text-slate-500">{API_BASE}</div>
            </div>
          </div>
          <div className="mt-4 flex flex-col gap-2 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-600 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <span className="font-medium text-slate-900">Terakhir diperbarui:</span> {formatLocalTimestamp(lastUpdated)}
            </div>
            <div className="text-xs uppercase tracking-[0.2em] text-slate-500">
              Auto-refresh aktif setiap {Math.round(POLL_INTERVAL_MS / 1000)} detik
            </div>
          </div>
          {refreshError ? (
            <div className="mt-4 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
              {refreshError}
            </div>
          ) : null}
        </header>

        <main className="space-y-6">
          <section className="dashboard-grid grid gap-4 md:grid-cols-3">
            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Pasar
              </div>
              <div className="dashboard-card__value mt-3 text-3xl font-semibold text-slate-900">
                {formatCurrency(summary.market.price)}
              </div>
              <div className="dashboard-card__meta mt-2 text-sm text-slate-600">
                {summary.market.symbol} · {headlineChange}
              </div>
            </article>

            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Watchlist
              </div>
              <div className="dashboard-card__value mt-3 text-3xl font-semibold text-slate-900">
                {watchlist.length}
              </div>
              <div className="dashboard-card__meta mt-2 text-sm text-slate-600">
                Simbol aktif yang sedang dipantau
              </div>
            </article>

            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Filing
              </div>
              <div className="dashboard-card__value mt-3 text-3xl font-semibold text-slate-900">
                {filings.length}
              </div>
              <div className="dashboard-card__meta mt-2 text-sm text-slate-600">
                Dokumen terbaru dari SEC
              </div>
            </article>
          </section>

          <section className="grid gap-6 lg:grid-cols-2">
            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Rekomendasi Live
              </div>

              <div className="mt-4 grid gap-3 sm:grid-cols-2">
                <select
                  value={selectedSector}
                  onChange={(event) => void handleSelectSector(event.target.value)}
                  className="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 outline-none focus:border-sky-400 focus:ring-2 focus:ring-sky-100"
                >
                  {sectors.map((sector) => (
                    <option key={sector.key} value={sector.key}>
                      {sector.label}
                    </option>
                  ))}
                </select>

                <select
                  value={recommendationSymbol}
                  onChange={(event) => setRecommendationSymbol(event.target.value)}
                  className="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 outline-none focus:border-sky-400 focus:ring-2 focus:ring-sky-100"
                >
                  {selectedSectorData.symbols.map((symbol) => (
                    <option key={symbol} value={symbol}>
                      {symbol}
                    </option>
                  ))}
                </select>
              </div>

              <div className="mt-3 flex flex-col gap-3 sm:flex-row">
                <button
                  type="button"
                  onClick={() => void handleLoadRecommendation(recommendationSymbol)}
                  className="rounded-2xl border border-sky-200 bg-sky-50 px-4 py-3 text-sm font-semibold text-sky-700 transition hover:bg-sky-100"
                >
                  {recommendationLoading ? 'Memuat…' : 'Segarkan rekomendasi'}
                </button>
                <div className="flex items-center text-sm text-slate-500">
                  Rekomendasi akan ikut diperbarui otomatis saat simbol berubah.
                </div>
              </div>

              <div className="mt-5 rounded-3xl border border-slate-200 bg-slate-50 p-5">
                <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <div className="text-xs uppercase tracking-[0.2em] text-slate-500">Simbol</div>
                    <div className="mt-1 text-2xl font-semibold text-slate-900">{displayedRecommendation.symbol}</div>
                  </div>
                  <div className={`rounded-full border px-4 py-2 text-sm font-semibold uppercase tracking-wide ${actionTone(displayedRecommendation.action)}`}>
                    {displayedRecommendation.action}
                  </div>
                </div>

                <div className="mt-5 grid gap-4 sm:grid-cols-2">
                  <div>
                    <div className="text-xs uppercase tracking-[0.2em] text-slate-500">Keyakinan</div>
                    <div className="mt-2 text-3xl font-semibold text-slate-900">{formatRecommendationScore(displayedRecommendation.confidence)}%</div>
                    <div className="mt-3 h-2 rounded-full bg-slate-200">
                      <div
                        className="h-2 rounded-full bg-sky-500 transition-all"
                        style={{ width: progressWidth(displayedRecommendation.confidence) }}
                      />
                    </div>
                  </div>
                  <div>
                    <div className="text-xs uppercase tracking-[0.2em] text-slate-500">Diperbarui</div>
                    <div className="mt-2 text-sm text-slate-600">{formatDate(displayedRecommendation.updatedAt)}</div>
                    <div className="mt-4 text-xs uppercase tracking-[0.2em] text-slate-500">Sumber</div>
                    <div className="mt-2 text-sm text-slate-600">
                      Market: {displayedRecommendation.sources.marketData}
                      <br />
                      News: {displayedRecommendation.sources.news}
                      <br />
                      Filing: {displayedRecommendation.sources.filings}
                    </div>
                    {displayedRecommendation.metadata?.scoreModel || displayedRecommendation.metadata?.scoreVersion ? (
                      <div className="mt-4 text-xs uppercase tracking-[0.2em] text-slate-500">
                        Metadata Skor
                        <div className="mt-2 normal-case tracking-normal text-slate-600">
                          {displayedRecommendation.metadata?.scoreModel ? `Model: ${displayedRecommendation.metadata.scoreModel}` : null}
                          {displayedRecommendation.metadata?.scoreModel && displayedRecommendation.metadata?.scoreVersion ? <br /> : null}
                          {displayedRecommendation.metadata?.scoreVersion ? `Versi: ${displayedRecommendation.metadata.scoreVersion}` : null}
                        </div>
                      </div>
                    ) : null}
                  </div>
                </div>

                <div className="mt-5 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
                  {[
                    { label: 'Teknikal', value: displayedRecommendation.scores.technical },
                    { label: 'Fundamental', value: displayedRecommendation.scores.fundamental },
                    { label: 'Berita', value: displayedRecommendation.scores.news },
                    { label: 'Risiko', value: displayedRecommendation.scores.risk }
                  ].map((score) => (
                    <div key={score.label} className="rounded-2xl border border-slate-200 bg-white p-4">
                      <div className="text-xs uppercase tracking-[0.2em] text-slate-500">{score.label}</div>
                      <div className="mt-2 text-2xl font-semibold text-slate-900">{formatRecommendationScore(score.value)}</div>
                    </div>
                  ))}
                </div>

                <div className="mt-5">
                  <div className="text-xs uppercase tracking-[0.2em] text-slate-500">Alasan</div>
                  <ul className="mt-3 space-y-2 text-sm text-slate-700">
                    {displayedRecommendation.reasons.map((reason) => (
                      <li key={reason} className="rounded-2xl border border-slate-200 bg-white px-4 py-3">
                        {translateReason(reason)}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </article>

            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Sinyal
              </div>
              <ul className="signal-list mt-4 space-y-3">
                {summary.signals.map((item) => (
                  <li key={item.symbol} className="signal-item flex items-start justify-between rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <div>
                      <div className="signal-item__symbol text-sm font-semibold text-slate-900">{item.symbol}</div>
                      <div className="mt-1 text-sm text-slate-600">{translateReason(item.reason)}</div>
                    </div>
                    <div className={`signal-item__signal rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-wide ring-1 ${signalTone(item.signal)}`}>
                      {item.signal} · {(item.confidence * 100).toFixed(0)}%
                    </div>
                  </li>
                ))}
              </ul>
            </article>
          </section>

          <section className="grid gap-6 lg:grid-cols-2">
            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Watchlist
              </div>
              <ul className="signal-list mt-4 space-y-3">
                {watchlist.map((item) => (
                  <li key={item.symbol} className="signal-item flex items-center justify-between rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <div>
                      <div className="signal-item__symbol text-sm font-semibold text-slate-900">{item.symbol}</div>
                      <div className="text-sm text-slate-600">{item.name}</div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-medium text-slate-900">{formatCurrency(item.price)}</div>
                      <div className="signal-item__signal mt-1 text-xs font-semibold uppercase tracking-wide text-slate-500">
                        {formatPercent(item.changePercent)} · {item.signal}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            </article>

            <article className="dashboard-card rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
              <div className="dashboard-card__title text-sm font-medium uppercase tracking-[0.2em] text-slate-500">
                Filing
              </div>
              <ul className="filing-list mt-4 space-y-3">
                {filings.map((item) => (
                  <li key={`${item.symbol}-${item.publishedAt}`} className="filing-item flex flex-col gap-2 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4 sm:flex-row sm:items-center sm:justify-between">
                    <div>
                      <div className="text-sm font-semibold text-slate-900">
                        {item.symbol} · {item.title}
                      </div>
                      <div className="mt-1 text-sm text-slate-600">
                        {item.source} · {formatDate(item.publishedAt)}
                      </div>
                    </div>
                    <a
                      href={item.url}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex items-center rounded-full border border-sky-200 bg-sky-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-sky-700 transition hover:bg-sky-100"
                    >
                      Buka
                    </a>
                  </li>
                ))}
              </ul>
            </article>
          </section>

          {loading ? <div className="text-sm text-slate-500">Memuat data dashboard…</div> : null}
          {!loading && allSymbols.length === 0 ? <div className="text-sm text-slate-500">Tidak ada simbol untuk ditampilkan.</div> : null}
        </main>
      </div>
    </div>
  )
}
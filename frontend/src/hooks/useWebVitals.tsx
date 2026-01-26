import { useEffect, useCallback, useRef } from 'react';
import { onCLS, onLCP, onFCP, onTTFB, onINP, type Metric } from 'web-vitals';

interface WebVitalsMetric {
  name: string;
  value: number;
  rating: 'good' | 'needs-improvement' | 'poor';
  delta: number;
  id: string;
}

interface WebVitalsState {
  CLS: WebVitalsMetric | null;
  INP: WebVitalsMetric | null;
  LCP: WebVitalsMetric | null;
  FCP: WebVitalsMetric | null;
  TTFB: WebVitalsMetric | null;
}

// Thresholds for Web Vitals (from Google)
const thresholds: Record<string, { good: number; poor: number }> = {
  CLS: { good: 0.1, poor: 0.25 },
  INP: { good: 200, poor: 500 },
  LCP: { good: 2500, poor: 4000 },
  FCP: { good: 1800, poor: 3000 },
  TTFB: { good: 800, poor: 1800 },
};

function getRating(name: string, value: number): 'good' | 'needs-improvement' | 'poor' {
  const threshold = thresholds[name as keyof typeof thresholds];
  if (!threshold) return 'good';

  if (value <= threshold.good) return 'good';
  if (value <= threshold.poor) return 'needs-improvement';
  return 'poor';
}

/**
 * Hook to monitor and report Web Vitals
 * Can be used to send metrics to analytics or display in a debug panel
 */
export function useWebVitals(
  onReport?: (metrics: WebVitalsState) => void
) {
  const metricsRef = useRef<WebVitalsState>({
    CLS: null,
    INP: null,
    LCP: null,
    FCP: null,
    TTFB: null,
  });

  const handleMetric = useCallback((metric: Metric) => {
    const webVitalMetric: WebVitalsMetric = {
      name: metric.name,
      value: metric.value,
      rating: getRating(metric.name, metric.value),
      delta: metric.delta,
      id: metric.id,
    };

    metricsRef.current = {
      ...metricsRef.current,
      [metric.name]: webVitalMetric,
    };

    // Log to console in development
    if (import.meta.env.DEV) {
      const color =
        webVitalMetric.rating === 'good'
          ? 'green'
          : webVitalMetric.rating === 'needs-improvement'
          ? 'orange'
          : 'red';

      console.log(
        `%c[Web Vital] ${metric.name}: ${metric.value.toFixed(2)} (${webVitalMetric.rating})`,
        `color: ${color}; font-weight: bold;`
      );
    }

    // Call the report callback
    onReport?.(metricsRef.current);
  }, [onReport]);

  useEffect(() => {
    // Register all Web Vitals observers
    onCLS(handleMetric);
    onINP(handleMetric);
    onLCP(handleMetric);
    onFCP(handleMetric);
    onTTFB(handleMetric);
  }, [handleMetric]);

  return metricsRef.current;
}

/**
 * Send Web Vitals to an analytics endpoint
 */
export function sendToAnalytics(metric: Metric) {
  const body = JSON.stringify({
    name: metric.name,
    value: metric.value,
    rating: getRating(metric.name, metric.value),
    delta: metric.delta,
    id: metric.id,
    page: window.location.pathname,
    timestamp: Date.now(),
  });

  // Use sendBeacon if available, fallback to fetch
  if (navigator.sendBeacon) {
    navigator.sendBeacon('/api/analytics/web-vitals', body);
  } else {
    fetch('/api/analytics/web-vitals', {
      method: 'POST',
      body,
      headers: { 'Content-Type': 'application/json' },
      keepalive: true,
    });
  }
}

/**
 * Component to display Web Vitals (development only)
 */
export function WebVitalsDebugPanel() {
  const metrics = useWebVitals();

  if (!import.meta.env.DEV) return null;

  const entries = Object.entries(metrics).filter(([, m]) => m !== null);

  if (entries.length === 0) return null;

  return (
    <div className="fixed bottom-4 left-4 bg-gray-900 text-white text-xs rounded-lg p-3 shadow-lg z-50 font-mono">
      <div className="font-bold mb-2">Web Vitals</div>
      {entries.map(([name, metric]) => (
        <div key={name} className="flex items-center gap-2">
          <span
            className={`w-2 h-2 rounded-full ${
              metric?.rating === 'good'
                ? 'bg-green-500'
                : metric?.rating === 'needs-improvement'
                ? 'bg-yellow-500'
                : 'bg-red-500'
            }`}
          />
          <span>{name}:</span>
          <span>{metric?.value.toFixed(2)}</span>
        </div>
      ))}
    </div>
  );
}

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import './index.css';

// Report Web Vitals in production
import { onCLS, onINP, onLCP, onFCP, onTTFB } from 'web-vitals';

if (import.meta.env.PROD) {
  const reportWebVitals = (metric: { name: string; value: number }) => {
    // Send to analytics endpoint
    console.log('[Web Vital]', metric.name, metric.value);
  };

  onCLS(reportWebVitals);
  onINP(reportWebVitals);
  onLCP(reportWebVitals);
  onFCP(reportWebVitals);
  onTTFB(reportWebVitals);
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
